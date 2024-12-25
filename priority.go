package behaviortree

// Priority represents a composite node in the behavior tree that attempts to execute its child nodes in priority order.
// If a child node fails, the Priority node moves to the next child. If a child succeeds, the Priority node succeeds.
type Priority[T any] struct {
	ControlNode Node[T]   // The control node managing this Priority node.
	Nodes       []Node[T] // The list of child nodes to execute in priority order.
	ActualTask  int       // The index of the currently executing child node.
	Object      T         // The object shared across nodes during execution.
}

// NewPriority creates a new Priority node with the specified child nodes.
func NewPriority[T any](nodes []Node[T]) *Priority[T] {
	return &Priority[T]{
		Nodes: nodes,
	}
}

// SetControl sets the control node for the Priority node.
func (p *Priority[T]) SetControl(control Node[T]) {
	p.ControlNode = control
}

// Start initializes the Priority node with the provided object.
func (p *Priority[T]) Start(object T) {
	p.Object = object
	p.ActualTask = 0
}

// Run executes the currently active child node. If a child node fails, the Priority node moves to the next child.
func (p *Priority[T]) Run(object T) {
	if p.ActualTask < len(p.Nodes) {
		currentNode := p.Nodes[p.ActualTask]
		currentNode.SetControl(p)
		currentNode.Start(object)
		currentNode.Run(object)
	}
}

// Success is called when a child node succeeds. It signals success to the control node.
func (p *Priority[T]) Success() {
	if p.ControlNode != nil {
		p.ControlNode.Success()
	}
}

// Fail is called when a child node fails. It advances to the next child node or signals failure to the control node
// if all children have been attempted.
func (p *Priority[T]) Fail() {
	p.ActualTask++

	if p.ActualTask < len(p.Nodes) {
		p.Run(p.Object)
	} else {
		if p.ControlNode != nil {
			p.ControlNode.Fail()
		}
	}
}

// Finish is a placeholder method for when the Priority node finishes execution.
func (p *Priority[T]) Finish(object T) {
}

// Running signals that the Priority node is still in progress to the control node.
func (p *Priority[T]) Running() {
	if p.ControlNode != nil {
		p.ControlNode.Running()
	}
}