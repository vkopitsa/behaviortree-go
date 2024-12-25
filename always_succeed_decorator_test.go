package behaviortree

import "testing"

func TestAlwaysSucceedDecorator_ChildSuccess(t *testing.T) {
	mockNode := &MockNode[int]{
		CustomRun: func(m *MockNode[int], obj int) {
			if m.Control != nil {
				m.Control.Success() // Simulate success
			}
		},
	}
	alwaysSucceedDecorator := NewAlwaysSucceedDecorator(mockNode)

	mockControl := &MockNode[int]{}
	alwaysSucceedDecorator.SetControl(mockControl)

	alwaysSucceedDecorator.Run(42)
	alwaysSucceedDecorator.Success()

	// Verify that the child node ran
	if !mockNode.RunCalled {
		t.Error("Expected child node's Run to be called")
	}

	// Verify that success was propagated as success
	if !mockControl.SuccessCalled {
		t.Error("Expected AlwaysSucceedDecorator to propagate success")
	}
}

func TestAlwaysSucceedDecorator_ChildFail(t *testing.T) {
	mockNode := &MockNode[int]{
		CustomRun: func(m *MockNode[int], obj int) {
			if m.Control != nil {
				m.Control.Fail() // Simulate failure
			}
		},
	}
	alwaysSucceedDecorator := NewAlwaysSucceedDecorator(mockNode)

	mockControl := &MockNode[int]{}
	alwaysSucceedDecorator.SetControl(mockControl)

	alwaysSucceedDecorator.Run(42)
	alwaysSucceedDecorator.Fail()

	// Verify that the child node ran
	if !mockNode.RunCalled {
		t.Error("Expected child node's Run to be called")
	}

	// Verify that failure was overridden as success
	if !mockControl.SuccessCalled {
		t.Error("Expected AlwaysSucceedDecorator to override failure to success")
	}
}
