package rtree

func newRect(min, max Point) rect {
	p := rect{}
	p.set(min, max)
	return p
}

type rect struct {
	min  Point
	max  Point
	area float64
}

func (r *rect) copy(b *rect) {
	r.min = b.min
	r.max = b.max
	r.area = b.area
}

func (r *rect) clone() rect {
	return rect{r.min, r.max, r.area}
}

func (r *rect) clean() {
	r.min.X = 0
	r.max.X = 0
	r.min.Y = 0
	r.max.Y = 0
	r.area = 0
}

// 设置区域
func (r *rect) set(min, max Point) {
	r.min = min
	r.max = max
	r.calcArea()
}

// 计算r区域面积(不一定是真实面积,而是必备数据)
func (r *rect) calcArea() {
	r.area = float64((r.max.Y - r.min.Y) * (r.max.X - r.min.X))
}

// r区域和b区域是否有重叠部分
func (r *rect) isOverlap(b *rect, withBorder bool) bool {
	if withBorder {
		if b.min.X > r.max.X || b.max.X < r.min.X {
			return false
		}
		if b.min.Y > r.max.Y || b.max.Y < r.min.Y {
			return false
		}
	} else {
		if b.min.X >= r.max.X || b.max.X <= r.min.X {
			return false
		}
		if b.min.Y >= r.max.Y || b.max.Y <= r.min.Y {
			return false
		}
	}
	return true
}

// 将多个区域合并到r区域
func (r *rect) merge(args ...*rect) bool {
	if len(args) == 0 {
		return false
	}
	grown := false
	for _, arg := range args {
		if arg.min.X < r.min.X {
			r.min.X = arg.min.X
			grown = true
		}
		if arg.min.Y < r.min.Y {
			r.min.Y = arg.min.Y
			grown = true
		}
		if arg.max.X > r.max.X {
			r.max.X = arg.max.X
			grown = true
		}
		if arg.max.Y > r.max.Y {
			r.max.Y = arg.max.Y
			grown = true
		}
	}
	if grown {
		r.calcArea()
	}
	return grown
}

// 检查预计合并后的区域
func (r *rect) checkMergeRect(b *rect) rect {
	nr := r.clone()
	nr.merge(b)
	return nr
}

func (r *rect) onEdge(b *rect) bool {
	if r.min.X == b.min.X || r.max.X == b.max.X {
		return true
	}
	if r.min.Y == b.min.Y || r.max.Y == b.max.Y {
		return true
	}
	return false
}
