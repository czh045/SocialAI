package backend

import (
	"context"
	"fmt"

	"socialai/constants"
	"socialai/util"

	"github.com/olivere/elastic/v7"
)

var ESBackend *ElasticsearchBackend

type ElasticsearchBackend struct {
	client *elastic.Client
}

// InitElasticsearchBackend connects to Elasticsearch and creates the two
// course indexes if they are not already present.
func InitElasticsearchBackend(config *util.ElasticsearchInfo) {
	options := []elastic.ClientOptionFunc{
		elastic.SetURL(config.Address),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	}
	if config.Username != "" {
		options = append(options, elastic.SetBasicAuth(config.Username, config.Password))
	}

	client, err := elastic.NewClient(options...)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	createIndexIfMissing(client, ctx, constants.PostIndex, `{
		"mappings": {
			"properties": {
				"id": {"type": "keyword"},
				"user": {"type": "keyword"},
				"message": {"type": "text"},
				"url": {"type": "keyword", "index": false},
				"type": {"type": "keyword", "index": false}
			}
		}
	}`)
	createIndexIfMissing(client, ctx, constants.UserIndex, `{
		"mappings": {
			"properties": {
				"username": {"type": "keyword"},
				"password": {"type": "keyword"},
				"age": {"type": "long", "index": false},
				"gender": {"type": "keyword", "index": false}
			}
		}
	}`)

	ESBackend = &ElasticsearchBackend{client: client}
	fmt.Println("Elasticsearch indexes are ready.")
}

func createIndexIfMissing(client *elastic.Client, ctx context.Context, index, mapping string) {
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		panic(err)
	}
	if exists {
		return
	}

	if _, err := client.CreateIndex(index).Body(mapping).Do(ctx); err != nil {
		panic(err)
	}
}

func (backend *ElasticsearchBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
	return backend.client.Search().
		Index(index).
		Query(query).
		Sort("_doc", false).
		Size(100).
		Do(context.Background())
}

func (backend *ElasticsearchBackend) SaveToES(document interface{}, index, id string) error {
	_, err := backend.client.Index().
		Index(index).
		Id(id).
		BodyJson(document).
		Refresh("wait_for").
		Do(context.Background())
	return err
}

func (backend *ElasticsearchBackend) DeleteFromES(query elastic.Query, index string) error {
	_, err := backend.client.DeleteByQuery().
		Index(index).
		Query(query).
		Refresh("true").
		Do(context.Background())
	return err
}
