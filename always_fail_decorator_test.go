package behaviortree

import "testing"

func TestAlwaysFailDecorator_ChildSuccess(t *testing.T) {
	mockNode := &MockNode[int]{
		CustomRun: func(m *MockNode[int], obj int) {
			if m.Control != nil {
				m.Control.Success() // Simulate success
			}
		},
	}
	alwaysFailDecorator := NewAlwaysFailDecorator(mockNode)

	mockControl := &MockNode[int]{}
	alwaysFailDecorator.SetControl(mockControl)

	alwaysFailDecorator.Run(42)
	alwaysFailDecorator.Success()

	// Verify that the child node ran
	if !mockNode.RunCalled {
		t.Error("Expected child node's Run to be called")
	}

	// Verify that success was overridden as failure
	if !mockControl.FailCalled {
		t.Error("Expected AlwaysFailDecorator to override success to failure")
	}
}

func TestAlwaysFailDecorator_ChildFail(t *testing.T) {
	mockNode := &MockNode[int]{
		CustomRun: func(m *MockNode[int], obj int) {
			if m.Control != nil {
				m.Control.Fail() // Simulate failure
			}
		},
	}
	alwaysFailDecorator := NewAlwaysFailDecorator(mockNode)

	mockControl := &MockNode[int]{}
	alwaysFailDecorator.SetControl(mockControl)

	alwaysFailDecorator.Run(42)
	alwaysFailDecorator.Success()
	alwaysFailDecorator.Fail()

	// Verify that the child node ran
	if !mockNode.RunCalled {
		t.Error("Expected child node's Run to be called")
	}

	// Verify that failure was propagated as failure
	if !mockControl.FailCalled {
		t.Error("Expected AlwaysFailDecorator to propagate failure")
	}
}
