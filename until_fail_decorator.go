package behaviortree

// UntilFailDecorator is a decorator node that repeatedly runs its child node until the child node signals failure.
// When the child node fails, the UntilFailDecorator signals success to its control node.
type UntilFailDecorator[T any] struct {
	Decorator[T] // Embeds the Decorator structure to wrap a single child node.
	NodeRunning bool // Indicates whether the child node is currently running.
}

// NewUntilFailDecorator creates a new UntilFailDecorator with the specified child node.
// The decorator will repeatedly execute the child node until it fails.
func NewUntilFailDecorator[T any](node Node[T]) *UntilFailDecorator[T] {
	decorator := &UntilFailDecorator[T]{}
	decorator.Node = node
	decorator.Node.SetControl(decorator)
	return decorator
}

// Start initializes the UntilFailDecorator and its child node with the provided object.
func (d *UntilFailDecorator[T]) Start(object T) {
	d.setObject(object)
	d.NodeRunning = false
	d.Node.SetControl(d)
}

// Run executes the child node. If the child node is not already running, it is started first.
func (d *UntilFailDecorator[T]) Run(object T) {
	if !d.NodeRunning {
		d.Node.Start(object)
	}
	d.Node.Run(object)
}

// Running signals that the child node is still running. It notifies the control node, if present.
func (d *UntilFailDecorator[T]) Running() {
	d.NodeRunning = true
	if d.ControlNode != nil {
		d.ControlNode.Running()
	}
}

// Success is called when the child node succeeds. It restarts the child node for another execution cycle.
func (d *UntilFailDecorator[T]) Success() {
	d.NodeRunning = false
	d.Node.Start(d.Object)
	d.Node.Run(d.Object)
}

// Fail is called when the child node fails. It signals success to the control node, indicating that
// the UntilFailDecorator has completed its operation.
func (d *UntilFailDecorator[T]) Fail() {
	d.NodeRunning = false
	d.Node.Finish(d.Object)
	if d.ControlNode != nil {
		d.ControlNode.Success()
	}
}