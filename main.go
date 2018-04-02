package main

import (
	"context"
	"encoding/json"
	"log"

	elastic "github.com/olivere/elastic"
)

type School struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Street      string    `json:"street"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	ZIP         string    `json:"zip"`
	Location    []float64 `json:"location"`
	Fees        int64     `json:"fees"`
	Tags        []string  `json:"tags"`
	Rating      string    `json:"rating"`
}

func main() {

	esURL := "http://127.0.0.1:9200"
	ctx := context.Background()

	client, err := elastic.NewClient(
		elastic.SetURL(esURL),
		elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}

	info, code, err := client.Ping(esURL).Do(ctx)
	if err != nil {
		panic(err)
	}

	log.Println(info, code)

	esversion, err := client.ElasticsearchVersion(esURL)
	if err != nil {
		panic(err)
	}

	log.Println("ElasticSearch Version:", esversion)

	// Get a indexed document
	get1, err := client.Get().
		Index("schools").
		Type("school").
		Id("1").
		Do(ctx)

	if err != nil {
		panic(err)
	}

	if get1.Found {
		log.Printf("got %+v\n", get1)
	}

	// Search with a term query
	termQuery := elastic.NewQueryStringQuery("CBSE")
	searchResult, err := client.Search().
		Index("schools").
		Query(termQuery).
		From(0).
		Size(10).
		Pretty(true).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("Search Result: %#v\n", searchResult)

	if searchResult.Hits.TotalHits > 0 {

		log.Printf("Found %d hits\n", searchResult.Hits.TotalHits)

		for _, hit := range searchResult.Hits.Hits {
			var s School
			err := json.Unmarshal(*hit.Source, &s)
			if err != nil {
				panic(err)
			}

			log.Printf("Got School: %#v\n", s)
		}
	} else {
		log.Println("No hit found")
	}
}
