package behaviortree

// InvertDecorator is a decorator node that inverts the result of its child node.
// When the child node succeeds, the decorator signals failure, and when the child node fails, the decorator signals success.
type InvertDecorator[T any] struct {
	Decorator[T] // Embeds the Decorator structure to wrap a single child node.
}

// NewInvertDecorator creates a new InvertDecorator with the specified child node.
// The child node's behavior is inverted by this decorator.
func NewInvertDecorator[T any](node Node[T]) *InvertDecorator[T] {
	decorator := &InvertDecorator[T]{}
	decorator.Node = node
	decorator.Node.SetControl(decorator)
	return decorator
}

// Success is called when the child node succeeds. It signals failure to the control node.
func (d *InvertDecorator[T]) Success() {
	if d.ControlNode != nil {
		d.ControlNode.Fail()
	}
}

// Fail is called when the child node fails. It signals success to the control node.
func (d *InvertDecorator[T]) Fail() {
	if d.ControlNode != nil {
		d.ControlNode.Success()
	}
}