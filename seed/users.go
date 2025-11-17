package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"processor/internal/domain"
	"processor/internal/es"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

const (
	defaultCSVPath = "./seed/users_input.csv"

	defaultUsername = "elastic"
	defaultPassword = "changeme"

	defaultIndexName  = "users"
	timoutTimeSeconds = 30
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Seed users failed: %v", err)
	}
}

func run() error {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Username: defaultUsername,
		Password: defaultPassword,
	})
	if err != nil {
		return err
	}

	userStore, err := es.NewUserStore(client, defaultIndexName)
	if err != nil {
		return err
	}

	users, err := readFileUsers(defaultCSVPath)
	if err != nil {
		return fmt.Errorf("read file users: %w", err)
	}

	context, cancel := context.WithTimeout(context.Background(), timoutTimeSeconds)

	defer cancel()

	err = userStore.IndexBulk(context, users)
	if err != nil {
		return fmt.Errorf("bulk users index: %w", err)
	}

	log.Printf("Successfully seeded %d users into index %q", len(users), defaultIndexName)

	return nil
}

func readFileUsers(path string) ([]domain.User, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read all csv users records: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no users found in csv file: %s", defaultCSVPath)
	}

	users := make([]domain.User, 0, len(records))

	for rowIdx, userRow := range records {
		id := int64(rowIdx + 1)
		name := strings.TrimSpace(userRow[0])
		email := strings.TrimSpace(userRow[1])

		ageNum, err := strconv.Atoi(userRow[2])
		if err != nil {
			return nil, err
		}

		age := uint64(ageNum)

		users = append(users, domain.User{
			Id:    id,
			Name:  name,
			Email: email,
			Age:   age,
		})
	}

	return users, nil
}
