package behaviortree

// BehaviorTree represents the root of a behavior tree. It manages the root node and handles
// execution flow, including starting, running, and finishing the tree.
type BehaviorTree[T any] struct {
	ControlNode Node[T] // The control node managing this behavior tree.
	RootNode    Node[T] // The root node of the behavior tree.
	Started     bool    // Indicates whether the behavior tree is currently running.
	Object      T       // The object shared across nodes during execution.
}

// NewBehaviorTree creates a new BehaviorTree with the specified root node.
func NewBehaviorTree[T any](rootNode Node[T]) *BehaviorTree[T] {
	return &BehaviorTree[T]{
		RootNode: rootNode,
	}
}

// SetControl sets the control node for the behavior tree.
func (bt *BehaviorTree[T]) SetControl(control Node[T]) {
	bt.ControlNode = control
}

// SetObject assigns the given object to the behavior tree. This object is shared
// across all nodes during execution.
func (bt *BehaviorTree[T]) SetObject(object T) {
	bt.Object = object
}

// Start initializes the behavior tree. Currently unused in this implementation.
func (bt *BehaviorTree[T]) Start(object T) {
	// not used in this implementation
}

// Finish finalizes the behavior tree. Currently unused in this implementation.
func (bt *BehaviorTree[T]) Finish(object T) {
	// not used in this implementation
}

// Run executes the root node of the behavior tree with the provided object.
func (bt *BehaviorTree[T]) Run(object T) {
	bt.Object = object
	bt.RootNode.SetControl(bt)
	bt.RootNode.Start(bt.Object)
	bt.RootNode.Run(bt.Object)
}

// Running signals that the behavior tree is still in progress. It notifies the control node, if present.
func (bt *BehaviorTree[T]) Running() {
	if bt.ControlNode != nil {
		bt.ControlNode.Running()
	}
	bt.Started = true
}

// Success is called when the root node succeeds. It signals success to the control node and finalizes the tree.
func (bt *BehaviorTree[T]) Success() {
	bt.RootNode.Finish(bt.Object)
	bt.Started = false
	if bt.ControlNode != nil {
		bt.ControlNode.Success()
	}
}

// Fail is called when the root node fails. It signals failure to the control node and finalizes the tree.
func (bt *BehaviorTree[T]) Fail() {
	bt.RootNode.Finish(bt.Object)
	bt.Started = false
	if bt.ControlNode != nil {
		bt.ControlNode.Fail()
	}
}
