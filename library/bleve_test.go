package library_test

import (
	"os"
	"testing"

	"github.com/ambientsound/visp/library"
	"github.com/blevesearch/bleve/v2"
)

func TestBleve(t *testing.T) {
	const path = "/tmp/bleve.index"

	os.RemoveAll(path)

	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(path, mapping)
	if err != nil {
		panic(err)
	}

	track := library.Track{
		Title: "hei du der",
	}
	err = index.Index("my-id", track)
	if err != nil {
		panic(err)
	}

	query := bleve.NewMatchQuery("hei")
	req := bleve.NewSearchRequest(query)
	res, _ := index.Search(req)
	for _, hit := range res.Hits {
		t.Log(hit)

		// hit.Fields
	}

	index.Close()
	// err = index
}
