# RingSlice
RingSlice is designed for slide window purposes, which provides a slice-like interface that supports random access, sorting, and binary search. Data can be dropped or appended from the head or tail like a queue but without any memory reallocate because of the data storage in a ring buffer.

# Concept
- RingSlice work like a slice, supports random access.
- RingSlice work like a slice, supports sort and search.
- RingSlice support drop data from both side (head and tail) of slice and all released space can be reused. (if you drop original slice head, those memory cannot reuse directly.)
- RingSlice can append data from both side (head or tail) of slice.
- RingSlice will not reallocate memory when append or remove item.

# Interface
```go
type RingSlice[T any] interface {
	// Interface sort.Interface contain
	// Len() int return items num that this ring-slice already hold
	// Less(i, j int) bool
	// Swap(i, j int)
	sort.Interface
	// Get item from index, index must less than Len()
	Get(index int) T
	// Set item to index, index must less than Len()
	Set(index int, item T)
	// AppendTail append items to tail,
    // Notice: panic if space not enough (Cap() - Len()) < len(items)
	AppendTail(items ...T)
	// AppendHead append item to head
    // Notice: panic if space not enough (Cap() - Len()) < 1
	AppendHead(item T)
	// RemoveHead remove n item from head
	RemoveHead(n int)
	// RemoveTail remove n item from tail
	RemoveTail(n int)
	// Cap return the max items num that this ring-slice can hold
	Cap() int
}
```

# Usage

## Only use as a slide window
```go
    r := NewRingSlice[int](5, nil)
    r.AppendTail([]int{1, 2, 3, 4, 5}...) // ring -> [1, 2, 3, 4, 5]
    r.RemoveHead(2) // ring -> [3, 4, 5]
    r.AppendTail(6, 7) // ring -> [3, 4, 5, 6, 7]
    r.RemoveTail(4) // ring -> [3]
    r.AppendHead(15) // ring -> [15, 3]
    r.AppendHead(16) // ring -> [16, 15, 3]
    r.AppendHead(17) // ring -> [17, 16, 15, 3]
    r.Get(0) // 17
    r.Get(-1) // 3
    r.Get(6) // panic: out of range
    r.Set(3, 14) // ring -> [17, 16, 15, 14]
```

## To support sort/search
To support sort/search, you need to set a comparators.
Comparators receive two items and compare them and return a Boolean to represent the result.

```go
import "sort"

// ...
    less := func(i, j float64) bool {
        return i < j
    }
    r := NewRingSlice[float64](5, less)
    r.AppendTail(5, 2, 3, 1, 4) // ring -> [5, 2, 3, 1, 4]
    sort.Sort(r) // ring -> [1, 2, 3, 4, 5]
```