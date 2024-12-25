package behaviortree

import (
	"testing"
)

func TestRandom_Start(t *testing.T) {
	node1 := &MockNode[int]{}
	node2 := &MockNode[int]{}
	node3 := &MockNode[int]{}
	nodes := []Node[int]{node1, node2, node3}

	random := NewRandom(nodes)
	random.Start(42)

	// Verify one and only one child node was selected
	selectedCount := 0
	for _, node := range nodes {
		if node.(*MockNode[int]).StartCalled {
			selectedCount++
		}
	}
	if selectedCount != 1 {
		t.Errorf("Expected only one node to be started, but got %d", selectedCount)
	}
}

func TestRandom_Success(t *testing.T) {
	node1 := &MockNode[int]{}
	node2 := &MockNode[int]{}
	node3 := &MockNode[int]{}
	nodes := []Node[int]{node1, node2, node3}

	control := &MockNode[int]{}

	random := NewRandom(nodes)
	random.SetControl(control)
	random.Start(42)
	random.Run(42)
	random.Success()

	// Verify Success was propagated correctly
	if !control.SuccessCalled {
		t.Error("Expected Success to be called on the control node")
	}
}

func TestRandom_Fail(t *testing.T) {
	node1 := &MockNode[int]{}
	node2 := &MockNode[int]{}
	node3 := &MockNode[int]{}
	nodes := []Node[int]{node1, node2, node3}

	control := &MockNode[int]{}

	random := NewRandom(nodes)
	random.SetControl(control)
	random.Start(42)
	random.Run(42)
	random.Fail()

	// Verify Fail was propagated correctly
	if !control.FailCalled {
		t.Error("Expected Fail to be called on the control node")
	}
}

func TestRandom_Randomness(t *testing.T) {
	node1 := &MockNode[int]{}
	node2 := &MockNode[int]{}
	node3 := &MockNode[int]{}
	nodes := []Node[int]{node1, node2, node3}

	random := NewRandom(nodes)

	// Collect selected nodes over multiple starts
	selectedNodes := make(map[int]int)
	for i := 0; i < 100; i++ {
		random.Start(42)
		if random.ActualTask >= 0 && random.ActualTask < len(nodes) {
			selectedNodes[random.ActualTask]++
		}

		// Reset mock nodes for the next iteration
		for _, node := range nodes {
			mock := node.(*MockNode[int])
			mock.StartCalled = false
		}
	}

	// Ensure randomness by checking all nodes are selected at least once
	for i := range nodes {
		if selectedNodes[i] == 0 {
			t.Errorf("Node %d was never selected", i)
		}
	}
}
