package behaviortree

// Sequence represents a composite node in the behavior tree that executes its child nodes sequentially.
// If one child fails, the sequence fails. If all children succeed, the sequence succeeds.
type Sequence[T any] struct {
	ControlNode Node[T] // The control node managing the sequence's behavior in the tree.
	Nodes       []Node[T] // The list of child nodes to execute in sequence.
	ActualTask  int // The index of the currently executing child node.
	Object      T // The object shared across nodes during execution.
}

// NewSequence creates a new Sequence node with the provided child nodes.
func NewSequence[T any](nodes []Node[T]) *Sequence[T] {
	return &Sequence[T]{
		Nodes: nodes,
	}
}

// SetControl sets the control node for the Sequence.
func (s *Sequence[T]) SetControl(control Node[T]) {
	s.ControlNode = control
}

// Start initializes the sequence and its child nodes with the provided object.
func (s *Sequence[T]) Start(object T) {
	s.Object = object
	s.ActualTask = 0
	for _, node := range s.Nodes {
		node.Start(object)
	}
}

// Run executes the current child node in the sequence. If all child nodes have been
// successfully executed, the sequence itself succeeds.
func (s *Sequence[T]) Run(object T) {
	if s.ActualTask >= len(s.Nodes) {
		s.Success()
		return
	}
	currentNode := s.Nodes[s.ActualTask]
	currentNode.SetControl(s)
	currentNode.Run(object)
}

// Success is called when a child node succeeds. It advances to the next child node
// or signals success to the control node if all children have succeeded.
func (s *Sequence[T]) Success() {
	s.ActualTask++
	if s.ActualTask < len(s.Nodes) {
		s.Run(s.Object)
	} else {
		if s.ControlNode != nil {
			s.ControlNode.Success()
		}
	}
}

// Fail is called when a child node fails. It signals failure to the control node.
func (s *Sequence[T]) Fail() {
	if s.ControlNode != nil {
		s.ControlNode.Fail()
	}
}

// Finish is a placeholder method for when the sequence finishes execution.
func (s *Sequence[T]) Finish(object T) {
	// not used in this implementation
}

// Running signals that the sequence is still in progress to the control node.
func (s *Sequence[T]) Running() {
	if s.ControlNode != nil {
		s.ControlNode.Running()
	}
}