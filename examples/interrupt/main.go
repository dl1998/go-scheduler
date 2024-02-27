// Example that shows how to interrupt task that is currently scheduled.
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
	fmt.Println(counter)
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

	time.Sleep(5 * time.Second)

	if err := newScheduler.StopTask(task); err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(task)
}
