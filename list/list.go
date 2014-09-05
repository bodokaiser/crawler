package list

import "container/list"

type List struct {
	list *list.List
}

func NewList() *List {
	return &List{
		list: list.New(),
	}
}

func (l *List) Add(i *Item) {
	if !l.Has(i) {
		l.list.PushFront(i)
	}
}

func (l *List) Has(i *Item) bool {
	for e := l.list.Front(); e != nil; e = e.Next() {
		if e.Value.(*Item).Origin() == i.Origin() {
			return true
		}
	}

	return false
}
