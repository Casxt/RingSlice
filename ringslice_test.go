package ringslice

import (
	"sort"
	"testing"
)

func TestRingSlice_AppendHead(t *testing.T) {
	tests := []struct {
		name  string
		size  int
		items []int
		want  []int
	}{
		{
			name:  "Full filled",
			size:  1,
			items: []int{1},
			want:  []int{1},
		}, {
			name:  "Not full",
			size:  5,
			items: []int{2, 3, 1},
			want:  []int{1, 3, 2},
		}, {
			name:  "Empty",
			size:  32,
			items: []int{},
			want:  []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRingSlice[int](tt.size, nil)
			for _, item := range tt.items {
				r.AppendHead(item)
			}
			check(t, r, tt.want)
		})
	}
}

func TestRingSlice_AppendTail(t *testing.T) {
	tests := []struct {
		name  string
		size  int
		items []int
		want  []int
	}{
		{
			name:  "Full filled",
			size:  1,
			items: []int{1},
			want:  []int{1},
		}, {
			name:  "Not full",
			size:  5,
			items: []int{2, 3, 1},
			want:  []int{2, 3, 1},
		}, {
			name:  "Empty",
			size:  32,
			items: []int{},
			want:  []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRingSlice[int](tt.size, nil)
			r.AppendTail(tt.items...)
			check(t, r, tt.want)
		})
	}
}

func TestRingSlice_Sort(t *testing.T) {
	tests := []struct {
		name  string
		size  int
		items []int
		want  []int
	}{
		{
			name:  "Ordered filled",
			size:  10,
			items: []int{1, 3, 6, 9, 10},
			want:  []int{1, 3, 6, 9, 10},
		}, {
			name:  "Random filled",
			size:  10,
			items: []int{3, 9, 1, 6, 10},
			want:  []int{1, 3, 6, 9, 10},
		}, {
			name:  "Reverse filled",
			size:  10,
			items: []int{10, 9, 6, 3, 1},
			want:  []int{1, 3, 6, 9, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRingSlice[int](tt.size, func(i, j int) bool {
				return i < j
			})
			r.AppendTail(tt.items...)
			sort.Sort(r)
			if !sort.IsSorted(r) {
				t.Errorf("got not sorted, expecting sorted")
			}
			check(t, r, tt.want)
		})
	}
}

func TestRingSlice_Stable(t *testing.T) {
	tests := []struct {
		name  string
		size  int
		items []int
		want  []int
	}{
		{
			name:  "Ordered filled",
			size:  10,
			items: []int{1, 3, 6, 9, 10},
			want:  []int{1, 3, 6, 9, 10},
		}, {
			name:  "Random filled",
			size:  10,
			items: []int{3, 9, 1, 6, 10},
			want:  []int{1, 3, 6, 9, 10},
		}, {
			name:  "Reverse filled",
			size:  10,
			items: []int{10, 9, 6, 3, 1},
			want:  []int{1, 3, 6, 9, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRingSlice[int](tt.size, func(i, j int) bool {
				return i < j
			})
			r.AppendTail(tt.items...)
			sort.Stable(r)
			if !sort.IsSorted(r) {
				t.Errorf("got not sorted, expecting sorted")
			}
			check(t, r, tt.want)
		})
	}
}

func TestRingSlice_RemoveHead(t *testing.T) {
	tests := []struct {
		name  string
		size  int
		items []int
		popn  int
		want  []int
	}{
		{
			name:  "Pop all",
			size:  1,
			items: []int{1},
			popn:  1,
			want:  []int{},
		}, {
			name:  "Pop some",
			size:  5,
			items: []int{2, 3, 1},
			popn:  2,
			want:  []int{1},
		}, {
			name:  "Pop more",
			size:  5,
			items: []int{2, 3, 1},
			popn:  6,
			want:  []int{},
		}, {
			name:  "Empty",
			size:  32,
			popn:  3,
			items: []int{},
			want:  []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRingSlice[int](tt.size, nil)
			r.AppendTail(tt.items...)
			r.RemoveHead(tt.popn)
			check(t, r, tt.want)
		})
	}
}

