package main

import (
	"context"
	"fmt"
	"processor/internal/domain"
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

	uuserStore, err := es.NewUserStore(client, "users")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	users := []domain.User{
		{Name: "John Doe", Email: "john@gm.com", Age: 30},
		{Name: "Jane Doe", Email: "tom@t.com", Age: 34},
	}

	err = uuserStore.IndexBulk(context.Background(), users)
	if err != nil {
		fmt.Println(err)
	}
}
