package ringslice

import (
	"errors"
	"fmt"
	"sort"
)

type RingSlice[T any] interface {
	// Interface sort.Interface contain
	// Len() int: return items num that this ring-slice already hold
	// Less(i, j int) bool: compare two items in sepcific position i/j, panic if comparator not set.
	// Swap(i, j int): swap items in sepcific position i/j.
	sort.Interface
	// Get item from index, index must less than length
	Get(index int) T
	// Set item to index, index must less than length
	Set(index int, item T)
	// AppendTail append items to tail,
	// Notice: panic if space not enough (Cap() - Len()) < len(items)
	AppendTail(items ...T)
	// AppendHead append item to head
	// Notice: panic if space not enough (Cap() - Len()) < 1
	AppendHead(item T)
	// RemoveHead remove n item from head, if n > Len() will remove all item.
	RemoveHead(n int)
	// RemoveTail remove n item from tail, if n > Len() will remove all item.
	RemoveTail(n int)
	// Cap returns the totla items num that this ring-slice can hold	Cap() int
	// RestSpace return the max items num that this ring-slice can append
	RestSpace() int
}

type ringSlice[T any] struct {
	head  int // point to head of array
	tail  int // point to tail of array
	size  int // size is the items num already hold
	slice []T
	less  func(i, j T) bool // Used for sort interface
}

var (
	NoSpaceError           = errors.New("No space")
	LessFuncNotImplemented = errors.New("Less func not implemented")
)

// NewRingSlice return a ring-slice instance,
// size is the max items num that this instance can hold
// less must be set if you want to use sort
func NewRingSlice[T any](size int, less func(i, j T) bool) RingSlice[T] {
	if size < 1 {
		panic(fmt.Errorf("%d is not a vaild space", size))
	}
	slice := make([]T, size, size)
	return &ringSlice[T]{
		head:  0,
		tail:  size - 1,
		slice: slice,
		less:  less,
		size:  0,
	}
}

func (r *ringSlice[T]) Get(index int) T {
	i := r.getAbsoluteIndex(index)
	h := int64(r.head)
	l := int64(r.Cap())
	return r.slice[int((h+i+l)%l)]
}

func (r *ringSlice[T]) Set(index int, item T) {
	i := r.getAbsoluteIndex(index)
	h := int64(r.head)
	l := int64(r.Cap())
	r.slice[int((h+i+l)%l)] = item
}

func (r *ringSlice[T]) AppendTail(items ...T) {
	if len(items) > r.RestSpace() {
		panic(NoSpaceError)
	}
	for i := range items {
		r.increaseTail()
		r.slice[r.tail] = items[i]
		r.size += 1
	}
}

func (r *ringSlice[T]) AppendHead(item T) {
	if 1 > r.RestSpace() {
		panic(NoSpaceError)
	}
	r.size += 1
	r.decreaseHead()
	r.slice[r.head] = item
}

func (r *ringSlice[T]) RemoveHead(n int) {
	if n > r.Len() {
		n = r.Len()
	}
	for i := 0; i < n; i++ {
		// clear data, avoid holding any unexpected pointer
		r.Set(0, *new(T))
		r.increaseHead()
		r.size -= 1
	}
}

func (r *ringSlice[T]) RemoveTail(n int) {
	if n > r.Len() {
		n = r.Len()
	}
	for i := 0; i < n; i++ {
		// clear data, avoid holding any unexpected pointer
		r.Set(r.Len()-1, *new(T))
		r.decreaseTail()
		r.size -= 1
	}
}

func (r *ringSlice[T]) Less(i, j int) bool {
	if r.less == nil {
		panic(LessFuncNotImplemented)
	}
	return r.less(r.Get(i), r.Get(j))
}

func (r *ringSlice[T]) Swap(rawI, rawJ int) {
	i := r.getAbsoluteIndex(rawI)
	j := r.getAbsoluteIndex(rawJ)
	h := int64(r.head)
	c := int64(r.Cap())
	realI := int((h + i) % c)
	realJ := int((h + j) % c)
	r.slice[realI], r.slice[realJ] = r.slice[realJ], r.slice[realI]
}

func (r *ringSlice[T]) Len() int {
	return r.size
}

func (r *ringSlice[T]) Cap() int {
	return len(r.slice)
}

func (r *ringSlice[T]) RestSpace() int {
	return r.Cap() - r.Len()
}

func (r *ringSlice[T]) increaseTail() {
	tail := int64(r.tail)
	l := int64(r.Cap())
	r.tail = int((tail + 1) % l)
}

func (r *ringSlice[T]) decreaseTail() {
	tail := int64(r.tail)
	l := int64(r.Cap())
	r.tail = int((tail - 1 + l) % l)
}

func (r *ringSlice[T]) increaseHead() {
	head := int64(r.head)
	l := int64(r.Cap())
	r.head = int((head + 1) % l)
}

func (r *ringSlice[T]) decreaseHead() {
	head := int64(r.head)
	l := int64(r.Cap())
	r.head = int((head - 1 + l) % l)
}

func (r *ringSlice[T]) checkIndex(index int) {
	if index >= r.Len() || index < -r.Len() {
		panic(fmt.Errorf("runtime error:index out of range [%d] with length %d", index, r.Len()))
	}
}

func (r *ringSlice[T]) getAbsoluteIndex(index int) int64 {
	r.checkIndex(index)
	return (int64(index) + int64(r.Len())) % int64(r.Len())
}
