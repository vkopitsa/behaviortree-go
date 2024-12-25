package behaviortree

import (
	"testing"
)

func TestSequence_AllNodesSucceed(t *testing.T) {
	node1 := &MockNode[int]{}
	node2 := &MockNode[int]{}
	node3 := &MockNode[int]{}
	nodes := []Node[int]{node1, node2, node3}

	sequence := NewSequence(nodes)
	sequence.Start(42)
	sequence.Run(42)

	// Verify all nodes ran successfully
	if !node1.RunCalled || !node2.RunCalled || !node3.RunCalled {
		t.Error("Expected all nodes to be executed")
	}

	if sequence.ActualTask != 3 {
		t.Errorf("Expected ActualTask to be 3, but got %d", sequence.ActualTask)
	}
}

func TestSequence_NodeFails(t *testing.T) {
	node1 := &MockNode[int]{}
	node2 := &MockNode[int]{
		CustomRun: func(m *MockNode[int], obj int) {
			m.RunCalled = true
			if m.Control != nil {
				m.Control.Fail()
			}
		},
	}
	node3 := &MockNode[int]{}

	nodes := []Node[int]{node1, node2, node3}
	sequence := NewSequence(nodes)

	sequence.Start(42)
	sequence.Run(42)

	// Verify only node1 and node2 ran, and sequence failed
	if !node1.RunCalled || !node2.RunCalled || node3.RunCalled {
		t.Error("Expected only the first two nodes to be executed")
	}

	if sequence.ActualTask != 1 {
		t.Errorf("Expected ActualTask to be 1, but got %d", sequence.ActualTask)
	}
}

func TestSequence_FailOnEmptyNodes(t *testing.T) {
	sequence := NewSequence[int](nil)

	sequence.Start(42)
	sequence.Run(42)

	// Verify sequence immediately fails
	if sequence.ControlNode != nil {
		t.Error("Expected sequence with no nodes to immediately fail")
	}
}

func TestSequence_SuccessTriggersNextNode(t *testing.T) {
	node1 := &MockNode[int]{}
	node2 := &MockNode[int]{}
	nodes := []Node[int]{node1, node2}

	sequence := NewSequence(nodes)
	sequence.Start(42)
	sequence.Run(42)

	// Verify node2 is triggered
	if !node2.RunCalled {
		t.Error("Expected node2 to run after node1 succeeds")
	}
}

func TestSequence_FailDoesNotProceed(t *testing.T) {
	node1 := &MockNode[int]{
		CustomRun: func(m *MockNode[int], obj int) {
			m.RunCalled = true
			if m.Control != nil {
				m.Control.Fail()
			}
		},
	}
	node2 := &MockNode[int]{}
	nodes := []Node[int]{node1, node2}

	sequence := NewSequence(nodes)
	sequence.Start(42)
	sequence.Run(42)

	// Verify node2 is not triggered
	if node2.RunCalled {
		t.Error("Expected node2 to not run after node1 fails")
	}
}

func TestSequence_Running(t *testing.T) {
	// Mock control node to verify if Running is called
	controlNode := &MockNode[int]{}

	// Mock sequence with a control node
	node1 := &MockNode[int]{}
	node2 := &MockNode[int]{}
	nodes := []Node[int]{node1, node2}

	sequence := NewSequence(nodes)
	sequence.SetControl(controlNode)

	// Trigger Running
	sequence.Running()

	// Verify that the Running method on the control node was called
	if !controlNode.RunningCalled {
		t.Error("Expected Running to be called on the control node")
	}
}
