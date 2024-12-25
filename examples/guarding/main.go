package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/vkopitsa/behaviortree-go"
)

// GuardDog represents a dog with guarding behaviors.
type GuardDog struct {
	Name             string
	GuardPoints      []string
	CurrentGuard     int
	IntruderDetected bool
	HuntCount        int
	BatteryLevel     int
}

// Bark simulates the dog barking.
func (d *GuardDog) Bark() {
	fmt.Printf("%s barks loudly to warn intruders!\n", d.Name)
}

// Patrol simulates the dog patrolling a specific guard point.
func (d *GuardDog) Patrol() {
	if d.CurrentGuard >= len(d.GuardPoints) {
		d.CurrentGuard = 0 // Loop back to start
	}
	currentPoint := d.GuardPoints[d.CurrentGuard]
	fmt.Printf("%s is patrolling the %s.\n", d.Name, currentPoint)
	d.CurrentGuard++
	time.Sleep(500 * time.Millisecond) // Simulate time taken to patrol
	d.BatteryLevel -= 10                  // Simulate battery drain
}

// DetectIntruder simulates detecting an intruder.
func (d *GuardDog) DetectIntruder() bool {
	// Randomly detect an intruder with 25% probability
	d.IntruderDetected = rand.Intn(100) < 25

	// Calculate the guard index safely
	guardIndex := d.CurrentGuard - 1
	if guardIndex < 0 {
		guardIndex = 0 // Default to the first guard point
	} else if guardIndex >= len(d.GuardPoints) {
		guardIndex = len(d.GuardPoints) - 1 // Default to the last guard point
	}

	if d.IntruderDetected {
		fmt.Printf("%s has detected an intruder at the %s!\n", d.Name, d.GuardPoints[guardIndex])
	} else {
		fmt.Printf("%s sees no intruders at the %s.\n", d.Name, d.GuardPoints[guardIndex])
	}
	return d.IntruderDetected
}

