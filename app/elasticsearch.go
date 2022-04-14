package app

import (
	"context"
	"fmt"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"gitlab.ahsoq.com/Raven/scaffold/util"
)

type ElasticsearchConfig struct {
	Addresses []string
}

var ElasticsearchClient *elasticsearch.Client

func InitElasticsearchClient(path string) (err error) {
	Conf := elasticsearch.Config{}
	err = util.ParseConfig(path, Conf)
	if err != nil {
		return err
	}
	ElasticsearchClient, err = elasticsearch.NewClient(Conf)
	return err
}

func CloseElasticsearch() {
	req := esapi.IndexRequest{
		Index:      "test",
		DocumentID: strconv.Itoa(1),
		Refresh:    "true",
	}
	_, err := req.Do(context.Background(), ElasticsearchClient)
	if err != nil {
		fmt.Println(err.Error())
	}
}
