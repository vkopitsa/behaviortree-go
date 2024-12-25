package behaviortree

// AlwaysFailDecorator is a decorator node that forces the result of its child node to always be failure.
// Regardless of whether the child node succeeds or fails, this decorator signals failure to its control node.
type AlwaysFailDecorator[T any] struct {
	Decorator[T] // Embeds the Decorator structure to wrap a single child node.
}

// NewAlwaysFailDecorator creates a new AlwaysFailDecorator with the specified child node.
// The child node's result will always be converted to failure by this decorator.
func NewAlwaysFailDecorator[T any](node Node[T]) *AlwaysFailDecorator[T] {
	decorator := &AlwaysFailDecorator[T]{}
	decorator.Node = node
	decorator.Node.SetControl(decorator)
	return decorator
}

// Success is called when the child node succeeds. Instead of signaling success, it signals failure to the control node.
func (d *AlwaysFailDecorator[T]) Success() {
	if d.ControlNode != nil {
		d.ControlNode.Fail()
	}
}

// Fail is called when the child node fails. It signals failure to the control node as usual.
func (d *AlwaysFailDecorator[T]) Fail() {
	if d.ControlNode != nil {
		d.ControlNode.Fail()
	}
}