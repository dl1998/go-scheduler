// Example that shows how to use Task context.
package main

import (
	"fmt"
	"github.com/dl1998/go-scheduler/pkg/scheduler"
	"time"
)

func DemoFunction(task *scheduler.Task) {
	contextCounterName := "counter"
	if task.GetFromContext(contextCounterName) == nil {
		task.SetToContext(contextCounterName, 0)
	}
	counter := task.GetFromContext(contextCounterName).(int)
	fmt.Printf("[%d] Hello, World!\n", counter)
	task.SetToContext(contextCounterName, counter+1)
}

func ScheduleDemoTask(scheduler *scheduler.Scheduler) *scheduler.Task {
	name := "Demo Task"
	duration := 10 * time.Second
	interval := time.Second

	return scheduler.ScheduleTask(name, nil, &duration, interval, DemoFunction)
}

func main() {
	newScheduler := scheduler.New()

	task := ScheduleDemoTask(newScheduler)

	task.Wait()

	fmt.Println()
	fmt.Println(task)
}
