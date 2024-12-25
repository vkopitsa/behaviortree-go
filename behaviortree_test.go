package behaviortree

import (
	"testing"
)

type MockNode[T any] struct {
	t                *testing.T
	debug            bool
	SetControlCalled bool
	StartCalled      bool
	RunCalled        bool
	FinishCalled     bool
	RunningCalled    bool
	FailCalled       bool
	SuccessCalled    bool
	Control          Node[T]
	CustomRun        func(m *MockNode[T], obj T)
}

func NewMockNode[T any](t *testing.T) *MockNode[T] {
	return &MockNode[T]{t: t}
}

func (m *MockNode[T]) SetControl(control Node[T]) {
	m.SetControlCalled = true
	m.Control = control
	if m.t != nil && m.debug {
		m.t.Logf("SetControl called with %T", control)
	}
}

func (m *MockNode[T]) Start(obj T) {
	m.StartCalled = true
	if m.t != nil && m.debug {
		m.t.Logf("Start called with object: %v", obj)
	}
}

func (m *MockNode[T]) Run(obj T) {
	m.RunCalled = true
	if m.t != nil && m.debug {
		m.t.Logf("Run called with object: %v", obj)
	}
	if m.CustomRun != nil {
		m.CustomRun(m, obj)
	} else if m.Control != nil {
		m.Control.Success()
	}
}

func (m *MockNode[T]) Finish(obj T) {
	m.FinishCalled = true
	if m.t != nil && m.debug {
		m.t.Logf("Finish called with object: %v", obj)
	}
}

func (m *MockNode[T]) Running() {
	m.RunningCalled = true
	if m.t != nil && m.debug {
		m.t.Log("Running called")
	}
	if m.Control != nil {
		m.Control.Running()
	}
}

func (m *MockNode[T]) Success() {
	m.SuccessCalled = true
	if m.t != nil && m.debug {
		m.t.Log("Success called")
	}
	if m.Control != nil {
		m.Control.Success()
	}
}

func (m *MockNode[T]) Fail() {
	m.FailCalled = true
	if m.t != nil && m.debug {
		m.t.Log("Fail called")
	}
	if m.Control != nil {
		m.Control.Fail()
	}
}

func TestBehaviorTree_SetObject(t *testing.T) {
	root := &MockNode[int]{}
	tree := NewBehaviorTree[int](root)

	tree.SetObject(42)

	if tree.Object != 42 {
		t.Errorf("Expected object to be 42, but got %d", tree.Object)
	}
}

func TestBehaviorTree_Run(t *testing.T) {
	root := &MockNode[int]{}
	tree := NewBehaviorTree(root)

	tree.Run(42)
	tree.Running()

	// Verify root node methods were called
	if !root.SetControlCalled {
		t.Error("Expected SetControl to be called on the root node")
	}
	if !root.StartCalled {
		t.Error("Expected Start to be called on the root node")
	}
	if !root.RunCalled {
		t.Error("Expected Run to be called on the root node")
	}

	// Verify the object was set
	if tree.Object != 42 {
		t.Errorf("Expected object to be 42, but got %d", tree.Object)
	}
}

func TestBehaviorTree_Success(t *testing.T) {
	root := &MockNode[int]{}
	tree := NewBehaviorTree(root)

	// Simulate a run
	tree.Run(42)

	// Call Success
	tree.Success()

	// Verify Finish was called on the root node
	if !root.FinishCalled {
		t.Error("Expected Finish to be called on the root node")
	}

	// Verify tree's state
	if tree.Started {
		t.Error("Expected tree to be not started after Success")
	}
}

func TestBehaviorTree_Fail(t *testing.T) {
	root := &MockNode[int]{}
	tree := NewBehaviorTree(root)

	// Simulate a run
	tree.Run(42)

	// Call Fail
	tree.Fail()

	// Verify Finish was called on the root node
	if !root.FinishCalled {
		t.Error("Expected Finish to be called on the root node")
	}

	// Verify tree's state
	if tree.Started {
		t.Error("Expected tree to be not started after Fail")
	}
}

func TestControl_WithBehaviorTree(t *testing.T) {
	control := &MockNode[int]{}
	root := &MockNode[int]{}
	tree := NewBehaviorTree(root)
	tree.SetControl(control)

	tree.Run(42)
	tree.Running()

	// Ensure Running state is reported
	if !control.RunningCalled {
		t.Error("Expected Running to be called on the Control during Run")
	}

	// Simulate Success in the root node
	tree.Success()

	if !control.SuccessCalled {
		t.Error("Expected Success to be called on the Control after Success")
	}

	// Simulate Failure in the root node
	tree.Fail()

	if !control.FailCalled {
		t.Error("Expected Fail to be called on the Control after Fail")
	}
}

