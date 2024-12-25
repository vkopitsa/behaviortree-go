package behaviortree

// BranchNode serves as a base for composite nodes, managing multiple child nodes and their execution flow.
// It provides the structure for running child nodes sequentially or in other specified orders.
type BranchNode[T any] struct {
	BaseNode[T]   // Embeds common behavior for behavior tree nodes.
	Nodes         []Node[T] // The list of child nodes managed by this branch node.
	Node          Node[T]   // The currently active child node.
	ActualTask    int       // Index of the currently executing child node.
	NodeRunning   bool      // Indicates whether a child node is currently running.
}

// NewBranchNode creates a new BranchNode with the specified child nodes.
func NewBranchNode[T any](nodes []Node[T]) *BranchNode[T] {
	return &BranchNode[T]{
		Nodes: nodes,
	}
}

// Start initializes the branch node with the provided object. It resets the active task index and prepares
// the node for execution.
func (b *BranchNode[T]) Start(object T) {
	if !b.NodeRunning {
		b.setObject(object)
		b.ActualTask = 0
	}
}

// Run executes the currently active child node. If there are remaining tasks, it delegates execution to `_run`.
func (b *BranchNode[T]) Run(object T) {
	if b.ActualTask < len(b.Nodes) {
		b._run(object)
	}
}

// _run handles the execution of the current child node. It starts the child node if it hasn't already started.
func (b *BranchNode[T]) _run(object T) {
	if !b.NodeRunning {
		b.Node = b.Nodes[b.ActualTask]
		b.Node.Start(object)
		b.Node.SetControl(b)
	}
	b.Node.Run(object)
}

// Running signals that the current child node is still running. It notifies the control node if one exists.
func (b *BranchNode[T]) Running() {
	b.NodeRunning = true
	if b.ControlNode != nil {
		b.ControlNode.Running()
	}
}

// Success is called when the current child node succeeds. It resets the state and prepares for the next task.
func (b *BranchNode[T]) Success() {
	b.NodeRunning = false
	if b.Node != nil {
		b.Node.Finish(b.Object)
	}
	b.Node = nil
}

// Fail is called when the current child node fails. It resets the state and notifies the control node if one exists.
func (b *BranchNode[T]) Fail() {
	b.NodeRunning = false
	if b.Node != nil {
		b.Node.Finish(b.Object)
	}
	b.Node = nil
}