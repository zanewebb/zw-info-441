package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/assignments-zanewebbUW/servers/gateway/models/users"
	"github.com/assignments-zanewebbUW/servers/gateway/sessions"
	"github.com/go-redis/redis"

	"github.com/assignments-zanewebbUW/servers/gateway/handlers"
)
import _ "github.com/go-sql-driver/mysql"

//Create docker  network
// docker network create __networkname__
// docker run -d --name redisServer --network __networkname__ redis

// Run docker container for mysql server ??
// sudo docker run -d --name mysqlServer --network gatewayNetwork -e MYSQL_ROOT_PASSWORD=PASS -e MYSQL_DATABASE=db zanewebb/zanemysql

//DSN will be something like username:password@protocol(address)/dbname
//							root:PASSWORD@TCP(dockerhostname)/dbname

func testHandler(w http.ResponseWriter, r *http.Request) {
	//log.Printf("Received a request and handled with testHandler")
	w.Write([]byte("Handled the test request"))
}

func main() {

	ADDR := os.Getenv("ADDR")
	if len(ADDR) == 0 {
		ADDR = ":443"
		//ADDR = ":8888"
	}

	TLSCERT := os.Getenv("TLSCERT")
	if len(TLSCERT) == 0 {
		fmt.Println("TLSCERT env variable was not set")
		os.Exit(1)
	}

	TLSKEY := os.Getenv("TLSKEY")
	if len(TLSKEY) == 0 {
		fmt.Println("TLSKEY env variable was not set")
		os.Exit(1)
	}

	sessionkey := os.Getenv("SESSIONKEY")
	if len(sessionkey) == 0 {
		fmt.Println("SESSIONKEY env variable was not set")
		os.Exit(1)
	}

	redisaddr := os.Getenv("REDISADDR")
	if len(redisaddr) == 0 {
		//redisaddr = "172.17.0.2:6379"
		redisaddr = "redisServer:6379"
	}

	//3306
	dsn := os.Getenv("DSN")
	if len(dsn) == 0 {
		fmt.Println("DSN env variable was not set")
		os.Exit(1)
	}

	//Create DB object from SQL DB
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Error opening the database: %v", err)
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("Error opening the database: %v", err)
		os.Exit(1)
	}

	//When comeplete, close the db
	defer db.Close()

	//Create mysqlstore
	usersStore := users.NewMySQLStore(db)

	//Create redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisaddr,
	})

	//Create redisstore
	sessionStore := sessions.NewRedisStore(redisClient, time.Hour)

	//Create context
	context := handlers.NewContext(sessionkey, sessionStore, usersStore)

	//Initialize the tree on server startup
	context.UsersStore.PopulateTrie()

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/summary", handlers.SummaryHandler)
	mux.HandleFunc("/v1/users", context.UsersHandler)
	mux.HandleFunc("/v1/users/", context.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", context.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", context.SpecificSessionHandler)
	mux.HandleFunc("/v1/test", testHandler)

	wrappedMux := handlers.NewCors(mux)

	log.Printf("Server running and listening on %s", ADDR)
	//log.Fatal(http.ListenAndServe(ADDR, wrappedMux))
	log.Fatal(http.ListenAndServeTLS(ADDR, TLSCERT, TLSKEY, wrappedMux))
}
