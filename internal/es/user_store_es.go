package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"processor/internal/domain"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type usersStoreEs struct {
	client    *elasticsearch.Client
	indexName string
}

func NewUsersStore(client *elasticsearch.Client, indexName string) (UsersStore, error) {
	s := &usersStoreEs{
		client:    client,
		indexName: indexName,
	}

	if err := s.ensureIndex(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *usersStoreEs) ensureIndex() error {
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

func (s *usersStoreEs) IndexBulk(ctx context.Context, users []domain.User) error {
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

func (s *usersStoreEs) GetAll(ctx context.Context) ([]domain.User, error) {
	batchSize := 2000
	scroll := time.Minute * 1

	users := make([]domain.User, 0)

	res, err := s.client.Search(s.client.Search.WithIndex(s.indexName), s.client.Search.WithScroll(time.Minute*1), s.client.Search.WithSize(batchSize))
	if err != nil || res.IsError() {
		return nil, fmt.Errorf("search users: %w", err)
	}

	defer res.Body.Close()

	type SearchResponse struct {
		ScrollId string `json:"_scroll_id"`
		Hits     struct {
			Hits []struct {
				Source domain.User `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	var searchRes SearchResponse

	if err := json.NewDecoder(res.Body).Decode(&searchRes); err != nil {
		return nil, fmt.Errorf("decode search users response: %w", err)
	}

	scrollId := searchRes.ScrollId

	for {
		if len(searchRes.Hits.Hits) == 0 {
			break
		}

		for _, hit := range searchRes.Hits.Hits {
			users = append(users, hit.Source)
		}

		scrollRes, err := s.client.Scroll(s.client.Scroll.WithScrollID(scrollId), s.client.Scroll.WithScroll(scroll))
		if err != nil {
			return nil, fmt.Errorf("scroll users error: %w", err)
		}

		defer scrollRes.Body.Close()

		if err := json.NewDecoder(scrollRes.Body).Decode(&searchRes); err != nil {
			return nil, fmt.Errorf("decode search users response: %w", err)
		}

		scrollId = searchRes.ScrollId

	}

	_, err = s.client.ClearScroll(s.client.ClearScroll.WithScrollID(scrollId))
	if err != nil {
		return nil, err
	}

	return users, nil
}
