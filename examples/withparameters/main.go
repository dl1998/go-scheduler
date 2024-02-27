// Example that shows how to schedule task and provide parameters to the
// scheduled function.
package main

import (
	"fmt"
	"github.com/dl1998/go-scheduler/pkg/scheduler"
	"os"
	"time"
)

func DemoFunction(task *scheduler.Task, name string) {
	fmt.Printf("Hello, %s!\n", name)
}

func ScheduleDemoTask(scheduler *scheduler.Scheduler) *scheduler.Task {
	name := "Demo Task"
	duration := 10 * time.Second
	interval := time.Second

	return scheduler.ScheduleTask(name, nil, &duration, interval, DemoFunction, os.Getenv("USER"))
}

func main() {
	newScheduler := scheduler.New()

	task := ScheduleDemoTask(newScheduler)

	task.Wait()

	fmt.Println()
	fmt.Println(task)
}