func TestSequenceSuccess(t *testing.T) {
	control := &MockNode[int]{}

	task1 := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})
	task2 := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})

	sequence := NewSequence[int]([]Node[int]{task1, task2})

	bt := NewBehaviorTree[int](sequence)
	bt.SetControl(control)
	bt.SetObject(0)
	bt.Run(0)

	if !control.SuccessCalled {
		t.Error("Expected Sequence to call Success on control")
	}
}

func TestSequenceFail(t *testing.T) {
	control := &MockNode[int]{}

	task1 := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})
	task2 := NewTask[int](func(task *Task[int], obj int) {
		task.Fail()
	})

	sequence := NewSequence[int]([]Node[int]{task1, task2})

	bt := NewBehaviorTree[int](sequence)
	bt.SetControl(control)
	bt.SetObject(0)
	bt.Run(1)

	if !control.FailCalled {
		t.Error("Expected Sequence to call Fail on control")
	}
}

func TestPrioritySuccess(t *testing.T) {
	control := &MockNode[int]{}

	task1 := NewTask[int](func(task *Task[int], obj int) {
		task.Fail()
	})
	task2 := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})

	priority := NewPriority[int]([]Node[int]{task1, task2})

	bt := NewBehaviorTree[int](priority)
	bt.SetControl(control)
	bt.SetObject(0)
	bt.Run(1)

	if !control.SuccessCalled {
		t.Error("Expected Priority to call Success on control")
	}
}

func TestPriorityFail(t *testing.T) {
	control := &MockNode[int]{}

	task1 := NewTask[int](func(task *Task[int], obj int) {
		task.Fail()
	})
	task2 := NewTask[int](func(task *Task[int], obj int) {
		task.Fail()
	})

	priority := NewPriority[int]([]Node[int]{task1, task2})

	bt := NewBehaviorTree[int](priority)
	bt.SetControl(control)
	bt.SetObject(0)
	bt.Run(1)

	if !control.FailCalled {
		t.Error("Expected Priority to call Fail on control")
	}
}

func TestSequenceWithSingleNodeSuccess(t *testing.T) {
	control := NewMockNode[int](t)

	task := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})

	sequence := NewSequence[int]([]Node[int]{task})

	bt := NewBehaviorTree[int](sequence)
	bt.SetControl(control)
	bt.SetObject(0)
	bt.Run(0)

	if !control.SuccessCalled {
		t.Error("Expected Sequence with single successful node to call Success on control")
	}
}

func TestSequenceWithSingleNodeFail(t *testing.T) {
	control := NewMockNode[int](t)

	task := NewTask[int](func(task *Task[int], obj int) {
		task.Fail()
	})

	sequence := NewSequence[int]([]Node[int]{task})

	bt := NewBehaviorTree[int](sequence)
	bt.SetControl(control)
	bt.SetObject(0)
	bt.Run(0)

	if !control.FailCalled {
		t.Error("Expected Sequence with single failing node to call Fail on control")
	}
}

func TestSequenceWithIntermediateFailure(t *testing.T) {
	control := NewMockNode[int](t)

	task1 := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})
	task2 := NewTask[int](func(task *Task[int], obj int) {
		task.Fail()
	})
	task3 := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})

	sequence := NewSequence[int]([]Node[int]{task1, task2, task3})

	bt := NewBehaviorTree[int](sequence)
	bt.SetControl(control)
	bt.SetObject(0)
	bt.Run(0)

	if !control.FailCalled {
		t.Error("Expected Sequence to call Fail on control when an intermediate node fails")
	}

	if task3.RunCalled {
		t.Error("Expected Sequence to not execute task3 after task2 fails")
	}
}

func TestPriorityWithFirstNodeSuccess(t *testing.T) {
	control := NewMockNode[int](t)

	task1 := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})
	task2 := NewTask[int](func(task *Task[int], obj int) {
		task.Fail() // Should not be executed
	})

	priority := NewPriority[int]([]Node[int]{task1, task2})

	bt := NewBehaviorTree[int](priority)
	bt.SetControl(control)
	bt.SetObject(0)
	bt.Run(0)

	if !control.SuccessCalled {
		t.Error("Expected Priority to call Success on control when first node succeeds")
	}

	if task2.RunCalled {
		t.Error("Expected Priority to not execute task2 after task1 succeeds")
	}
}