// ChaseIntruder simulates the dog chasing the intruder.
func (d *GuardDog) ChaseIntruder() bool {
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
func (d *GuardDog) AlertOwner() {
	fmt.Printf("%s alerts the owner about the intruder!\n", d.Name)
}

// Recharge simulates the dog recharging its battery.
func (d *GuardDog) Recharge() {
	fmt.Printf("%s is recharging its battery.\n", d.Name)
	time.Sleep(1 * time.Second) // Simulate time taken to recharge
	d.BatteryLevel = 100
	fmt.Printf("%s has recharged its battery to %d%%.\n", d.Name, d.BatteryLevel)
}

// CheckBattery simulates checking the battery level.
func (d *GuardDog) CheckBattery() bool {
	if d.BatteryLevel < 20 {
		fmt.Printf("%s has low battery!\n", d.Name)
		return false
	}
	fmt.Printf("%s battery level is sufficient (%d%%).\n", d.Name, d.BatteryLevel)
	return true
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Create a GuardDog instance.
	dog := &GuardDog{
		Name:          "Max",
		GuardPoints:   []string{"North Gate", "East Wing", "South Gate", "West Wing"},
		CurrentGuard:  0, // Initialize to 0
		BatteryLevel:  100,
		HuntCount:     0,
	}

	// Define the behavior tree in a single structure
	tree := behaviortree.NewBehaviorTree[*GuardDog](
		behaviortree.NewSequence[*GuardDog]([]behaviortree.Node[*GuardDog]{
			// Low Battery Sequence: Check battery and recharge if necessary
			behaviortree.NewSequence[*GuardDog]([]behaviortree.Node[*GuardDog]{
				// Task to check battery
				behaviortree.NewTask[*GuardDog](func(task *behaviortree.Task[*GuardDog], dog *GuardDog) {
					if dog.CheckBattery() {
						task.Success()
					} else {
						task.Fail()
					}
				}),
				// Task to recharge, wrapped with AlwaysSucceedDecorator
				behaviortree.NewAlwaysSucceedDecorator[*GuardDog](
					behaviortree.NewTask[*GuardDog](func(task *behaviortree.Task[*GuardDog], dog *GuardDog) {
						dog.Recharge()
						task.Success()
					}),
				),
			}),
			// Task to bark
			behaviortree.NewTask[*GuardDog](func(task *behaviortree.Task[*GuardDog], dog *GuardDog) {
				dog.Bark()
				task.Success()
			}),
			// Main Priority Selector: Choose between handling intruders or patrolling
			behaviortree.NewPriority[*GuardDog]([]behaviortree.Node[*GuardDog]{
				// Patrol and Handle Sequence wrapped with Random Selector
				behaviortree.NewRandom[*GuardDog]([]behaviortree.Node[*GuardDog]{
					// Task to patrol
					behaviortree.NewTask[*GuardDog](func(task *behaviortree.Task[*GuardDog], dog *GuardDog) {
						dog.Patrol()
						task.Success()
					}),
					// Task to bark
					behaviortree.NewTask[*GuardDog](func(task *behaviortree.Task[*GuardDog], dog *GuardDog) {
						dog.Bark()
						task.Success()
					}),
					// Handle or Patrol Selector
					behaviortree.NewPriority[*GuardDog]([]behaviortree.Node[*GuardDog]{
						// Handle Intruder with Retry Sequence
						behaviortree.NewSequence[*GuardDog]([]behaviortree.Node[*GuardDog]{
							// Task to detect intruder
							behaviortree.NewTask[*GuardDog](func(task *behaviortree.Task[*GuardDog], dog *GuardDog) {
								detected := dog.DetectIntruder()
								if detected {
									task.Success()
								} else {
									task.Fail()
								}
							}),
							// Decorator to keep chasing until intruder is caught
							behaviortree.NewUntilFailDecorator[*GuardDog](
								behaviortree.NewTask[*GuardDog](func(task *behaviortree.Task[*GuardDog], dog *GuardDog) {
									caught := dog.ChaseIntruder()
									if caught {
										task.Success()
									} else {
										task.Running()
									}
								}),
							),
							// Task to alert owner
							behaviortree.NewTask[*GuardDog](func(task *behaviortree.Task[*GuardDog], dog *GuardDog) {
								dog.AlertOwner()
								task.Success()
							}),
						}),
						// Handle No Intruder Sequence using InvertDecorator
						behaviortree.NewSequence[*GuardDog]([]behaviortree.Node[*GuardDog]{
							// Invert decorator to handle no intruder detection
							behaviortree.NewInvertDecorator[*GuardDog](
								behaviortree.NewTask[*GuardDog](func(task *behaviortree.Task[*GuardDog], dog *GuardDog) {
									if dog.DetectIntruder() {
										task.Success()
									} else {
										task.Running()
									}
								}),
							),
							// Task to patrol
							behaviortree.NewTask[*GuardDog](func(task *behaviortree.Task[*GuardDog], dog *GuardDog) {
								dog.Patrol()
								task.Success()
							}),
						}),
					}),
				}),
				// Sequence with AlwaysFailDecorator to demonstrate decorator behavior
				behaviortree.NewSequence[*GuardDog]([]behaviortree.Node[*GuardDog]{
					behaviortree.NewAlwaysFailDecorator[*GuardDog](
						behaviortree.NewTask[*GuardDog](func(task *behaviortree.Task[*GuardDog], dog *GuardDog) {
							fmt.Printf("%s attempts an impossible task and fails.\n", dog.Name)
							task.Fail()
						}),
					),
				}),
			}),
		}),
	)

	// Assign the dog to the behavior tree
	tree.SetObject(dog)

	// Run the behavior tree
	fmt.Println("=== Running Guarding Behavior Tree ===")
	for i := 0; i < 10; i++ {
		fmt.Printf("\n-- Tick %d --\n", i+1)
		tree.Run(dog)
		time.Sleep(1 * time.Second) // Simulate time between ticks
	}
}
