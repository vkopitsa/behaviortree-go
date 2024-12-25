package behaviortree

import "testing"

func TestDecorator_Start(t *testing.T) {
	mockNode := &MockNode[int]{}
	decorator := NewDecorator(mockNode)

	decorator.Start(42)

	// Verify Start is called on the child node
	if !mockNode.StartCalled {
		t.Error("Expected Start to be called on the child node")
	}
}

func TestDecorator_Run(t *testing.T) {
	mockNode := &MockNode[int]{}
	decorator := NewDecorator(mockNode)

	decorator.Run(42)

	// Verify Run is called on the child node
	if !mockNode.RunCalled {
		t.Error("Expected Run to be called on the child node")
	}

}

func TestDecorator_Finish(t *testing.T) {
	mockNode := &MockNode[int]{}
	decorator := NewDecorator(mockNode)

	decorator.Finish(42)

	// Verify Finish is called on the child node
	if !mockNode.FinishCalled {
		t.Error("Expected Finish to be called on the child node")
	}
}
