package rtree

import (
	"math"
)

// newNode 新节点
func newNode(h int) *node {
	nodeInstId++
	n := &node{
		id:       nodeInstId,
		children: make(map[uint]box, maxObject),
	}
	n.father = nil
	n.height = h
	return n
}

type node struct {
	id       uint
	box      rect
	father   *node
	height   int
	children map[uint]box
}

var nodeInstId = uint(0)

func (n *node) ID() uint {
	return n.id
}

// Box 返回节点的区域范围
func (n *node) Box() *rect {
	return &n.box
}

func (n *node) ReCalcArea() {
	first := true
	for _, c := range n.children {
		if first {
			n.box.set(c.Box().min, c.Box().max)
			first = false
		} else {
			n.box.merge(c.Box())
		}
	}
}

func (n *node) SetFather(f *node) {
	n.father = f
}

func (n *node) SearchOverlap(r *rect, cb func(Item) bool, withBorder bool) bool {
	if n.box.isOverlap(r, withBorder) {
		for _, c := range n.children {
			if !c.SearchOverlap(r, cb, withBorder) {
				return false
			}
		}
	}
	return true
}

func (n *node) Free() {
	n.father = nil
	n.box.clean()
	n.height = 0
}

// 插入数据
func (n *node) insert(obj *item) bool {
	if n.height == 0 {
		n.addChildren(obj)
		return n.box.merge(obj.Box())
	}
	best := n.chooseBest(obj.Box())
	grown := best.insert(obj)
	if grown {
		grown = n.box.merge(obj.Box())
	}
	if len(best.children) == maxObject {
		nn := newNode(best.height)
		best.splitTo(nn)
		n.addChildren(nn)
	}
	return grown
}

func (n *node) delete(obj *item) {
	n.delChildren(obj.ID())
}

func (n *node) addChildren(b box) {
	n.children[b.ID()] = b
	b.SetFather(n)
}

func (n *node) delChildren(id uint) {
	if c, ok := n.children[id]; ok {
		delete(n.children, id)
		n.checkReCalcArea(c.Box())
		c.Free()
		n.tryDown()
	}
}

func (n *node) checkReCalcArea(r *rect) {
	isOnEdge := n.box.onEdge(r)
	if isOnEdge {
		n.ReCalcArea()
		if n.father != nil {
			n.father.checkReCalcArea(r)
		}
	}
}

// 尝试降低树高
func (n *node) tryDown() {
	if n.father != nil {
		if len(n.children) == 0 {
			n.father.delChildren(n.id)
		} else {
			n.father.tryMergeChildren(n)
		}
	}
}

func (n *node) tryMergeChildren(c *node) {
	if len(c.children) <= minObject {
		if _, ok := n.children[c.id]; ok && len(n.children) > 1 {
			delete(n.children, c.id)
			best := n.chooseBest(c.Box())
			if len(best.children) < halfMaxObject {
				for id, cc := range c.children {
					best.addChildren(cc)
					delete(c.children, id)
				}
				best.Box().merge(c.Box())
				c.Free()
				n.tryDown()
			} else {
				n.children[c.id] = c
			}
		}
	}
}

// 选择最佳插入节点
func (n *node) chooseBest(box *rect) (best *node) {
	selectDiff := math.Inf(1)
	for _, nn := range n.children {
		nd := nn.(*node)
		c := nd.box.checkMergeRect(box)
		diff := c.area - nd.box.area

		//选择面积扩展最小且原始面积最小的且包含对象最少的
		if diff > selectDiff {
			continue
		}

		if diff == selectDiff {
			if best == nil {
				selectDiff = diff
				best = nd
				continue
			}
			if nd.box.area > best.box.area {
				continue
			}
			if nd.box.area == best.box.area {
				if len(nd.children) >= len(nd.children) {
					continue
				}
			}
		}
		selectDiff = diff
		best = nd
	}
	return
}

// 将节点下的数据分离到目标节点
func (n *node) splitTo(tn *node) {
	tn.box.set(n.box.max, n.box.max)
	n.box.set(n.box.min, n.box.min)

	w := make(map[uint]box, halfMaxObject)
	for id, c := range n.children {
		nMin := n.box.checkMergeRect(c.Box())
		nMax := tn.box.checkMergeRect(c.Box())
		if nMax.area > nMin.area {
			n.box = nMin
			continue
		}
		if nMax.area == nMin.area {
			w[id] = c
		} else {
			delete(n.children, id)
			tn.addChildren(c)
			tn.box = nMax
		}
	}
	for id, c := range w {
		if len(tn.children) < len(n.children) {
			delete(n.children, id)
			tn.addChildren(c)
			tn.box.merge(c.Box())
		} else {
			n.box.merge(c.Box())
		}
	}
}
