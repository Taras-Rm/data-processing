package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"processor/internal/es"
	"processor/internal/infra/csv"
	"processor/internal/service"

	"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Username: "elastic",
		Password: "changeme",
	})
	if err != nil {
		panic(err)
	}

	usersStore, err := es.NewUsersStore(client, "users")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	f, err := os.Create("users_output.csv")
	if err != nil {
		log.Fatalf("create csv: %v", err)
	}

	defer f.Close()

	usersCsv := csv.NewUsersCsv(f)

	usersService := service.NewUsersService(usersStore, &usersCsv)

	err = usersService.ProcessAll(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}
