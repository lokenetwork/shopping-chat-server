package main

import (
	"log"
	"net/http"

	"./chat"

)


func main() {
	log.SetFlags(log.Lshortfile)


	// websocket server
	server := chat.NewServer("/buyer-entry","/shoper-entry")
	go server.Listen()

	// static files
	http.Handle("/", http.FileServer(http.Dir("/home/loken/Applications/Go/src/webroot")))

	log.Fatal(http.ListenAndServe(":9080", nil))

}
