package behaviortree

import "testing"

func TestTask_RunSuccess(t *testing.T) {
	task := NewTask(func(task *Task[int], object int) {
		if task.ControlNode != nil {
			task.ControlNode.Success() // Simulate success
		}
	})

	mockControl := &MockNode[int]{}
	task.SetControl(mockControl)

	task.Run(42)

	// Verify success was propagated
	if !mockControl.SuccessCalled {
		t.Error("Expected Task to propagate success")
	}
}

func TestTask_RunFail(t *testing.T) {
	task := NewTask(func(task *Task[int], object int) {
		if task.ControlNode != nil {
			task.ControlNode.Fail() // Simulate failure
		}
	})

	mockControl := &MockNode[int]{}
	task.SetControl(mockControl)

	task.Run(42)

	// Verify failure was propagated
	if !mockControl.FailCalled {
		t.Error("Expected Task to propagate failure")
	}
}

func TestTask_CustomLogic(t *testing.T) {
	executed := false
	task := NewTask(func(task *Task[int], object int) {
		executed = true
		if task.ControlNode != nil {
			task.ControlNode.Success() // Simulate success
		}
	})

	mockControl := &MockNode[int]{}
	task.SetControl(mockControl)

	task.Run(42)

	// Verify the custom logic was executed
	if !executed {
		t.Error("Expected Task's custom logic to be executed")
	}

	// Verify success was propagated
	if !mockControl.SuccessCalled {
		t.Error("Expected Task to propagate success")
	}
}
