# RingSlice
RingSlice provides a slice-like interface that supports queue io, random access, sorting, and binary search. Data can be dropped or appended from both head and tail like a queue but without any memory reallocate because of the data storage in a ring buffer.

# Concept
- RingSlice works like a slice and supports random access.
- RingSlice works like a slice and supports `sort` and `search`.
- RingSlice supports drop data from both sides (head and tail) of the slice and the released space can be reused. (If you drop the original slice head, that memory cannot be reused directly.)
- RingSlice can append data from both sides (head or tail) of the slice.
- RingSlice will not reallocate memory when `append` or `remove` items.

# Interface
```go
type RingSlice[T any] interface {
	// Interface sort.Interface contain
	// Len() int return items num that this ring-slice already hold
	// Less(i, j int) bool
	// Swap(i, j int)
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
```

# Usage

## Use as a slide window
```go
    r := NewRingSlice[int](5, nil)
    
    r.AppendTail([]int{1, 2, 3, 4, 5}...) // ring -> [1, 2, 3, 4, 5]
    r.RemoveHead(2) // ring -> [3, 4, 5]
    r.AppendTail(6, 7) // ring -> [3, 4, 5, 6, 7]
    r.RemoveTail(4) // ring -> [3]

    if r.RestSpace() >= 3 { // remember to check rest space if not sure how much space left
        r.AppendHead(15) // ring -> [15, 3]
        r.AppendHead(16) // ring -> [16, 15, 3]
        r.AppendHead(17) // ring -> [17, 16, 15, 3]
        r.Get(0) // 17
        r.Get(-1) // 3
        r.Get(6) // panic: out of range
        r.Set(3, 14) // ring -> [17, 16, 15, 14]
    }

    for i := 0; i < r.Len(); i++ {
        _ = r.Get(i)
    }

```

## To support sort/search
To support sort/search, you need to set a comparator.
Comparator receive two items to compare and return a boolean to represent the result.

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
