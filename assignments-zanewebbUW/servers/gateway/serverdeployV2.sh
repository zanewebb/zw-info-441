sudo docker rm -f zanewebb
sudo docker pull zanewebb/zanewebbuw
sudo docker pull zanewebb/zanemysql
export TLSCERT=/etc/letsencrypt/live/api.zanewebb.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.zanewebb.me/privkey.pem
export DSN="root:GoodbyeWeekend@tcp(mysqlServer:3306)/db"
export SESSIONKEY=sessionkey
export MYSQL_ROOT_PASSWORD=GoodbyeWeekend
sudo docker run -d --name redisServer --network __networkname__ redis
sudo docker run -d --name mysqlServer --network gatewayNetwork -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD -e MYSQL_DATABASE=db zanewebb/zanemysql
sudo docker run -d -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY -e DSN=$DSN -e SESSIONKEY=$SESSIONKEY --name gateway zanewebb/zanewebbuw