func TestPriorityWithMiddleNodeSuccess(t *testing.T) {
	control := NewMockNode[int](t)

	task1 := NewTask[int](func(task *Task[int], obj int) {
		task.Fail()
	})
	task2 := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})
	task3 := NewTask[int](func(task *Task[int], obj int) {
		task.Fail()
	})

	priority := NewPriority[int]([]Node[int]{task1, task2, task3})

	bt := NewBehaviorTree[int](priority)
	bt.SetControl(control)
	bt.SetObject(0)
	bt.Run(0)

	if !control.SuccessCalled {
		t.Error("Expected Priority to call Success on control when task2 succeeds")
	}

	if task3.RunCalled {
		t.Error("Expected Priority to not execute task3 after task2 succeeds")
	}
}

func TestPriorityWithAllNodesFail(t *testing.T) {
	control := NewMockNode[int](t)

	task1 := NewTask[int](func(task *Task[int], obj int) {
		task.Fail()
	})
	task2 := NewTask[int](func(task *Task[int], obj int) {
		task.Fail()
	})
	task3 := NewTask[int](func(task *Task[int], obj int) {
		task.Fail()
	})

	priority := NewPriority[int]([]Node[int]{task1, task2, task3})

	bt := NewBehaviorTree[int](priority)
	bt.SetControl(control)
	bt.SetObject(0)
	bt.Run(0)

	if !control.FailCalled {
		t.Error("Expected Priority to call Fail on control when all nodes fail")
	}
}

func TestAlwaysFailDecorator_AlwaysFails(t *testing.T) {
	control := NewMockNode[int](t)

	task1 := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})
	task2 := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})

	sequence := NewSequence[int]([]Node[int]{task1, task2})
	alwaysFailDecorator := NewAlwaysFailDecorator(sequence)

	bt := NewBehaviorTree[int](alwaysFailDecorator)
	bt.SetControl(control)
	bt.SetObject(0)
	bt.Run(0)

	if !control.FailCalled {
		t.Error("Expected AlwaysFailDecorator to call Fail on control regardless of child success")
	}
}

func TestAlwaysSucceedDecorator_AlwaysSucceeds(t *testing.T) {
	control := NewMockNode[int](t)

	task := NewTask[int](func(task *Task[int], obj int) {
		task.Fail() // Child would normally fail
	})

	alwaysSucceedDecorator := NewAlwaysSucceedDecorator(task)

	bt := NewBehaviorTree[int](alwaysSucceedDecorator)
	bt.SetControl(control)
	bt.SetObject(0)
	bt.Run(0)

	if !control.SuccessCalled {
		t.Error("Expected AlwaysSucceedDecorator to call Success on control regardless of child failure")
	}
}

func TestInvertDecorator_InvertsSuccessToFail(t *testing.T) {
	control := NewMockNode[int](t)

	task := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})

	invertDecorator := NewInvertDecorator(task)

	bt := NewBehaviorTree[int](invertDecorator)
	bt.SetControl(control)
	bt.SetObject(0)
	bt.Run(0)

	if !control.FailCalled {
		t.Error("Expected InvertDecorator to invert Success to Fail on control")
	}
}

func TestInvertDecorator_InvertsFailToSuccess(t *testing.T) {
	control := NewMockNode[int](t)

	task := NewTask[int](func(task *Task[int], obj int) {
		task.Fail()
	})

	invertDecorator := NewInvertDecorator(task)

	bt := NewBehaviorTree[int](invertDecorator)
	bt.SetControl(control)
	bt.SetObject(0)
	bt.Run(0)

	if !control.SuccessCalled {
		t.Error("Expected InvertDecorator to invert Fail to Success on control")
	}
}

func TestUntilFailDecorator_RunsUntilChildFails(t *testing.T) {
	control := NewMockNode[int](t)
	callCount := 0

	mockNode := NewMockNode[int](t)
	mockNode.CustomRun = func(m *MockNode[int], obj int) {
		callCount++
		if callCount < 3 {
			m.Control.Success() // Simulate success for the first two runs
		} else {
			m.Control.Fail() // Simulate failure on the third run
		}
	}

	untilFailDecorator := NewUntilFailDecorator[int](mockNode)

	bt := NewBehaviorTree[int](untilFailDecorator)
	bt.SetControl(control)
	bt.SetObject(0)
	bt.Run(0)

	if callCount != 3 {
		t.Errorf("Expected child node to run 3 times, but ran %d times", callCount)
	}

	if !control.SuccessCalled {
		t.Error("Expected UntilFailDecorator to call Success on control after child fails")
	}
}

