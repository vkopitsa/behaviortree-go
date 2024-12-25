package main

import (
	"behaviortree"
	"fmt"
	"time"
)

// Dog represents an entity with behaviors.
type RunningDog struct {
	Name string
}

func (d *RunningDog) Bark() {
	fmt.Println(d.Name, "barks!")
}

func (d *RunningDog) RandomlyWalk() {
	fmt.Println(d.Name, "walks randomly.")
}

func (d *RunningDog) StandBesideATree() bool {
	return true
}

func (d *RunningDog) LiftALeg() {
	fmt.Println(d.Name, "lifts a leg.")
}

func (d *RunningDog) Pee() {
	fmt.Println(d.Name, "pees on the tree.")
}

func main() {
	count := 0

	// Define the 'looking' task.
	looking := behaviortree.NewTask[*RunningDog](func(task *behaviortree.Task[*RunningDog], dog *RunningDog) {
		dog.Bark()
		task.Success()
	})

	// Define the 'runningTask'.
	runningTask := behaviortree.NewTask[*RunningDog](func(task *behaviortree.Task[*RunningDog], dog *RunningDog) {
		fmt.Println("Running task")
		if count == 10 {
			task.Success()
		} else {
			task.Running()
		}
	})

	// Create a sequence node.
	sequence := behaviortree.NewSequence[*RunningDog]([]behaviortree.Node[*RunningDog]{looking, runningTask})

	// Create the behavior tree.
	tree := behaviortree.NewBehaviorTree[*RunningDog](sequence)

	// Create a Dog instance.
	dog := &RunningDog{Name: "Frank"}

	// Assign the Dog to the behavior tree.
	tree.SetObject(dog)

	// Run the behavior tree
	for i := 0; i < 5; i++ {
		fmt.Printf("\n-- Tick %d --\n", i+1)
		tree.Run(dog)
		time.Sleep(500 * time.Millisecond) // Simulate time between ticks
	}
}
