package rtree

type box interface {
	ID() uint
	Box() *rect
	ReCalcArea()
	SetFather(*node)
	Free()
	SearchOverlap(*rect, func(Item) bool, bool) bool
}
