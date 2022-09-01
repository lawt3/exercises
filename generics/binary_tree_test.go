package generics

import (
	"strings"
	"testing"
)

func TestIntBinaryTree(t *testing.T) {
	tree := NewBinaryTree(compareInt)
	tree.Add(10)
	tree.Add(20)
	tree.Add(30)
	tree.Add(5)
	tree.Add(15)

	if !tree.Contains(30) {
		t.Error("Couldn't find 30 in the tree")
	}
	if !tree.Contains(15) {
		t.Error("Couldn't find 15 in the tree")
	}

	flat := tree.Flatten()
	if len(flat) != 5 {
		t.Fatal("Didn't flatten entire tree", len(flat))
	}

	expected := []int{5, 10, 15, 20, 30}
	for i, v := range flat {
		if v != expected[i] {
			t.Errorf("got %d but expected %d", v, expected[i])
		}
	}
}

func TestStructBinaryTree(t *testing.T) {
	tree := NewBinaryTree(comparePeople)
	tree.Add(Person{"Bob", 30})
	tree.Add(Person{"Maria", 25})
	tree.Add(Person{"Bob", 50})

	if !tree.Contains(Person{"Bob", 50}) {
		t.Error("Couldn't find {Bob, 50} in the tree")
	}
	if tree.Contains(Person{"Fred", 25}) {
		t.Error("Found non-existent entry in the tree")
	}

	expected := []Person{
		{"Bob", 30},
		{"Bob", 50},
		{"Maria", 25},
	}
	for i, v := range tree.Flatten() {
		if v != expected[i] {
			t.Errorf("Got %v but expected %v", v, expected[i])
		}
	}
}

func compareInt(t1, t2 int) int {
	if t1 < t2 {
		return -1
	}
	if t1 > t2 {
		return 1
	}
	return 0
}

type Person struct {
	Name string
	Age  int
}

func comparePeople(p1, p2 Person) int {
	out := strings.Compare(p1.Name, p2.Name)
	if out == 0 {
		out = p1.Age - p2.Age
	}
	return out
}
