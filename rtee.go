package rtree

type rtree2D struct {
	root *node
	objs map[uint]*item
}

func (r *rtree2D) Insert(min, max Point, ctx Context) Item {
	obj := newItem(newRect(min, max), ctx)
	grown := r.root.insert(obj)
	if grown {
		r.root.Box().merge(obj.Box())
	}
	if len(r.root.children) == maxObject { //根节点满了,拆分
		left := r.root
		right := newNode(left.height)
		root := newNode(left.height + 1)

		left.splitTo(right)
		root.addChildren(left)
		root.addChildren(right)
		root.Box().merge(left.Box(), right.Box())
		r.root = root
	}

	r.objs[obj.id] = obj
	return obj
}

func (r *rtree2D) Delete(obj Item) (Context, bool) {
	return r.DeleteWithID(obj.ID())
}

func (r *rtree2D) DeleteWithID(id uint) (Context, bool) {
	obj, ok := r.objs[id]
	if !ok {
		return nil, false
	}
	ctx := obj.ctx
	obj.father.delete(obj)
	if r.root.height > 0 && len(r.root.children) == 1 {
		old := r.root
		var root *node
		for _, c := range old.children {
			root = c.(*node)
		}
		delete(old.children, root.id)
		root.father = nil
		r.root = root
		old.Free()
	}
	if r.root.height == 0 && len(r.root.children) == 0 {
		r.root.box.clean()
	}
	delete(r.objs, id)
	return ctx, true
}

func (r *rtree2D) SearchOverlap(min, max Point, cb func(Item) bool) {
	r.search(min, max, cb, false)
}

func (r *rtree2D) SearchOverlapAndBorder(min, max Point, cb func(Item) bool) {
	r.search(min, max, cb, true)
}

func (r *rtree2D) HasOverlap(min, max Point) bool {
	return r.hasOverlap(min, max, false)
}

func (r *rtree2D) HasOverlapAndBorder(min, max Point) bool {
	return r.hasOverlap(min, max, true)
}

func (r *rtree2D) Range(cb func(Item) bool) {
	for _, obj := range r.objs {
		if !cb(obj) {
			return
		}
	}
}

func (r *rtree2D) search(min, max Point, cb func(Item) bool, withBorder bool) {
	if len(r.objs) == 0 {
		return
	}
	b := newRect(min, max)
	if !r.root.SearchOverlap(&b, cb, withBorder) {
		return
	}
}

func (r *rtree2D) hasOverlap(min, max Point, withBorder bool) bool {
	if len(r.objs) == 0 {
		return false
	}
	b := newRect(min, max)
	if r.root.box.isOverlap(&b, withBorder) {
		return true
	}
	return false
}
