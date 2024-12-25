package behaviortree

import (
	"testing"
)

func TestBranchNode_Start(t *testing.T) {
	nodes := []Node[int]{
		&MockNode[int]{},
		&MockNode[int]{},
	}
	branchNode := NewBranchNode(nodes)

	branchNode.Start(42)

	if branchNode.Object != 42 {
		t.Errorf("Expected object to be 42, but got %d", branchNode.Object)
	}

	if branchNode.ActualTask != 0 {
		t.Errorf("Expected ActualTask to be 0, but got %d", branchNode.ActualTask)
	}
}

func TestBranchNode_Run(t *testing.T) {
	mockNode1 := &MockNode[int]{}
	mockNode2 := &MockNode[int]{}
	nodes := []Node[int]{mockNode1, mockNode2}
	branchNode := NewBranchNode(nodes)

	branchNode.Start(42)
	branchNode.Run(42)

	// Verify the first node was started and run
	if !mockNode1.StartCalled {
		t.Error("Expected Start to be called on the first node")
	}
	if !mockNode1.RunCalled {
		t.Error("Expected Run to be called on the first node")
	}

	// Verify the second node was not run yet
	if mockNode2.StartCalled || mockNode2.RunCalled {
		t.Error("Expected second node to not be run")
	}
}

func TestBranchNode_Success(t *testing.T) {
	mockNode := &MockNode[int]{}
	nodes := []Node[int]{mockNode}
	branchNode := NewBranchNode(nodes)

	branchNode.Start(42)
	branchNode.Run(42)
	branchNode.Success()

	// Verify Success behavior
	if !mockNode.FinishCalled {
		t.Error("Expected Finish to be called on the node")
	}

	if branchNode.NodeRunning {
		t.Error("Expected NodeRunning to be false after Success")
	}

	if branchNode.Node != nil {
		t.Error("Expected Node to be nil after Success")
	}
}

func TestBranchNode_Fail(t *testing.T) {
	mockNode := &MockNode[int]{}
	nodes := []Node[int]{mockNode}
	branchNode := NewBranchNode(nodes)

	branchNode.Start(42)
	branchNode.Run(42)
	branchNode.Fail()

	// Verify Fail behavior
	if !mockNode.FinishCalled {
		t.Error("Expected Finish to be called on the node")
	}

	if branchNode.NodeRunning {
		t.Error("Expected NodeRunning to be false after Fail")
	}

	if branchNode.Node != nil {
		t.Error("Expected Node to be nil after Fail")
	}
}

func TestBranchNode_Running(t *testing.T) {
	mockNode := &MockNode[int]{}
	nodes := []Node[int]{mockNode}
	branchNode := NewBranchNode(nodes)

	control := &MockNode[int]{}
	branchNode.SetControl(control)

	branchNode.Start(42)
	branchNode.Run(42)
	branchNode.Running()

	// Verify Running behavior
	if !branchNode.NodeRunning {
		t.Error("Expected NodeRunning to be true during Running")
	}

	if !control.RunningCalled {
		t.Error("Expected Running to be called on the control")
	}
}

func TestBranchNode_Fail_NodeFinishCalled(t *testing.T) {
	// MockNode to verify Finish is called
	mockNode := &MockNode[int]{}
	mockNode.CustomRun = func(m *MockNode[int], obj int) {
		m.Control.Fail() // Simulate immediate failure
	}

	nodes := []Node[int]{mockNode}
	branchNode := NewBranchNode(nodes)

	// Start and run the branch node
	branchNode.Start(42)
	branchNode.Run(42)

	// Fail the branch node
	branchNode.Fail()

	// Verify that Finish was called on the node
	if !mockNode.FinishCalled {
		t.Error("Expected Finish to be called on the node during Fail")
	}

	// Verify that the branch node state is reset
	if branchNode.NodeRunning {
		t.Error("Expected NodeRunning to be false after Fail")
	}

	if branchNode.Node != nil {
		t.Error("Expected Node to be nil after Fail")
	}
}
