package behaviortree

import "testing"

func TestPriority_AllFail(t *testing.T) {
	node1 := &MockNode[int]{
		CustomRun: func(m *MockNode[int], bj int) {
			m.Fail()
		},
	}
	node2 := &MockNode[int]{
		CustomRun: func(m *MockNode[int], bj int) {
			m.Fail()
		},
	}
	node3 := &MockNode[int]{
		CustomRun: func(m *MockNode[int], bj int) {
			m.Fail()
		},
	}

	nodes := []Node[int]{node1, node2, node3}
	priority := NewPriority(nodes)

	control := &MockNode[int]{}
	priority.SetControl(control)

	priority.Start(42)
	priority.Run(42)

	// Verify all nodes ran
	if !node1.RunCalled || !node2.RunCalled || !node3.RunCalled {
		t.Error("Expected all nodes to be executed")
	}

	// Verify the priority node fails after all child nodes fail
	if control.SuccessCalled {
		t.Error("Expected priority to propagate failure when all nodes fail")
	}
}

func TestPriority_FirstSucceeds(t *testing.T) {
	node1 := &MockNode[int]{
		CustomRun: func(m *MockNode[int], bj int) {
			m.RunCalled = true
			if m.Control != nil {
				m.Control.Success() // Trigger success for node1
			}
		},
	}
	node2 := &MockNode[int]{}
	node3 := &MockNode[int]{}

	nodes := []Node[int]{node1, node2, node3}
	priority := NewPriority(nodes)

	// Set a mock control for the priority node
	control := &MockNode[int]{}
	priority.SetControl(control)

	priority.Start(42)
	priority.Run(42)

	priority.Success()

	// Verify only the first node runs and succeeds
	if !node1.RunCalled {
		t.Error("Expected node1 to be executed")
	}
	if node2.RunCalled || node3.RunCalled {
		t.Error("Expected node2 and node3 to not be executed")
	}

	// Verify success is propagated
	if !control.SuccessCalled {
		t.Error("Expected priority to propagate success after first node succeeds")
	}
}

func TestPriority_SecondSucceeds(t *testing.T) {
	node1 := &MockNode[int]{}
	node2 := &MockNode[int]{
		CustomRun: func(m *MockNode[int], obj int) {
			m.RunCalled = true
			if m.Control != nil {
				m.Control.Success()
			}
		},
	}
	node3 := &MockNode[int]{}

	nodes := []Node[int]{node1, node2, node3}
	priority := NewPriority(nodes)

	control := &MockNode[int]{}
	priority.SetControl(control)

	priority.Start(42)
	priority.Run(42)

	// Simulate failure in the first node
	priority.Fail()
	priority.Success()

	// Verify second node succeeds
	if !node1.RunCalled {
		t.Error("Expected node1 to be executed")
	}
	if !node2.RunCalled {
		t.Error("Expected node2 to be executed")
	}
	if node3.RunCalled {
		t.Error("Expected node3 to not be executed")
	}

	// Verify success is propagated
	if !control.SuccessCalled {
		t.Error("Expected priority to propagate success after second node succeeds")
	}
}

func TestPriority_Running(t *testing.T) {
	node1 := &MockNode[int]{
		CustomRun: func(m *MockNode[int], obj int) {
			m.RunCalled = true
			if m.Control != nil {
				m.Control.Running() // Simulate running state for node1
			}
		},
	}
	node2 := &MockNode[int]{}
	node3 := &MockNode[int]{}

	nodes := []Node[int]{node1, node2, node3}
	priority := NewPriority(nodes)

	control := &MockNode[int]{}
	priority.SetControl(control)

	priority.Start(42)
	priority.Run(42)

	// Verify node1 is in a running state
	if !node1.RunCalled {
		t.Error("Expected node1 to be executed")
	}
	if !control.RunningCalled {
		t.Error("Expected priority to propagate running state")
	}

	// Verify subsequent nodes do not run while node1 is running
	if node2.RunCalled || node3.RunCalled {
		t.Error("Expected node2 and node3 to not be executed while node1 is running")
	}
}
