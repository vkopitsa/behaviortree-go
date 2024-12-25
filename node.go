package behaviortree

// Node defines the interface for behavior tree nodes. All nodes in a behavior tree must implement this interface.
type Node[T any] interface {
	// Start initializes the node with the given object. Called when the node begins execution.
	Start(object T)

	// Finish finalizes the node with the given object. Called when the node finishes execution.
	Finish(object T)

	// Run executes the node's logic with the given object.
	Run(object T)

	// SetControl sets the parent or controlling node, enabling communication between child and parent nodes.
	SetControl(control Node[T])

	// Running signals that the node is still in progress and not yet finished.
	Running()

	// Success signals that the node has successfully completed its execution.
	Success()

	// Fail signals that the node has failed during its execution.
	Fail()
}

// BaseNode provides a default implementation of the Node interface.
// It can be embedded in custom node types to simplify their development.
type BaseNode[T any] struct {
	ControlNode Node[T] // The parent or controlling node managing this node.
	Object      T       // The object passed during the node's execution.
}

// SetControl sets the control node for the current node.
// The control node is typically the parent node in the behavior tree.
func (n *BaseNode[T]) SetControl(control Node[T]) {
	n.ControlNode = control
}

// setObject assigns the given object to the node.
func (n *BaseNode[T]) setObject(object T) {
	n.Object = object
}

// Start initializes the node with the given object.
func (n *BaseNode[T]) Start(object T) {
}

// Finish finalizes the node with the given object.
func (n *BaseNode[T]) Finish(object T) {
}

// Run executes the node's logic with the given object.
func (n *BaseNode[T]) Run(object T) {
}

// Running signals that the node is still in progress and not yet finished.
func (n *BaseNode[T]) Running() {
	if n.ControlNode != nil {
		n.ControlNode.Running()
	}
}

// Success signals that the node has successfully completed its execution.
func (n *BaseNode[T]) Success() {
	if n.ControlNode != nil {
		n.ControlNode.Success()
	}
}

// Fail signals that the node has failed during its execution.
func (n *BaseNode[T]) Fail() {
	if n.ControlNode != nil {
		n.ControlNode.Fail()
	}
}