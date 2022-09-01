package generics

type OrderableFunc[T any] func(t1, t2 T) int

type BinaryTree[T any] struct {
	f    OrderableFunc[T]
	root *BinaryNode[T]
}

type BinaryNode[T any] struct {
	left, right *BinaryNode[T]
	val         T
}

func NewBinaryTree[T any](f OrderableFunc[T]) *BinaryTree[T] {
	return &BinaryTree[T]{
		f: f,
	}
}

func (t *BinaryTree[T]) Add(v T) {
	t.root = t.root.Add(t.f, v)
}

func (t *BinaryTree[T]) Contains(v T) bool {
	return t.root.Contains(t.f, v)
}

func (t *BinaryTree[T]) Flatten() []T {
	var out []T // don't know how big the tree is ahead of time
	return t.root.Flatten(out)
}

func (n *BinaryNode[T]) Add(f OrderableFunc[T], v T) *BinaryNode[T] {
	if n == nil {
		return &BinaryNode[T]{val: v}
	}
	switch r := f(v, n.val); {
	case r <= -1:
		n.left = n.left.Add(f, v)
	case r >= 1:
		n.right = n.right.Add(f, v)
		// Don't do anything if the value is already in the tree
	}
	return n
}

func (n *BinaryNode[T]) Contains(f OrderableFunc[T], v T) bool {
	if n == nil {
		return false
	}
	switch r := f(v, n.val); {
	case r <= -1:
		return n.left.Contains(f, v)
	case r >= 1:
		return n.right.Contains(f, v)
	}
	return true
}

func (n *BinaryNode[T]) Flatten(out []T) []T {
	if n == nil {
		return out
	}
	if n.left != nil {
		out = n.left.Flatten(out)
	}
	out = append(out, n.val)
	if n.right != nil {
		out = n.right.Flatten(out)
	}
	return out
}
