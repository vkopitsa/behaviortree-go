package main

import (
	"fmt"
	"time"

	"github.com/vkopitsa/behaviortree-go"
)

// HuntingDog represents a dog with hunting behaviors.
type HuntingDog struct {
	Name       string
	PreyFound  bool
	ChaseCount int
}

// Bark simulates the dog barking.
func (d *HuntingDog) Bark() {
	fmt.Println(d.Name, "barks loudly!")
}

// SearchPrey simulates the dog searching for prey.
func (d *HuntingDog) SearchPrey() bool {
	fmt.Println(d.Name, "is searching for prey.")
	// Simulate prey search with a delay
	time.Sleep(1 * time.Second)
	// For demonstration, we assume prey is found
	d.PreyFound = true
	return d.PreyFound
}

// ChasePrey simulates the dog chasing the prey.
func (d *HuntingDog) ChasePrey() bool {
	fmt.Println(d.Name, "is chasing the prey.")
	d.ChaseCount++
	// Simulate chase with a delay
	time.Sleep(1 * time.Second)
	if d.ChaseCount >= 3 {
		fmt.Println(d.Name, "has caught the prey!")
		return true
	}
	fmt.Println(d.Name, "is still chasing...")
	return false
}

// Rest simulates the dog resting.
func (d *HuntingDog) Rest() {
	fmt.Println(d.Name, "is resting after a long day.")
}

func main() {
	// Create a HuntingDog instance.
	dog := &HuntingDog{Name: "Hunter"}

	// Define tasks.

	// Task to bark
	barkTask := behaviortree.NewTask[*HuntingDog](func(task *behaviortree.Task[*HuntingDog], dog *HuntingDog) {
		dog.Bark()
		task.Success()
	})

	// Task to search for prey
	searchPreyTask := behaviortree.NewTask[*HuntingDog](func(task *behaviortree.Task[*HuntingDog], dog *HuntingDog) {
		found := dog.SearchPrey()
		if found {
			task.Success()
		} else {
			task.Fail()
		}
	})

	// Task to chase prey
	chasePreyTask := behaviortree.NewTask[*HuntingDog](func(task *behaviortree.Task[*HuntingDog], dog *HuntingDog) {
		caught := dog.ChasePrey()
		if caught {
			task.Success()
		} else {
			task.Running()
		}
	})

	// Task to rest
	restTask := behaviortree.NewTask[*HuntingDog](func(task *behaviortree.Task[*HuntingDog], dog *HuntingDog) {
		dog.Rest()
		task.Success()
	})

	// Sequence to search and chase prey
	searchAndChaseSequence := behaviortree.NewSequence[*HuntingDog]([]behaviortree.Node[*HuntingDog]{searchPreyTask, chasePreyTask})

	// Priority (Selector) to choose between hunting and resting
	huntOrRestSelector := behaviortree.NewPriority[*HuntingDog]([]behaviortree.Node[*HuntingDog]{searchAndChaseSequence, restTask})

	// Main sequence: bark and then hunt or rest
	mainSequence := behaviortree.NewSequence[*HuntingDog]([]behaviortree.Node[*HuntingDog]{barkTask, huntOrRestSelector})

	// Create the behavior tree
	tree := behaviortree.NewBehaviorTree[*HuntingDog](mainSequence)

	// Assign the dog to the behavior tree
	tree.SetObject(dog)

	// Optionally, assign a control node if needed
	// For simplicity, omit control node here.

	// Run the behavior tree
	fmt.Println("=== Running Hunting Behavior Tree ===")
	for i := 0; i < 5; i++ {
		fmt.Printf("\n-- Tick %d --\n", i+1)
		tree.Run(dog)
		time.Sleep(500 * time.Millisecond) // Simulate time between ticks
	}
}
