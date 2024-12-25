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
