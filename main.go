package main

import (
	"log"
	"github.com/havells/nlp/router"
	"github.com/micro/go-micro/web"
)

func main() {
	service := web.NewService(
		web.Name("nlp"),
		web.Address(servicePort),
	)
	//Initialize service
	service.Init()
	router := router.SetupRouter()

	//Setup http handler
	service.Handle("/", router)
	// Run server	
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

