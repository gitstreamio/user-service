package client

import (
	"context"
	"gopkg.in/olivere/elastic.v5"
)

type ElasticClient struct {
	Client *elastic.Client
	Ctx    context.Context
}

func NewElasticClient(ctx context.Context, index string) (*ElasticClient, error) {
	// Create a client
	client, err := elastic.NewClient(elastic.SetURL("http://elastic:9200"))
	if err != nil {
		return nil, err
	}

	indexExists, err := client.IndexExists().Index([]string{index}).Do(ctx)
	if err != nil {
		return nil, err
	}
	if !indexExists {
		_, err = client.CreateIndex(index).Do(ctx)
		if err != nil {
			return nil, err
		}
	}

	return &ElasticClient{Client: client, Ctx: ctx}, err

}
