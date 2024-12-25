package main

import (
	"fmt"
	"time"
	"math/rand"
	"behaviortree"
)

// ExampleObject represents the context in which the behavior tree operates.
type ExampleObject struct {}

// SuccessTask creates a task that always succeeds.
// It prints a success message and signals the task's success.
func SuccessTask(name string) *behaviortree.Task[*ExampleObject] {
	return behaviortree.NewTask(func(task *behaviortree.Task[*ExampleObject], obj *ExampleObject) {
		fmt.Printf("%s is performing the success task.\n", name)
		task.Success()
	})
}

// FailTask creates a task that always fails.
// It prints a failure message and signals the task's failure.
func FailTask(name string) *behaviortree.Task[*ExampleObject] {
	return behaviortree.NewTask(func(task *behaviortree.Task[*ExampleObject], obj *ExampleObject) {
		fmt.Printf("%s is performing the fail task.\n", name)
		task.Fail()
	})
}

// RandomTask creates a task that randomly succeeds or fails.
// It prints either a success or failure message based on a random value.
func RandomTask(name string) *behaviortree.Task[*ExampleObject] {
	return behaviortree.NewTask(func(task *behaviortree.Task[*ExampleObject], obj *ExampleObject) {
		if rand.Intn(2) == 1 {
			fmt.Printf("%s is performing success of the random task.\n", name)
			task.Success()
		} else {
			fmt.Printf("%s is performing fail of the random task.\n", name)
			task.Fail()
		}
	})
}

// main initializes and executes the behavior tree example.
// The tree consists of a priority node with multiple children, including sequences, decorators, and tasks.
//
// Behavior Tree Structure:
//   - Root Priority Node
//     - RandomTask("Priority 1")
//     - InvertDecorator wrapping RandomTask("InvertDecorator 1")
//     - Sequence Node
//       - SuccessTask("Sequence 1")
//       - SuccessTask("Sequence 2")
//       - FailTask("Sequence 3")
//       - SuccessTask("Sequence 4")
//     - FailTask("Priority 2")
//     - AlwaysFailDecorator wrapping RandomTask("AlwaysFailDecorator 1")
//     - AlwaysSucceedDecorator wrapping RandomTask("AlwaysSucceedDecorator 1")
//     - SuccessTask("Priority 3")
func main() {
	rand.Seed(time.Now().UnixNano()) // Initialize random seed

	// Define the behavior tree structure
	priority := behaviortree.NewPriority([]behaviortree.Node[*ExampleObject]{
		RandomTask("Priority 1"),
		behaviortree.NewInvertDecorator(RandomTask("InvertDecorator 1")),
		behaviortree.NewSequence([]behaviortree.Node[*ExampleObject]{
			SuccessTask("Sequence 1"),
			SuccessTask("Sequence 2"),
			FailTask("Sequence 3"),
			SuccessTask("Sequence 4"),
		}),
		FailTask("Priority 2"),
		behaviortree.NewAlwaysFailDecorator(RandomTask("AlwaysFailDecorator 1")),
		behaviortree.NewAlwaysSucceedDecorator(RandomTask("AlwaysSucceedDecorator 1")),
		SuccessTask("Priority 3"),
	})

	// Create the behavior tree with the priority node as the root
	tree := behaviortree.NewBehaviorTree(priority)

	// Create an example object representing the agent
	obj := &ExampleObject{}

	// Assign the object to the behavior tree
	tree.SetObject(obj)

	// Run the behavior tree for multiple ticks
	for i := 0; i < 5; i++ {
		fmt.Printf("\n-- Tick %d --\n", i+1)
		tree.Run(obj)
		time.Sleep(500 * time.Millisecond) // Simulate time between ticks
	}
}