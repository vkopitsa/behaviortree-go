package main

import (
	"fmt"
	"time"

	"github.com/vkopitsa/behaviortree-go"
)

// RandomDog represents an entity with behaviors.
type RandomDog struct {
	Name string
}

func (d *RandomDog) Bark() {
	fmt.Println(d.Name, "barks!")
}

func (d *RandomDog) WalkLeft() {
	fmt.Println(d.Name, "walks left.")
}

func (d *RandomDog) WalkUp() {
	fmt.Println(d.Name, "walks up.")
}

func (d *RandomDog) WalkRight() {
	fmt.Println(d.Name, "walks right.")
}

func (d *RandomDog) WalkDown() {
	fmt.Println(d.Name, "walks down.")
}

func main() {
	// Define tasks.
	looking := behaviortree.NewTask[*RandomDog](func(task *behaviortree.Task[*RandomDog], dog *RandomDog) {
		dog.Bark()
		task.Success()
	})

	walkLeft := behaviortree.NewTask[*RandomDog](func(task *behaviortree.Task[*RandomDog], dog *RandomDog) {
		dog.WalkLeft()
		task.Success()
	})

	walkUp := behaviortree.NewTask[*RandomDog](func(task *behaviortree.Task[*RandomDog], dog *RandomDog) {
		dog.WalkUp()
		task.Success()
	})

	walkRight := behaviortree.NewTask[*RandomDog](func(task *behaviortree.Task[*RandomDog], dog *RandomDog) {
		dog.WalkRight()
		task.Success()
	})

	walkDown := behaviortree.NewTask[*RandomDog](func(task *behaviortree.Task[*RandomDog], dog *RandomDog) {
		dog.WalkDown()
		task.Success()
	})

	// Create a Random selector node.
	randomWalk := behaviortree.NewRandom[*RandomDog]([]behaviortree.Node[*RandomDog]{walkLeft, walkUp, walkRight, walkDown})

	// Create a sequence node.
	sequence := behaviortree.NewSequence[*RandomDog]([]behaviortree.Node[*RandomDog]{looking, randomWalk})

	// Create the behavior tree.
	tree := behaviortree.NewBehaviorTree[*RandomDog](sequence)

	// Create a Dog instance.
	dog := &RandomDog{Name: "Frank"}

	// Assign the Dog to the behavior tree.
	tree.SetObject(dog)

	// Run the behavior tree
	for i := 0; i < 5; i++ {
		fmt.Printf("\n-- Tick %d --\n", i+1)
		tree.Run(dog)
		time.Sleep(500 * time.Millisecond) // Simulate time between ticks
	}
}
