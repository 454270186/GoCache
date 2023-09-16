package list

import "fmt"

type List struct {
	head *Element
	tail *Element
	size int
}

type Element struct {
	Key       string
	Val       string
	Pre, Next *Element
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

func (l *List) Len() int {
	return l.size
}

func (l *List) Add(node *Element) {
	last := l.tail.Pre
	last.Next = node
	node.Pre = last
	node.Next = l.tail
	l.tail.Pre = node

	l.size++
}

func (l *List) Remove(node *Element) {
	pre := node.Pre
	next := node.Next
	pre.Next = next
	next.Pre = pre

	l.size--
}

func (l *List) RemoveHead() {
	if l.size <= 0 {
		return
	}

	l.Remove(l.head.Next)
}

func (l *List) MoveToTail(node *Element) {
	l.Remove(node)
	l.Add(node)
}

func (l *List) GetFirst() *Element {
	if l.Len() <= 1 {
		panic("list is empty")
	}

	return l.head.Next
}

func (l *List) Print() {
	cur := l.head.Next

	for cur != l.tail {
		fmt.Print("<", cur.Key,", ", cur.Val, ">")
		fmt.Print(" ")
		cur = cur.Next
	}

	fmt.Println()
}