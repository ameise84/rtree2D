package rtree

type Context = any

type Item interface {
	ID() uint
	Context() Context
	Rect() (min, max Point)
}

type RTree interface {
	// Insert 插入节点
	Insert(min, max Point, ctx Context) Item
	// Delete 移除节点
	Delete(obj Item) (Context, bool)
	// DeleteWithID 根据ID移除节点
	DeleteWithID(id uint) (Context, bool)
	// SearchOverlap 搜索重叠区域(不含相邻)
	SearchOverlap(min, max Point, cb func(Item) bool)
	// SearchOverlapAndBorder 搜索重叠区域以及接壤区域
	SearchOverlapAndBorder(min, max Point, cb func(Item) bool)
	// HasOverlap 是否有重叠区域(不含相邻)
	HasOverlap(min, max Point) bool
	// HasOverlapAndBorder 是否有重叠区域以及接壤区域
	HasOverlapAndBorder(min, max Point) bool
	// Range 扫描所有节点,cb返回参数标识是否继续扫描
	Range(cb func(Item) bool)
}

func New2D() RTree {
	return &rtree2D{
		root: newNode(0),
		objs: make(map[uint]*item, 64),
	}
}
