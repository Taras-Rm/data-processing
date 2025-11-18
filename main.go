package main

import (
	"context"
	"fmt"
	"processor/internal/es"

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

	_, err = usersStore.GetAll(context.Background())
	if err != nil {
		fmt.Println(err)
	}

}
