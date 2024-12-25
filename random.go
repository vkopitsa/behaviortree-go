package behaviortree

import (
	"math/rand"
	"time"
)

// Random represents a composite node that selects one child node at random to execute.
// It ensures that one of its child nodes is chosen and run each time it starts.
type Random[T any] struct {
	BranchNode[T] // Embeds the BranchNode structure to manage child nodes.
}

// NewRandom creates a new Random node with the specified child nodes.
// Each time the node starts, it selects a child node at random.
func NewRandom[T any](nodes []Node[T]) *Random[T] {
	rand.Seed(time.Now().UnixNano())
	bn := NewBranchNode(nodes)
	return &Random[T]{
		BranchNode: *bn,
	}
}

// Start initializes the Random node and selects a random child node to execute.
func (r *Random[T]) Start(object T) {
	r.BranchNode.Start(object)
	if len(r.Nodes) > 0 {
		r.ActualTask = rand.Intn(len(r.Nodes)) // Select a random child node
		r.Nodes[r.ActualTask].Start(object)   // Start the selected child node
	}
}

// Success is called when the selected child node succeeds.
func (r *Random[T]) Success() {
	r.BaseNode.Success()
	if r.ControlNode != nil {
		r.ControlNode.Success()
	}
}

// Fail is called when the selected child node fails.
func (r *Random[T]) Fail() {
	r.BaseNode.Fail()
	if r.ControlNode != nil {
		r.ControlNode.Fail()
	}
}