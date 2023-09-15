package list

type List struct {
	head *Element
	tail *Element
	size int
}

type Element struct {
	Key       string
	Val       Value
	Pre, Next *Element
}

type Value interface {
	Len()
}

func New() *List {
	h, t := &Element{}, &Element{}
	h.Next, t.Pre = t, h

	return &List{
		size: 0,
		head: h,
		tail: t,
	}
}