func TestRingSlice_NegativeGet(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		index int
		want  int
	}{
		{
			name:  "Negative Get 1",
			items: []int{1, 2, 3},
			index: -1,
			want:  3,
		}, {
			name:  "Negative Get 2",
			items: []int{1, 2, 3},
			index: -2,
			want:  2,
		}, {
			name:  "Negative Get 3",
			items: []int{1, 2, 3},
			index: -3,
			want:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRingSlice[int](len(tt.items), nil)
			r.AppendTail(tt.items...)
			g := r.Get(tt.index)
			if g != tt.want {
				t.Errorf("got data %d in pos %d, expecting %d", g, tt.index, tt.want)
			}
		})
	}
}

func TestRingSlice_NegativeSet(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		index int
		want  []int
	}{
		{
			name:  "Negative Set 1",
			items: []int{1, 2, 3},
			index: -1,
			want:  []int{1, 2, -1},
		}, {
			name:  "Negative Set 2",
			items: []int{1, 2, 3},
			index: -2,
			want:  []int{1, -2, 3},
		}, {
			name:  "Negative Set 3",
			items: []int{1, 2, 3},
			index: -3,
			want:  []int{-3, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRingSlice[int](len(tt.items), nil)
			r.AppendTail(tt.items...)
			r.Set(tt.index, tt.index)
			check(t, r, tt.want)
		})
	}
}

func TestRingSlice_RemoveTail(t *testing.T) {
	tests := []struct {
		name  string
		size  int
		items []int
		popn  int
		want  []int
	}{
		{
			name:  "Pop all",
			size:  1,
			items: []int{1},
			popn:  1,
			want:  []int{},
		}, {
			name:  "Pop some",
			size:  5,
			items: []int{2, 3, 1},
			popn:  2,
			want:  []int{2},
		}, {
			name:  "Pop more",
			size:  5,
			items: []int{2, 3, 1},
			popn:  6,
			want:  []int{},
		}, {
			name:  "Empty",
			size:  32,
			popn:  3,
			items: []int{},
			want:  []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRingSlice[int](tt.size, nil)
			r.AppendTail(tt.items...)
			r.RemoveTail(tt.popn)
			check(t, r, tt.want)
		})
	}
}

func TestRingSlice_Dummy1(t *testing.T) {
	t.Run("dummy1", func(t *testing.T) {
		r := NewRingSlice[int](5, nil)
		r.AppendTail([]int{1, 2, 3, 4, 5}...)
		check(t, r, []int{1, 2, 3, 4, 5})
		r.RemoveHead(2)
		check(t, r, []int{3, 4, 5})
		r.AppendHead(6)
		check(t, r, []int{6, 3, 4, 5})
		r.AppendHead(7)
		check(t, r, []int{7, 6, 3, 4, 5})
		r.RemoveHead(2)
		check(t, r, []int{3, 4, 5})
		r.AppendTail([]int{8, 9}...)
		check(t, r, []int{3, 4, 5, 8, 9})
		r.RemoveHead(5)
		check(t, r, []int{})
		r.AppendTail([]int{10, 11, 12, 13, 14}...)
		check(t, r, []int{10, 11, 12, 13, 14})
		r.RemoveTail(5)
		check(t, r, []int{})
		r.AppendHead(15)
		r.AppendHead(16)
		r.AppendHead(17)
		check(t, r, []int{17, 16, 15})
	})
}

func check(t *testing.T, r RingSlice[int], wanted []int) {
	for i, d := range wanted {
		g := r.Get(i)
		if g != d {
			t.Errorf("got data %d in pos %d, expecting %d", g, i, d)
		}
	}
	if r.Len() != len(wanted) {
		t.Errorf("got len %d, expecting %d", r.Len(), wanted)
	}
}
