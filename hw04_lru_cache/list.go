package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	m_front *ListItem
	m_back  *ListItem
	m_size  int
}

func (l list) Len() int {
	return l.m_size
}

func (l list) Front() *ListItem {
	return l.m_front
}

func (l list) Back() *ListItem {
	return l.m_back
}

func (l list) PushFront(v interface{}) *ListItem {
	l.m_size++
	node := new(ListItem)
	node.Value = v
	node.Next = l.m_front
	node.Prev = nil
	l.m_front.Prev = node
	return node
}

func (l list) PushBack(v interface{}) *ListItem {
	l.m_size++
	node := new(ListItem)
	node.Value = v
	node.Next = nil
	node.Prev = l.m_back
	l.m_back.Next = node
	return node
}

func (l list) Remove(i *ListItem) {
	l.m_size--
	if i.Prev == nil {
		l.m_front.Next.Prev = nil
		l.m_front = l.m_front.Next
	} else if i.Next == nil {
		l.m_back.Prev.Next = nil
		l.m_back = l.m_back.Prev
	}
}

func (l list) MoveToFront(i *ListItem) {

}

func NewList() List {
	return new(list)
}
