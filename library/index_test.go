package library_test

import (
	"testing"

	"github.com/ambientsound/visp/library"
	"github.com/ambientsound/visp/list"
	"github.com/stretchr/testify/assert"
)

func TestBleve(t *testing.T) {
	idx, err := library.NewInMemory()
	if err != nil {
		panic(err)
	}

	input := list.New()
	input.Add(list.NewRow("foobar", list.DataTypeTrack, map[string]string{
		"some": "verynicestring",
		"data": "should-be-stored",
	}))
	input.Add(list.NewRow("baz", list.DataTypeTrack, map[string]string{
		"some": "not so nice string",
		"data": "save this please",
	}))

	err = idx.Add(input)
	if err != nil {
		panic(err)
	}

	// Perform a query for a sample string.
	result, err := idx.Query("ver")
	if err != nil {
		panic(err)
	}

	t.Logf("%d results", result.Len())
	for i, row := range result.All() {
		t.Logf("%003d %-15s %+v", i+1, row.ID(), row.Fields())
	}

	assert.Equal(t, 1, result.Len())

	// Query for row data and check that it is equal to what we stored.
	row, err := idx.QueryID("baz")
	if err != nil {
		panic(err)
	}

	assert.Equal(t, input.RowByID("baz").Fields(), row.Fields())
}
