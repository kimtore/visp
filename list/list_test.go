package list_test

import (
	"strconv"
	"testing"

	"github.com/ambientsound/visp/list"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	lst := list.New()

	dataset := []list.Row{
		list.NewRow("", map[string]string{
			"foo": "foo content",
			"bar": "nope",
		}),
		list.NewRow("", map[string]string{
			"baz": "foo",
			"bar": "nopenope",
		}),
		list.NewRow("", map[string]string{
			"bar": "hope",
		}),
	}

	for _, row := range dataset {
		lst.Add(row)
	}

	t.Run("dataset content is correct", func(t *testing.T) {
		assert.Equal(t, len(dataset), lst.Len())

		for i := 0; i < lst.Len(); i++ {
			assert.Equal(t, dataset[i], lst.Row(i))
		}

	})

	t.Run("column names are correct", func(t *testing.T) {
		names := lst.ColumnNames()
		assert.ElementsMatch(t, []string{"bar", "baz", "foo"}, names)
	})

	t.Run("columns have correct content and sizes", func(t *testing.T) {
		names := []string{"foo", "bar", "baz"}
		cols := lst.Columns(names)

		assert.Len(t, cols, len(names))

		assert.Equal(t, 11, cols[0].Median())
		assert.Equal(t, 11, cols[0].Avg())

		assert.Equal(t, 4, cols[1].Median())
		assert.Equal(t, 5, cols[1].Avg())

		assert.Equal(t, 3, cols[2].Median())
		assert.Equal(t, 3, cols[2].Avg())
	})
}

func TestListNextOf(t *testing.T) {
	lst := list.New()

	dataset := []list.Row{
		list.NewRow("", map[string]string{
			"foo": "x",
			"bar": "x",
		}),
		list.NewRow("", map[string]string{
			"foo": "x",
			"bar": "xyz",
		}),
		list.NewRow("", map[string]string{
			"foo": "foo",
			"bar": "x",
		}),
		list.NewRow("", map[string]string{
			"foo": "x",
			"bar": "x",
		}),
		list.NewRow("", map[string]string{
			"foo": "x",
			"bar": "x",
		}),
	}

	for _, row := range dataset {
		lst.Add(row)
	}

	t.Run("stop in the middle at positive direction", func(t *testing.T) {
		next := lst.NextOf([]string{"foo"}, 0, 1)
		assert.Equal(t, 2, next)

		next = lst.NextOf([]string{"foo"}, 1, 1)
		assert.Equal(t, 2, next)
	})

	t.Run("stop in the middle at negative direction", func(t *testing.T) {
		next := lst.NextOf([]string{"foo"}, 4, -1)
		assert.Equal(t, 3, next)
	})

	t.Run("stop at the first multi-tag diff", func(t *testing.T) {
		next := lst.NextOf([]string{"foo", "bar"}, 0, 1)
		assert.Equal(t, 1, next)
	})

	t.Run("stop at the end if no diff found", func(t *testing.T) {
		next := lst.NextOf([]string{"bar"}, 2, 1)
		assert.Equal(t, 5, next)
	})
}

func TestListInsert(t *testing.T) {
	dataset := make([]list.Row, 0)
	for i := 0; i < 6; i++ {
		id := strconv.Itoa(i)
		row := list.NewRow(id, map[string]string{
			"name": id,
		})
		dataset = append(dataset, row)
	}

	setup := func() (list.List, list.List) {
		lst := list.New()
		inserts := list.New()

		for i := 0; i < 3; i++ {
			// adds 0, 1, 2 to the list
			lst.Add(dataset[i])
		}

		for i := 3; i < 6; i++ {
			// adds 3, 4, 5 to the list
			inserts.Add(dataset[i])
		}

		return lst, inserts
	}

	assertOrder := func(t *testing.T, lst list.List, order []int) {
		actualOrder := make([]int, 0)
		for _, row := range lst.All() {
			id, _ := strconv.Atoi(row.ID())
			actualOrder = append(actualOrder, id)
		}
		assert.Equal(t, order, actualOrder)
	}

	t.Run("test insertion at beginning", func(t *testing.T) {
		lst, inserts := setup()
		assert.NoError(t, lst.InsertList(inserts, 0))
		assertOrder(t, lst, []int{3, 4, 5, 0, 1, 2})
	})

	t.Run("test insertion in the middle", func(t *testing.T) {
		lst, inserts := setup()
		assert.NoError(t, lst.InsertList(inserts, 1))
		assertOrder(t, lst, []int{0, 3, 4, 5, 1, 2})
	})

	t.Run("test insertion at end", func(t *testing.T) {
		lst, inserts := setup()
		assert.NoError(t, lst.InsertList(inserts, lst.Len()))
		assertOrder(t, lst, []int{0, 1, 2, 3, 4, 5})
	})

	t.Run("insert into empty list", func(t *testing.T) {
		_, inserts := setup()
		lst := list.New()
		assert.NoError(t, lst.InsertList(inserts, 0))
		assertOrder(t, lst, []int{3, 4, 5})
	})

	t.Run("negative out of range", func(t *testing.T) {
		lst, inserts := setup()
		assert.Error(t, lst.InsertList(inserts, -1))
	})

	t.Run("positive out of range", func(t *testing.T) {
		lst, inserts := setup()
		assert.Error(t, lst.InsertList(inserts, lst.Len()+1))
	})
}
