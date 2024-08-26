package example

import (
	"github.com/ameise84/rtree"
	"math/rand"
	"testing"
	"time"
)

func TestRtree(t *testing.T) {
	rand.Seed(time.Now().Unix())

	r := rtree.New2D()
	var item []rtree.Item
	row, col := 20, 50
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			it := r.Insert(rtree.Point{X: i * 10, Y: j * 10}, rtree.Point{X: i*10 + 10, Y: j*10 + 10}, struct {
				x int
				y int
				a int
				b int
			}{i * 10, j * 10, i*10 + 10, j*10 + 10})
			item = append(item, it)
		}
	}

	r.SearchOverlapAndBorder(rtree.Point{}, rtree.Point{X: 10, Y: 10}, func(item rtree.Item) bool {
		return true
	})

	rand.Shuffle(len(item), func(i, j int) {
		item[i], item[j] = item[j], item[i]
	})

	for _, it := range item {
		r.DeleteWithID(it.ID())
	}
}
