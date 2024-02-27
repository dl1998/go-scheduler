// Example that shows how to schedule simple task.
package main

import (
	"fmt"
	"github.com/dl1998/go-scheduler/pkg/scheduler"
	"time"
)

func DemoFunction(task *scheduler.Task) {
	fmt.Println("Hello, World!")
}

func ScheduleDemoTask(scheduler scheduler.Scheduler) *scheduler.Task {
	name := "Demo Task"
	duration := 10 * time.Second
	interval := time.Second

	return scheduler.ScheduleTask(name, nil, &duration, interval, DemoFunction)
}

func main() {
	newScheduler := scheduler.New()

	task := ScheduleDemoTask(*newScheduler)

	task.Wait()

	fmt.Println()
	fmt.Println(task)
}
