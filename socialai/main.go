package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"socialai/backend"
	"socialai/handler"
	"socialai/util"
)

func main() {
	fmt.Println("started-service")

	configPath := os.Getenv("SOCIALAI_CONFIG")
	if configPath == "" {
		configPath = "conf/deploy.yml"
	}

	config, err := util.LoadApplicationConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	backend.InitElasticsearchBackend(config.ElasticsearchConfig)
	if err := backend.InitMediaBackend(config); err != nil {
		log.Fatalf("failed to initialize media storage: %v", err)
	}

	address := os.Getenv("PORT")
	if address == "" {
		address = "8080"
	}

	log.Printf("Social AI backend is listening on :%s", address)
	log.Fatal(http.ListenAndServe(":"+address, handler.InitRouter(config.TokenConfig)))
}
