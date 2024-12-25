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
