package behaviortree

// Task represents a node in a behavior tree that encapsulates a specific action or task to be executed.
// It is generic, allowing it to work with any type `T`.
type Task[T any] struct {
	BaseNode[T] // Inherits functionality from BaseNode for tree-related operations.

	// RunFunc defines the function to be executed when this task is run.
	// It receives the current task instance and an object of type `T`.
	RunFunc func(task *Task[T], object T)

	// RunCalled indicates whether the Run method has been called at least once.
	RunCalled bool
}

// NewTask creates a new instance of Task with the specified run function.
func NewTask[T any](runFunc func(task *Task[T], object T)) *Task[T] {
	return &Task[T]{
		RunFunc: runFunc,
	}
}

// Run executes the task's `RunFunc` with the provided object.
func (t *Task[T]) Run(object T) {
	if t.RunFunc != nil {
		t.RunFunc(t, object)
		t.RunCalled = true
	}
}