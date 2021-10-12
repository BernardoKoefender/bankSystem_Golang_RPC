package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"../administracao"
)

func main(){
	adm := new(administracao.Adm)

	// Register methods with rpc.
	err := rpc.Register(adm)
	if err != nil {
		log.Fatal("Error registering Adm ", err)
	}


	// allows lib to use http
	rpc.HandleHTTP()

	// Initialize the process that listens to all
	//comunication at give port, following
	//TCP protocoll.
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		log.Fatal("Listener error ", err)
	}
	log.Printf("Serving rpc on port: %d", 4040)


	// Activate the server defined above
	http.Serve(listener, nil)
	if err != nil{
		log.Fatal("Error serving: ", err)
	}

}