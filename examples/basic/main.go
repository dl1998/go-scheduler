package main

import (
	"fmt"
	"github.com/dl1998/go-scheduler/pkg/scheduler"
	"time"
)

func DemoFunction(task *scheduler.Task, message string) {
	contextCounterName := "counter"
	if task.GetFromContext(contextCounterName) == nil {
		task.SetToContext(contextCounterName, 0)
	}
	counter := task.GetFromContext(contextCounterName).(int)
	fmt.Printf("[%d] %s\n", counter, message)
	task.SetToContext(contextCounterName, counter+1)
}

func ScheduleDemoTask(scheduler scheduler.Scheduler) *scheduler.Task {
	name := "Demo Task"
	duration := 10 * time.Second
	interval := time.Second

	return scheduler.ScheduleTask(name, nil, &duration, interval, DemoFunction, "Hello, World!")
}

func main() {
	newScheduler := scheduler.New()

	task := ScheduleDemoTask(*newScheduler)

	task.Wait()

	fmt.Println()
	fmt.Println(task)
}
