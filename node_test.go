package behaviortree

import (
	"testing"
)

func TestBaseNode_SetControl(t *testing.T) {
	node := &BaseNode[int]{}
	control := &MockNode[int]{}

	node.SetControl(control)

	if node.ControlNode != control {
		t.Error("Expected ControlNode to be set correctly")
	}
}

func TestBaseNode_Running(t *testing.T) {
	node := &BaseNode[int]{}
	control := &MockNode[int]{}
	node.SetControl(control)

	node.Running()

	if !control.RunningCalled {
		t.Error("Expected Running to be called on the control")
	}
}

func TestBaseNode_Success(t *testing.T) {
	node := &BaseNode[int]{}
	control := &MockNode[int]{}
	node.SetControl(control)

	node.Success()

	if !control.SuccessCalled {
		t.Error("Expected Success to be called on the control")
	}
}

func TestBaseNode_Fail(t *testing.T) {
	node := &BaseNode[int]{}
	control := &MockNode[int]{}
	node.SetControl(control)

	node.Fail()

	if !control.FailCalled {
		t.Error("Expected Fail to be called on the control")
	}
}

func TestBaseNode_SetObject(t *testing.T) {
	node := &BaseNode[int]{}
	object := 42

	node.setObject(object)

	if node.Object != object {
		t.Errorf("Expected object to be %d, but got %d", object, node.Object)
	}
}

func TestControl_Running(t *testing.T) {
	control := &MockNode[int]{}
	node := &MockNode[int]{}
	node.SetControl(control)

	node.Running()

	if !control.RunningCalled {
		t.Error("Expected Running to be called on the Control")
	}
}

func TestControl_Success(t *testing.T) {
	control := &MockNode[int]{}
	node := &MockNode[int]{}
	node.SetControl(control)

	node.Success()

	if !control.SuccessCalled {
		t.Error("Expected Success to be called on the Control")
	}
}

func TestControl_Fail(t *testing.T) {
	control := &MockNode[int]{}
	node := &MockNode[int]{}
	node.SetControl(control)

	node.Fail()

	if !control.FailCalled {
		t.Error("Expected Fail to be called on the Control")
	}
}