func TestUntilFailDecorator_ChildFailsImmediately1(t *testing.T) {
	control := NewMockNode[int](t)

	mockNode := NewMockNode[int](t)
	mockNode.CustomRun = func(m *MockNode[int], obj int) {
		m.Control.Fail() // Simulate immediate failure
	}

	untilFailDecorator := NewUntilFailDecorator[int](mockNode)

	bt := NewBehaviorTree[int](untilFailDecorator)
	bt.SetControl(control)
	bt.SetObject(0)
	bt.Run(0)

	if mockNode.RunCalled != true {
		t.Error("Expected child node to run once")
	}

	if !control.SuccessCalled {
		t.Error("Expected UntilFailDecorator to call Success on control after immediate failure")
	}
}

func TestTask_RunSuccess1(t *testing.T) {
	control := NewMockNode[int](t)

	task := NewTask[int](func(task *Task[int], object int) {
		task.Success()
	})

	task.SetControl(control)
	task.Run(0)

	if !control.SuccessCalled {
		t.Error("Expected Task to propagate Success to control")
	}
}

func TestTask_RunFail1(t *testing.T) {
	control := NewMockNode[int](t)

	task := NewTask[int](func(task *Task[int], object int) {
		task.Fail()
	})

	task.SetControl(control)
	task.Run(0)

	if !control.FailCalled {
		t.Error("Expected Task to propagate Fail to control")
	}
}

func TestBehaviorTree_RunMultipleSequences(t *testing.T) {
	control := NewMockNode[int](t)

	sequence1 := NewSequence[int]([]Node[int]{
		NewTask[int](func(task *Task[int], obj int) { task.Success() }),
		NewTask[int](func(task *Task[int], obj int) { task.Success() }),
	})

	sequence2 := NewSequence[int]([]Node[int]{
		NewTask[int](func(task *Task[int], obj int) { task.Success() }),
		NewTask[int](func(task *Task[int], obj int) { task.Fail() }),
	})

	// First run: sequence1 should succeed
	bt := NewBehaviorTree[int](sequence1)
	bt.SetControl(control)
	bt.SetObject(1)
	bt.Run(1)

	if !control.SuccessCalled {
		t.Error("Expected first sequence to call Success on control")
	}

	// Reset control flags for second run
	control.SuccessCalled = false
	control.FailCalled = false

	// Second run: sequence2 should fail
	bt = NewBehaviorTree[int](sequence2)
	bt.SetControl(control)
	bt.SetObject(2)
	bt.Run(2)

	if !control.FailCalled {
		t.Error("Expected second sequence to call Fail on control")
	}
}

func TestBehaviorTree_ResetStateAfterRun(t *testing.T) {
	control := NewMockNode[int](t)

	task := NewTask[int](func(task *Task[int], object int) {
		task.Success()
	})

	sequence := NewSequence[int]([]Node[int]{task})

	bt := NewBehaviorTree[int](sequence)
	bt.SetControl(control)
	bt.SetObject(10)
	bt.Run(10)

	if !control.SuccessCalled {
		t.Error("Expected first run to call Success on control")
	}

	// Reset control flags for second run
	control.SuccessCalled = false
	control.FailCalled = false

	bt.SetObject(20)
	bt.Run(20)

	if !control.SuccessCalled {
		t.Error("Expected second run to call Success on control")
	}
}

func TestNestedBehaviorTrees(t *testing.T) {
	outerControl := NewMockNode[int](t)

	innerTask := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})

	innerSequence := NewSequence[int]([]Node[int]{innerTask})
	innerBT := NewBehaviorTree[int](innerSequence)

	outerSequence := NewSequence[int]([]Node[int]{innerBT})
	outerBT := NewBehaviorTree[int](outerSequence)
	outerBT.SetControl(outerControl)
	outerBT.SetObject(100)
	outerBT.Run(100)

	if !outerControl.SuccessCalled {
		t.Error("Expected outer BehaviorTree to call Success on outer control")
	}
}

func TestBehaviorTree_RunWithoutControlNode(t *testing.T) {
	sequence := NewSequence[int]([]Node[int]{
		NewTask[int](func(task *Task[int], obj int) { task.Success() }),
	})

	bt := NewBehaviorTree[int](sequence)
	bt.SetObject(0)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("BehaviorTree.Run() panicked when ControlNode was not set: %v", r)
		}
	}()

	bt.Run(0)
}
