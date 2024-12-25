package main

import (
	"fmt"
	"math/rand"
	"time"

	"behaviortree"
)

// SmartDog represents a dog with comprehensive behaviors.
type SmartDog struct {
	Name             string
	BatteryLevel     int
	PatrolPoints     []string
	CurrentPatrol    int
	IntruderDetected bool
	HuntCount        int
}

// Bark simulates the dog barking.
func (d *SmartDog) Bark() {
	fmt.Printf("%s barks loudly!\n", d.Name)
}

// Patrol simulates the dog patrolling to the next point.
func (d *SmartDog) Patrol() {
	if d.CurrentPatrol >= len(d.PatrolPoints) {
		d.CurrentPatrol = 0 // Loop back to start
	}
	fmt.Printf("%s is patrolling to %s.\n", d.Name, d.PatrolPoints[d.CurrentPatrol])
	d.CurrentPatrol++
	time.Sleep(500 * time.Millisecond) // Simulate time taken to patrol
	d.BatteryLevel -= 10                  // Simulate battery drain
}

// DetectIntruder simulates detecting an intruder.
func (d *SmartDog) DetectIntruder() bool {
	// Randomly detect an intruder with 30% probability
	d.IntruderDetected = rand.Intn(100) < 30
	if d.IntruderDetected {
		fmt.Printf("%s has detected an intruder!\n", d.Name)
	} else {
		fmt.Printf("%s found nothing suspicious.\n", d.Name)
	}
	return d.IntruderDetected
}

// ChaseIntruder simulates the dog chasing the intruder.
func (d *SmartDog) ChaseIntruder() bool {
	d.HuntCount++
	fmt.Printf("%s is chasing the intruder. Attempt %d.\n", d.Name, d.HuntCount)
	time.Sleep(500 * time.Millisecond) // Simulate time taken to chase
	if d.HuntCount >= 2 {
		fmt.Printf("%s has caught the intruder!\n", d.Name)
		d.HuntCount = 0 // Reset hunt count after catching
		return true
	}
	fmt.Printf("%s is still chasing...\n", d.Name)
	return false
}

// AlertOwner simulates the dog alerting the owner.
func (d *SmartDog) AlertOwner() {
	fmt.Printf("%s alerts the owner about the intruder!\n", d.Name)
}

// Recharge simulates the dog recharging its battery.
func (d *SmartDog) Recharge() {
	fmt.Printf("%s is recharging its battery.\n", d.Name)
	time.Sleep(1 * time.Second) // Simulate time taken to recharge
	d.BatteryLevel = 100
	fmt.Printf("%s has recharged its battery to %d%%.\n", d.Name, d.BatteryLevel)
}

// CheckBattery simulates checking the battery level.
func (d *SmartDog) CheckBattery() bool {
	if d.BatteryLevel < 20 {
		fmt.Printf("%s has low battery!\n", d.Name)
		return false
	}
	fmt.Printf("%s battery level is sufficient (%d%%).\n", d.Name, d.BatteryLevel)
	return true
}

