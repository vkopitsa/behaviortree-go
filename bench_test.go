package behaviortree

import (
	"testing"
)

func BenchmarkSequence_Success(b *testing.B) {
	task := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})
	sequence := NewSequence[int]([]Node[int]{task, task, task})

	bt := NewBehaviorTree[int](sequence)
	bt.SetObject(0)

	for i := 0; i < b.N; i++ {
		bt.Run(0)
	}
}

func BenchmarkSequence_Fail(b *testing.B) {
	taskFail := NewTask[int](func(task *Task[int], obj int) {
		task.Fail()
	})
	taskSuccess := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})
	sequence := NewSequence[int]([]Node[int]{taskSuccess, taskFail, taskSuccess})

	bt := NewBehaviorTree[int](sequence)
	bt.SetObject(0)

	for i := 0; i < b.N; i++ {
		bt.Run(0)
	}
}

func BenchmarkPriority_Success(b *testing.B) {
	taskSuccess := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})
	taskFail := NewTask[int](func(task *Task[int], obj int) {
		task.Fail()
	})
	priority := NewPriority[int]([]Node[int]{taskFail, taskFail, taskSuccess})

	bt := NewBehaviorTree[int](priority)
	bt.SetObject(0)

	for i := 0; i < b.N; i++ {
		bt.Run(0)
	}
}

func BenchmarkPriority_Fail(b *testing.B) {
	taskFail := NewTask[int](func(task *Task[int], obj int) {
		task.Fail()
	})
	priority := NewPriority[int]([]Node[int]{taskFail, taskFail, taskFail})

	bt := NewBehaviorTree[int](priority)
	bt.SetObject(0)

	for i := 0; i < b.N; i++ {
		bt.Run(0)
	}
}

func BenchmarkAlwaysSucceedDecorator(b *testing.B) {
	task := NewTask[int](func(task *Task[int], obj int) {
		task.Fail() // Normally fails
	})
	alwaysSucceed := NewAlwaysSucceedDecorator(task)

	bt := NewBehaviorTree[int](alwaysSucceed)
	bt.SetObject(0)

	for i := 0; i < b.N; i++ {
		bt.Run(0)
	}
}

func BenchmarkAlwaysFailDecorator(b *testing.B) {
	task := NewTask[int](func(task *Task[int], obj int) {
		task.Success() // Normally succeeds
	})
	alwaysFail := NewAlwaysFailDecorator(task)

	bt := NewBehaviorTree[int](alwaysFail)
	bt.SetObject(0)

	for i := 0; i < b.N; i++ {
		bt.Run(0)
	}
}

func BenchmarkInvertDecorator(b *testing.B) {
	task := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})
	invertDecorator := NewInvertDecorator(task)

	bt := NewBehaviorTree[int](invertDecorator)
	bt.SetObject(0)

	for i := 0; i < b.N; i++ {
		bt.Run(0)
	}
}

func BenchmarkUntilFailDecorator(b *testing.B) {
	callCount := 0
	mockNode := NewMockNode[int](nil)
	mockNode.CustomRun = func(m *MockNode[int], obj int) {
		callCount++
		if callCount < 5 {
			m.Control.Success()
		} else {
			m.Control.Fail()
		}
	}

	untilFail := NewUntilFailDecorator[int](mockNode)
	bt := NewBehaviorTree[int](untilFail)
	bt.SetObject(0)

	for i := 0; i < b.N; i++ {
		bt.Run(0)
		callCount = 0 // Reset call count for each iteration
	}
}

var sink interface{} // Global variable to prevent compiler optimizations

func BenchmarkCreateTask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sink = NewTask(func(task *Task[int], obj int) {
			task.Success()
		})
	}
}

func BenchmarkCreateSequence(b *testing.B) {
	task := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})
	for i := 0; i < b.N; i++ {
		sink = NewSequence[int]([]Node[int]{task, task, task})
	}
}

func BenchmarkCreatePriority(b *testing.B) {
	task := NewTask[int](func(task *Task[int], obj int) {
		task.Fail()
	})
	for i := 0; i < b.N; i++ {
		sink = NewPriority[int]([]Node[int]{task, task, task})
	}
}

func BenchmarkCreateAlwaysSucceedDecorator(b *testing.B) {
	task := NewTask[int](func(task *Task[int], obj int) {
		task.Fail()
	})
	for i := 0; i < b.N; i++ {
		sink = NewAlwaysSucceedDecorator(task)
	}
}

func BenchmarkCreateAlwaysFailDecorator(b *testing.B) {
	task := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})
	for i := 0; i < b.N; i++ {
		sink = NewAlwaysFailDecorator(task)
	}
}

func BenchmarkCreateInvertDecorator(b *testing.B) {
	task := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})
	for i := 0; i < b.N; i++ {
		sink = NewInvertDecorator(task)
	}
}

func BenchmarkCreateUntilFailDecorator(b *testing.B) {
	task := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})
	for i := 0; i < b.N; i++ {
		sink = NewUntilFailDecorator(task)
	}
}

func BenchmarkCreateBehaviorTree(b *testing.B) {
	task := NewTask[int](func(task *Task[int], obj int) {
		task.Success()
	})
	for i := 0; i < b.N; i++ {
		sink = NewBehaviorTree[int](task)
	}
}
