package rtree

type item struct {
	id     uint
	box    rect
	father *node
	ctx    Context
}

var itemInstId = uint(0)

func newItem(box rect, ctx Context) *item {
	itemInstId++
	return &item{
		id:  itemInstId,
		box: box,
		ctx: ctx,
	}
}

func (i *item) ID() uint {
	return i.id
}

func (i *item) Context() Context {
	return i.ctx
}

func (i *item) Rect() (Point, Point) {
	return i.box.min, i.box.max
}

func (i *item) Box() *rect {
	return &i.box
}

func (i *item) ReCalcArea() {
}

func (i *item) SetFather(f *node) {
	i.father = f
}

func (i *item) SearchOverlap(r *rect, cb func(Item) bool, withBorder bool) bool {
	if i.box.isOverlap(r, withBorder) {
		return cb(i)
	}
	return true
}

func (i *item) Free() {
	i.box.clean()
	i.father = nil
	i.ctx = nil
}
