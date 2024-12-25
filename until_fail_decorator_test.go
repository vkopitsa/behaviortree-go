package behaviortree

import "testing"

func TestUntilFailDecorator_ChildKeepsRunningUntilFail(t *testing.T) {
	callCount := 0
	mockNode := &MockNode[int]{
		CustomRun: func(m *MockNode[int], obj int) {
			callCount++
			if callCount < 3 {
				if m.Control != nil {
					m.Control.Success() // Simulate success for the first two runs
				}
			} else {
				if m.Control != nil {
					m.Control.Fail() // Simulate failure on the third run
				}
			}
		},
	}
	untilFailDecorator := NewUntilFailDecorator(mockNode)

	mockControl := &MockNode[int]{}
	untilFailDecorator.SetControl(mockControl)

	untilFailDecorator.Start(42)
	untilFailDecorator.Run(42)
	untilFailDecorator.Running()

	// Verify that the child node runs multiple times
	if callCount != 3 {
		t.Errorf("Expected child node to run 3 times, but ran %d times", callCount)
	}

	// Verify that failure stops the loop and propagates success to the parent
	if !mockNode.FinishCalled {
		t.Error("Expected child node's Finish to be called after failure")
	}
	if !mockControl.SuccessCalled {
		t.Error("Expected UntilFailDecorator to propagate success after child node fails")
	}
}

func TestUntilFailDecorator_ChildFailsImmediately(t *testing.T) {
	mockNode := &MockNode[int]{
		CustomRun: func(m *MockNode[int], obj int) {
			if m.Control != nil {
				m.Control.Fail() // Simulate immediate failure
			}
		},
	}
	untilFailDecorator := NewUntilFailDecorator(mockNode)

	mockControl := &MockNode[int]{}
	untilFailDecorator.SetControl(mockControl)

	untilFailDecorator.Start(42)
	untilFailDecorator.Run(42)

	// Verify that the child node runs only once
	if !mockNode.RunCalled {
		t.Error("Expected child node's Run to be called once")
	}

	// Verify that failure stops the loop and propagates success to the parent
	if !mockNode.FinishCalled {
		t.Error("Expected child node's Finish to be called after failure")
	}
	if !mockControl.SuccessCalled {
		t.Error("Expected UntilFailDecorator to propagate success after child node fails")
	}
}