// Main function to set up and run the Behavior Tree.
func main() {
	rand.Seed(time.Now().UnixNano())

	// Create a SmartDog instance.
	dog := &SmartDog{
		Name:          "Rex",
		BatteryLevel:  100,
		PatrolPoints:  []string{"North Gate", "East Wing", "South Gate", "West Wing"},
		CurrentPatrol: 0,
		HuntCount:     0,
	}

	// Define tasks.

	// Task to check battery
	checkBatteryTask := behaviortree.NewTask[*SmartDog](func(task *behaviortree.Task[*SmartDog], dog *SmartDog) {
		if dog.CheckBattery() {
			task.Success()
		} else {
			task.Fail()
		}
	})

	// Task to recharge
	rechargeTask := behaviortree.NewTask[*SmartDog](func(task *behaviortree.Task[*SmartDog], dog *SmartDog) {
		dog.Recharge()
		task.Success()
	})

	// Decorator to always succeed (used to ensure recharge always succeeds)
	alwaysSucceedRecharge := behaviortree.NewAlwaysSucceedDecorator[*SmartDog](rechargeTask)

	// Sequence to handle low battery
	lowBatterySequence := behaviortree.NewSequence[*SmartDog]([]behaviortree.Node[*SmartDog]{checkBatteryTask, alwaysSucceedRecharge})

	// Task to bark
	barkTask := behaviortree.NewTask[*SmartDog](func(task *behaviortree.Task[*SmartDog], dog *SmartDog) {
		dog.Bark()
		task.Success()
	})

	// Task to patrol
	patrolTask := behaviortree.NewTask[*SmartDog](func(task *behaviortree.Task[*SmartDog], dog *SmartDog) {
		dog.Patrol()
		task.Success()
	})

	// Task to detect intruder
	detectIntruderTask := behaviortree.NewTask[*SmartDog](func(task *behaviortree.Task[*SmartDog], dog *SmartDog) {
		detected := dog.DetectIntruder()
		if detected {
			task.Success()
		} else {
			task.Fail()
		}
	})

	// Task to chase intruder
	chaseIntruderTask := behaviortree.NewTask[*SmartDog](func(task *behaviortree.Task[*SmartDog], dog *SmartDog) {
		caught := dog.ChaseIntruder()
		if caught {
			task.Success()
		} else {
			task.Running()
		}
	})

	// Task to alert owner
	alertOwnerTask := behaviortree.NewTask[*SmartDog](func(task *behaviortree.Task[*SmartDog], dog *SmartDog) {
		dog.AlertOwner()
		task.Success()
	})

	// Invert decorator to handle cases when no intruder is detected
	invertDetect := behaviortree.NewInvertDecorator[*SmartDog](detectIntruderTask)

	// Decorator to keep chasing until intruder is caught
	untilIntruderCaught := behaviortree.NewUntilFailDecorator[*SmartDog](chaseIntruderTask)

	// Sequence to handle intruder detection and chase until caught
	handleIntruderWithRetrySequence := behaviortree.NewSequence[*SmartDog]([]behaviortree.Node[*SmartDog]{detectIntruderTask, untilIntruderCaught, alertOwnerTask})

	// Sequence to handle no intruder detected and continue patrolling
	handleNoIntruderSequence := behaviortree.NewSequence[*SmartDog]([]behaviortree.Node[*SmartDog]{invertDetect, patrolTask})

	// Priority (Selector) to choose between handling intruder or patrolling
	handleOrPatrolSelector := behaviortree.NewPriority[*SmartDog]([]behaviortree.Node[*SmartDog]{handleIntruderWithRetrySequence, handleNoIntruderSequence})

	// Sequence to perform patrol and then handle intruders
	patrolAndHandleSequence := behaviortree.NewRandom[*SmartDog]([]behaviortree.Node[*SmartDog]{patrolTask, barkTask, handleOrPatrolSelector})

	// Define a task that always fails, wrapped with AlwaysFailDecorator
	alwaysFailTask := behaviortree.NewTask[*SmartDog](func(task *behaviortree.Task[*SmartDog], dog *SmartDog) {
		fmt.Printf("%s attempts an impossible task and fails.\n", dog.Name)
		task.Fail()
	})
	alwaysFailDecorator := behaviortree.NewAlwaysFailDecorator[*SmartDog](alwaysFailTask)

	// Integrate the AlwaysFailDecorator into the Behavior Tree
	// For demonstration, we'll add it to the patrolAndHandleSequence
	// without affecting the main sequence.
	// To do this, we'll create a separate sequence that includes the AlwaysFailDecorator
	// and integrate it as a separate branch using a Priority node.

	// Sequence with AlwaysFailDecorator
	failSequence := behaviortree.NewSequence[*SmartDog]([]behaviortree.Node[*SmartDog]{alwaysFailDecorator})

	// Priority node to run patrolAndHandleSequence or failSequence
	// This ensures that failSequence does not block the main sequence
	mainPriority := behaviortree.NewPriority[*SmartDog]([]behaviortree.Node[*SmartDog]{patrolAndHandleSequence, failSequence})

	// Main sequence: check battery, recharge if necessary, bark, then perform main priority actions
	mainSequence := behaviortree.NewSequence[*SmartDog]([]behaviortree.Node[*SmartDog]{lowBatterySequence, barkTask, mainPriority})

	// Create the behavior tree
	tree := behaviortree.NewBehaviorTree[*SmartDog](mainSequence)

	// Assign the dog to the behavior tree
	tree.SetObject(dog)

	// Run the behavior tree
	fmt.Println("=== Running Comprehensive Behavior Tree ===")
	for i := 0; i < 10; i++ {
		fmt.Printf("\n-- Tick %d --\n", i+1)
		tree.Run(dog)
		time.Sleep(1 * time.Second) // Simulate time between ticks
	}
}
