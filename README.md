
# BehaviorTree-Go

[![GitHub Workflow Status](https://github.com/vkopitsa/behaviortree-go/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/vkopitsa/behaviortree-go/actions/workflows/test.yml)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/vkopitsa/behaviortree-go)
[![License](http://img.shields.io/badge/license-apache-blue.svg?style=flat-square)](https://raw.githubusercontent.com/vkopitsa/behaviortree-go/master/LICENSE)


BehaviorTree-Go is a Go implementation of Behavior Trees, a powerful tool for organizing complex decision-making processes. Inspired by the Lua implementation from [BehaviourTree.lua](https://github.com/tanema/behaviourtree.lua), this library facilitates the creation of sophisticated behaviors for applications such as AI in video games, robotics, and [more](https://arxiv.org/pdf/1709.00084).

## Features

- **Flexible Structure**: Easily create complex behavior trees using sequences, selectors, and decorators.
- **Extensible**: Add custom nodes and decorators tailored to your specific needs.
- **Well-Tested**: Comprehensive unit tests ensure reliability and stability.

## Benchmark

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

## Installation

To include BehaviorTree-Go in your project, you can use `go get`:

```bash
go get github.com/vkopitsa/behaviortree-go
```

Alternatively, you can clone the repository and include it manually:

```bash
git clone https://github.com/vkopitsa/behaviortree-go.git
```

## Usage

### Creating a Simple Task

A `Task` is a basic unit in a Behavior Tree that performs an action. Implement the `Run` method to define the task's behavior.

```go
package main

import (
	"fmt"
	"time"

	"github.com/vkopitsa/behaviortree-go"
)

// ExampleObject represents the context in which the behavior tree operates.
type ExampleObject struct {
	Name      string
	Condition bool
}

// MainTask is a simple task that prints a message and succeeds.
func MainTask(task *behaviortree.Task[*ExampleObject], obj *ExampleObject) {
	fmt.Printf("%s is performing the main task.", obj.Name)
	task.Success()
}

func main() {
	// Create a new task.
	mainTask := behaviortree.NewTask(MainTask)

	// Create the behavior tree with the main task as the root.
	tree := behaviortree.NewBehaviorTree(mainTask)

	// Create an example object.
	obj := &ExampleObject{Name: "Agent", Condition: true}

	// Assign the object to the behavior tree.
	tree.SetObject(obj)

	// Run the behavior tree.
	tree.Run(obj)
}
```

### Creating a Sequence

A `Sequence` node executes its child nodes in order. It succeeds only if all child nodes succeed. If any child fails, the sequence fails immediately.

```go
package main

import (
	"fmt"

	"github.com/vkopitsa/behaviortree-go"
)

// ExampleObject represents the context in which the behavior tree operates.
type ExampleObject struct {
	Name      string
	Condition bool
}

// TaskOne is the first task in the sequence.
func TaskOne(task *behaviortree.Task[*ExampleObject], obj *ExampleObject) {
	fmt.Printf("%s is executing Task One.", obj.Name)
	task.Success()
}

// TaskTwo is the second task in the sequence.
func TaskTwo(task *behaviortree.Task[*ExampleObject], obj *ExampleObject) {
	fmt.Printf("%s is executing Task Two.", obj.Name)
	task.Success()
}

func main() {
	// Create tasks.
	taskOne := behaviortree.NewTask(TaskOne)
	taskTwo := behaviortree.NewTask(TaskTwo)

	// Create a sequence with the tasks.
	sequence := behaviortree.NewSequence([]behaviortree.Node[*ExampleObject]{taskOne, taskTwo})

	// Create the behavior tree with the sequence as the root.
	tree := behaviortree.NewBehaviorTree(sequence)

	// Create an example object.
	obj := &ExampleObject{Name: "Agent", Condition: true}

	// Assign the object to the behavior tree.
	tree.SetObject(obj)

	// Run the behavior tree.
	tree.Run(obj)
}
```

### Custom Decorators

To implement a custom decorator, embed the `Decorator` struct and override the required methods. For example:

```go
package main

import (
	"fmt"

	"github.com/vkopitsa/behaviortree-go"
)

// LogDecorator logs the execution of its child node.
type LogDecorator[T any] struct {
	behaviortree.Decorator[T]
}

func NewLogDecorator[T any](node behaviortree.Node[T]) *LogDecorator[T] {
	decorator := &LogDecorator[T]{}
	decorator.Node = node
	decorator.Node.SetControl(decorator)
	return decorator
}

func (d *LogDecorator[T]) Success() {
	fmt.Println("LogDecorator: Child succeeded.")
	if d.ControlNode != nil {
		d.ControlNode.Success()
	}
}

func (d *LogDecorator[T]) Fail() {
	fmt.Println("LogDecorator: Child failed.")
	if d.ControlNode != nil {
		d.ControlNode.Fail()
	}
}

func main() {
	// Create a simple task.
	task := behaviortree.NewTask(func(task *behaviortree.Task[ExampleObject], obj *ExampleObject) {
		fmt.Println("Executing task...")
		task.Success()
	})

	// Wrap the task with a LogDecorator.
	logDecorator := NewLogDecorator(task)

	// Create the behavior tree with the decorator as the root.
	tree := behaviortree.NewBehaviorTree(logDecorator)

	// Create an example object.
	obj := &ExampleObject{Name: "Agent"}

	// Assign the object to the behavior tree.
	tree.SetObject(obj)

	// Run the behavior tree.
	tree.Run(obj)
}
```

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Commit your changes with clear messages.
4. Open a pull request detailing your changes.

Please ensure all tests pass and adhere to the project's coding standards.

## License

BehaviorTree-Go is released under the [Apache License 2.0](LICENSE).

## Acknowledgments

- Inspired by [BehaviourTree.lua](https://github.com/tanema/behaviourtree.lua)
- Thanks to the open-source community for their invaluable contributions and support.
