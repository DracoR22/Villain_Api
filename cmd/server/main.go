package main

import (
	"log"

	"github.com/DracoR22/villain_api/app/api"
	"github.com/DracoR22/villain_api/storage"
)

// ------------------------------------------------//RUN SERVER//--------------------------------------------//
func main() {
	store, err := storage.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":3000", store)
	server.Run()
}
