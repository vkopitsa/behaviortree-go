/*
Package behaviortree implements behavior tree nodes for creating complex behavioral structures.
This file defines the BehaviorTree structure, which serves as the root of a behavior tree.

Example Usage:

	package main

	import (
		"fmt"
		"behaviortree"
	)

	type ExampleObject struct {
		Name      string
		Condition bool
	}

	func MainTask(task *behaviortree.Task[*ExampleObject], obj *ExampleObject) {
		fmt.Printf("%s is performing the main task.\n", obj.Name)
		task.Success()
	}

	func main() {
		mainTask := behaviortree.NewTask(MainTask)

		tree := behaviortree.NewBehaviorTree(mainTask)

		obj := &ExampleObject{Name: "Agent", Condition: true}

		tree.SetObject(obj)

		tree.Run(obj)
	}

Benchmark tests:
	$ go test -bench=. -benchtime=5s -test.benchmem
	BenchmarkSequence_Success-11                    282116228               21.77 ns/op            0 B/op          0 allocs/op
	BenchmarkSequence_Fail-11                       356140072               16.94 ns/op            0 B/op          0 allocs/op
	BenchmarkPriority_Success-11                    290776261               20.25 ns/op            0 B/op          0 allocs/op
	BenchmarkPriority_Fail-11                       292416607               20.38 ns/op            0 B/op          0 allocs/op
	BenchmarkAlwaysSucceedDecorator-11              749381368                8.040 ns/op           0 B/op          0 allocs/op
	BenchmarkAlwaysFailDecorator-11                 760693929                8.036 ns/op           0 B/op          0 allocs/op
	BenchmarkInvertDecorator-11                     767313714                7.817 ns/op           0 B/op          0 allocs/op
	BenchmarkUntilFailDecorator-11                  179043154               33.36 ns/op            0 B/op          0 allocs/op
	BenchmarkCreateTask-11                          392671045               15.31 ns/op           48 B/op          1 allocs/op
	BenchmarkCreateSequence-11                      189564435               31.66 ns/op          112 B/op          2 allocs/op
	BenchmarkCreatePriority-11                      187834677               32.30 ns/op          112 B/op          2 allocs/op
	BenchmarkCreateAlwaysSucceedDecorator-11        362021664               17.07 ns/op           48 B/op          1 allocs/op
	BenchmarkCreateAlwaysFailDecorator-11           359627308               16.40 ns/op           48 B/op          1 allocs/op
	BenchmarkCreateInvertDecorator-11               366092287               16.36 ns/op           48 B/op          1 allocs/op
	BenchmarkCreateUntilFailDecorator-11            368578960               16.33 ns/op           48 B/op          1 allocs/op
	BenchmarkCreateBehaviorTree-11                  390397416               15.67 ns/op           48 B/op          1 allocs/op
*/
package behaviortree