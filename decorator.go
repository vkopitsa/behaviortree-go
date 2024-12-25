package behaviortree

// Decorator is a node in the behavior tree that wraps and modifies the behavior of a single child node.
// It delegates its lifecycle methods (Start, Run, Finish) to the child node, allowing for additional behavior
// to be added before or after the child's execution.
type Decorator[T any] struct {
	BaseNode[T] // Embeds common behavior for behavior tree nodes.
	Node        Node[T] // The child node whose behavior is being decorated.
}

// NewDecorator creates a new Decorator node with the specified child node.
func NewDecorator[T any](node Node[T]) *Decorator[T] {
	return &Decorator[T]{
		Node: node,
	}
}

// Start initializes the decorator and its child node with the provided object.
// This method allows the decorator to perform any setup logic before starting the child.
func (d *Decorator[T]) Start(object T) {
	d.Node.Start(object)
}

// Finish finalizes the decorator and its child node with the provided object.
// This method allows the decorator to perform any cleanup logic after the child has finished.
func (d *Decorator[T]) Finish(object T) {
	d.Node.Finish(object)
}

// Run executes the child node's Run method. This method can be overridden to add custom logic
// before or after delegating execution to the child.
func (d *Decorator[T]) Run(object T) {
	d.Node.Run(object)
}
