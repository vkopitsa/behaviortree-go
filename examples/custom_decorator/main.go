package main

import (
	"fmt"

	"behaviortree"
)

// ExampleObject represents the context in which the behavior tree operates.
type ExampleObject struct {
	Name      string
	Condition bool
}

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
	task := behaviortree.NewTask(func(task *behaviortree.Task[*ExampleObject], obj *ExampleObject) {
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
