package behaviortree

import "testing"

func TestInvertDecorator_SuccessInvertsToFail(t *testing.T) {
	mockNode := &MockNode[int]{
		CustomRun: func(m *MockNode[int], obj int) {
			if m.Control != nil {
				m.Control.Success() // Simulate success
			}
		},
	}
	invertDecorator := NewInvertDecorator(mockNode)

	mockControl := &MockNode[int]{}
	invertDecorator.SetControl(mockControl)

	invertDecorator.Run(42)
	invertDecorator.Success()

	// Verify that the child node ran
	if !mockNode.RunCalled {
		t.Error("Expected child node's Run to be called")
	}

	// Verify that success from the child node was inverted to fail
	if !mockControl.FailCalled {
		t.Error("Expected InvertDecorator to invert success to fail")
	}
}

func TestInvertDecorator_FailInvertsToSuccess(t *testing.T) {
	mockNode := &MockNode[int]{
		CustomRun: func(m *MockNode[int], obj int) {
			if m.Control != nil {
				m.Control.Fail() // Simulate failure
			}
		},
	}
	invertDecorator := NewInvertDecorator(mockNode)

	mockControl := &MockNode[int]{}
	invertDecorator.SetControl(mockControl)

	invertDecorator.Run(42)
	invertDecorator.Fail()

	// Verify that the child node ran
	if !mockNode.RunCalled {
		t.Error("Expected child node's Run to be called")
	}

	// Verify that failure from the child node was inverted to success
	if !mockControl.SuccessCalled {
		t.Error("Expected InvertDecorator to invert fail to success")
	}
}
