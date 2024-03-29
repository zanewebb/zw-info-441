package main

import (
    "log"
    "net/http"
	"os"
    "strings"
    "fmt"
    "path"
    "image/png"
)

//Thanks Mr. TA
const headerCORS = "Access-Control-Allow-Origin"
const corsAnyOrigin = "*"

//HelloHandler handles requests for the `/hello` resource
func identiconHandler(w http.ResponseWriter, r *http.Request) {
    //Thank you Mr. TA
    w.Header().Add(headerCORS, corsAnyOrigin)
    w.Header().Add("Content-Type", "image/png")

    fmt.Println("Path was:", r.URL.Path)

    gimmethatnameboi := path.Base(r.URL.Path)

    image := identicon(gimmethatnameboi)

    png.Encode(w,image)
    //w.Write([]byte(png.Encode(w,image)))
}

func nameHandler(w http.ResponseWriter, r *http.Request){
	//Thank you Mr. TA
	w.Header().Add(headerCORS, corsAnyOrigin)
    log.Printf("Received a request")
	//Not sure if this will work, found out online 
	//that you can use this to access the GET query params
	name := r.URL.Query().Get("name")
	if name == ""{ 
		w.Write([]byte("Hello World!"))
	} else{ 
		//incredibly inefficient concatenation because I couldn't find a really easy way
		var s []string
		s = append(s, "Hello, ")
		s = append(s, name)
		s = append(s, "!")
		//hopefully this join method actually concatenates these
		w.Write([]byte(strings.Join(s,"")))
	}
}
/*func NewMethodMux() *MethodMux {
    return &MethodMux{
        HandlerFuncs: map[string]func(http.ResponseWriter, *http.Request){},
    }
}*/

func main() {
    //get the value of the ADDR environment variable
    addr := os.Getenv("ADDR")

    //if it's blank, default to ":80", which means
    //listen port 80 for requests addressed to any host
    if len(addr) == 0 {
        addr = ":80"
    }

    //create a new mux (router)
    //the mux calls different functions for
    //different resource paths
    mux := http.NewServeMux()

    //tell it to call the HelloHandler() function
    //when someone requests the resource path `/hello`
    mux.HandleFunc("/", nameHandler)
    mux.HandleFunc("/identicon/",identiconHandler)

    //start the web server using the mux as the root handler,
    //and report any errors that occur.
    //the ListenAndServe() function will block so
    //this program will continue to run until killed
    log.Printf("server is listening at %s...", addr)
    log.Fatal(http.ListenAndServe(addr, mux))
}