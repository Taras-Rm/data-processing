package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"processor/internal/domain"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type userStoreEs struct {
	client    *elasticsearch.Client
	indexName string
}

func NewUserStore(client *elasticsearch.Client, indexName string) (UserStore, error) {
	s := &userStoreEs{
		client:    client,
		indexName: indexName,
	}

	if err := s.ensureIndex(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *userStoreEs) ensureIndex() error {
	res, err := s.client.Indices.Exists([]string{s.indexName})
	if err != nil {
		return err
	}

	res.Body.Close()

	if res.StatusCode == 200 {
		return nil
	}

	// new index creating
	body := map[string]any{
		"mappings": map[string]any{
			"properties": map[string]any{
				"id":    map[string]any{"type": "keyword"},
				"name":  map[string]any{"type": "text"},
				"email": map[string]any{"type": "keyword"},
			},
		},
	}

	req := esapi.IndicesCreateRequest{
		Index: s.indexName,
		Body:  esutil.NewJSONReader(body),
	}

	createRes, err := req.Do(context.Background(), s.client)
	if err != nil {
		return err
	}

	defer createRes.Body.Close()

	if createRes.IsError() {
		return fmt.Errorf("error creating index: %s", createRes.String())
	}

	return nil
}

func (s *userStoreEs) IndexBulk(ctx context.Context, users []domain.User) error {
	var buf bytes.Buffer

	for _, user := range users {
		meta := []byte(fmt.Sprintf(`{ "index" : { "_index" : "%s" } }%s`, s.indexName, "\n"))

		data, err := json.Marshal(user)
		if err != nil {
			fmt.Printf("error marshaling user: %s\n", err)
			return err
		}

		data = append(data, byte('\n'))

		buf.Grow(len(meta) + len(data))

		buf.Write(meta)
		buf.Write(data)
	}

	res, err := s.client.Bulk(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		fmt.Printf("bulk indexing error: %s\n", res.String())
		return err
	}

	return nil
}

func (s *userStoreEs) Search(ctx context.Context) ([]domain.User, error) {
	return nil, nil
}
