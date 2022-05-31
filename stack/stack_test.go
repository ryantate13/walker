package stack

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		it     string
		assert func(t *testing.T, s *Stack[int])
	}{
		{
			it: "instantiates a new empty stack",
			assert: func(t *testing.T, s *Stack[int]) {
				require.NotNil(t, s)
				require.True(t, s.Empty())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			s := New[int]()
			tt.assert(t, s)
		})
	}
}

func TestStack_Empty(t *testing.T) {
	tests := []struct {
		it     string
		setup  func(s *Stack[int])
		assert func(t *testing.T, s *Stack[int])
	}{
		{
			it:    "returns true if the stack has zero items",
			setup: func(s *Stack[int]) {},
			assert: func(t *testing.T, s *Stack[int]) {
				require.True(t, s.Empty())
			},
		},
		{
			it: "returns false if the stack has one or more items",
			setup: func(s *Stack[int]) {
				s.Push(123)
			},
			assert: func(t *testing.T, s *Stack[int]) {
				require.False(t, s.Empty())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			s := New[int]()
			tt.setup(s)
			tt.assert(t, s)
		})
	}
}

func TestStack_Peek(t *testing.T) {
	type list struct {
		i    int
		next *list
	}
	tests := []struct {
		it     string
		setup  func(s *Stack[list])
		assert func(t *testing.T, s *Stack[list])
	}{
		{
			it:    "returns a zero value if the stack is empty",
			setup: func(s *Stack[list]) {},
			assert: func(t *testing.T, s *Stack[list]) {
				require.Equal(t, list{i: 0, next: nil}, s.Peek())
			},
		},
		{
			it: "returns the item at the top of the stack if the stack is not empty",
			setup: func(s *Stack[list]) {
				head := &list{
					i: 0,
					next: &list{
						i: 1,
						next: &list{
							i:    2,
							next: nil,
						},
					},
				}
				for el := head; el != nil; el = el.next {
					s.Push(*el)
				}
			},
			assert: func(t *testing.T, s *Stack[list]) {
				require.Equal(t, 2, s.Peek().i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			s := New[list]()
			tt.setup(s)
			tt.assert(t, s)
		})
	}
}

type T int

func TestStack_Pop(t *testing.T) {
	tests := []struct {
		it     string
		setup  func(s *Stack[int])
		assert func(t *testing.T, s *Stack[int])
	}{
		{
			it:    "returns a zero value if the stack is empty",
			setup: func(s *Stack[int]) {},
			assert: func(t *testing.T, s *Stack[int]) {
				require.Equal(t, 0, s.Pop())
			},
		},
		{
			it: "pops the item off the top of the stack if the stack is not empty",
			setup: func(s *Stack[int]) {
				s.Push(123)
			},
			assert: func(t *testing.T, s *Stack[int]) {
				require.Equal(t, 123, s.Pop())
				require.True(t, s.Empty())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			s := New[int]()
			tt.setup(s)
			tt.assert(t, s)
		})
	}
}

func TestStack_Push(t *testing.T) {
	tests := []struct {
		it     string
		setup  func(s *Stack[int])
		assert func(t *testing.T, s *Stack[int])
	}{
		{
			it: "puts an item at the top of the stack",
			setup: func(s *Stack[int]) {
				for i := 1; i <= 100; i++ {
					s.Push(i)
				}
			},
			assert: func(t *testing.T, s *Stack[int]) {
				for i := 100; i > 0; i-- {
					require.Equal(t, i, s.Pop())
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			s := New[int]()
			tt.setup(s)
			tt.assert(t, s)
		})
	}
}
