package behaviortree

// AlwaysSucceedDecorator is a decorator node that forces the result of its child node to always be success.
// Regardless of whether the child node succeeds or fails, this decorator signals success to its control node.
type AlwaysSucceedDecorator[T any] struct {
	Decorator[T] // Embeds the Decorator structure to wrap a single child node.
}

// NewAlwaysSucceedDecorator creates a new AlwaysSucceedDecorator with the specified child node.
// The child node's result will always be converted to success by this decorator.
func NewAlwaysSucceedDecorator[T any](node Node[T]) *AlwaysSucceedDecorator[T] {
	decorator := &AlwaysSucceedDecorator[T]{}
	decorator.Node = node
	decorator.Node.SetControl(decorator)
	return decorator
}

// Success is called when the child node succeeds. It signals success to the control node.
func (d *AlwaysSucceedDecorator[T]) Success() {
	if d.ControlNode != nil {
		d.ControlNode.Success()
	}
}

// Fail is called when the child node fails. Instead of signaling failure, it signals success to the control node.
func (d *AlwaysSucceedDecorator[T]) Fail() {
	if d.ControlNode != nil {
		d.ControlNode.Success()
	}
}