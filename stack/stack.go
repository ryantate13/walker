package stack

// A Stack is a generic LIFO queue
type Stack[T any] struct {
	entries []T
}

// New instantiates a new stack of a given type
func New[T any]() *Stack[T] {
	return &Stack[T]{
		entries: make([]T, 0),
	}
}

// Empty returns true if there are no items in the stack
func (s *Stack[T]) Empty() bool {
	return len(s.entries) == 0
}

// Push adds an item to the stack
func (s *Stack[T]) Push(el T) {
	s.entries = append(s.entries, el)
}

func (s *Stack[T]) zero() T {
	var zeroVal T
	return zeroVal
}

// Peek returns the item at the top of the stack or the corresponding zero value if the stack is empty
func (s *Stack[T]) Peek() T {
	if len(s.entries) == 0 {
		return s.zero()
	}
	return s.entries[len(s.entries)-1]
}

// Pop removes and returns the item at the top of the stack or the corresponding zero value if the stack is empty.
func (s *Stack[T]) Pop() T {
	i := len(s.entries) - 1
	if i == -1 {
		return s.zero()
	}
	el := s.entries[i]
	s.entries = s.entries[:i]
	return el
}
