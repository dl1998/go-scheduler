// Package scheduler provides Scheduler implementation in Go. It consists of the
// Scheduler itself and Task that represents one scheduled item.
package scheduler

import (
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"strings"
	"time"
)

// Scheduler struct that stores scheduled task.
type Scheduler struct {
	// Tasks stores list of scheduled Task.
	Tasks []*Task
}

// New creates a new Scheduler object.
func New() *Scheduler {
	return &Scheduler{}
}

// Task represent a thing that could be scheduled using Scheduler.
type Task struct {
	// ID unique value to distinguish different tasks, could be custom, but it is
	// recommended to leave it empty.
	ID string `json:"id,omitempty"`
	// Name is used to distinguish tasks in human-readable format, it is not unique.
	Name string `json:"name"`
	// Start represents time after which task could be triggered (first schedule
	// time).
	Start *time.Time `json:"start,omitempty"`
	// Duration stores how long this task shall be keep alive by the scheduler. It is
	// not a task execution duration, but rather time during which task exists in the
	// scheduler.
	Duration *time.Duration `json:"duration,omitempty"`
	// Interval stores information how often this task shall be triggered by
	// scheduler.
	Interval time.Duration `json:"interval"`
	// stopSignal stores channel for task termination, it terminates the whole task,
	// not only current execution.
	stopSignal chan bool
	// context stores additional key-value data that are shared between different
	// task executions.
	context map[string]interface{}
}

// NewTask creates a new task struct, it handles default value initialization for
// empty values (nil or blank string).
func NewTask(id string, name string, start *time.Time, duration *time.Duration, interval time.Duration, stopSignal chan bool, context map[string]interface{}) *Task {
	if id == "" {
		id = uuid.New().String()
	}

	if start == nil {
		startTime := time.Now()
		start = &startTime
	}

	if stopSignal == nil {
		stopSignal = make(chan bool)
	}

	if context == nil {
		context = make(map[string]interface{})
	}

	return &Task{
		ID:         id,
		Name:       name,
		Start:      start,
		Duration:   duration,
		Interval:   interval,
		stopSignal: stopSignal,
		context:    context,
	}
}

// NewSimpleTask creates a new task struct where all parameters have default
// value (nil or empty string) except of name and interval.
func NewSimpleTask(name string, interval time.Duration) *Task {
	return NewTask("", name, nil, nil, interval, nil, nil)
}

// String returns human-readable string with information about the Task.
func (task *Task) String() string {
	stringBuilder := strings.Builder{}
	stringBuilder.WriteString(fmt.Sprintf("ID: %s\n", task.ID))
	stringBuilder.WriteString(fmt.Sprintf("Name: %s\n", task.Name))
	stringBuilder.WriteString(fmt.Sprintf("Start: %s\n", task.Start.String()))
	stringBuilder.WriteString(fmt.Sprintf("Duration: %s\n", task.Duration.String()))
	stringBuilder.WriteString(fmt.Sprintf("Interval: %s\n", task.Interval.String()))
	stringBuilder.WriteString(fmt.Sprintf("Context: %v", task.context))
	return stringBuilder.String()
}

// FindTaskIndex searches through tasks array for the scheduled task and returns
// its index.
func (scheduler *Scheduler) FindTaskIndex(scheduledTask *Task) int {
	for index, task := range scheduler.Tasks {
		if task == scheduledTask {
			return index
		}
	}
	return -1
}

// FindTaskByName searches through tasks array for the scheduled task with
// provided name and returns it. If task with provided name was not found, then
// returns nil.
func (scheduler *Scheduler) FindTaskByName(name string) *Task {
	for _, task := range scheduler.Tasks {
		if task.Name == name {
			return task
		}
	}
	return nil
}

// FindTaskByID searches through tasks array for the scheduled task with provided
// ID and returns it. If task with provided ID was not found, then returns nil.
func (scheduler *Scheduler) FindTaskByID(id string) *Task {
	for _, task := range scheduler.Tasks {
		if task.ID == id {
			return task
		}
	}
	return nil
}

// ScheduleTask starts a Go routine that runs a provided function with given
// parameters. It stops either after a specified duration or when a stop signal
// is received, whichever comes first. If duration is nil, it only stops when a
// stop signal is received.
func (scheduler *Scheduler) ScheduleTask(name string, startTime *time.Time, duration *time.Duration, interval time.Duration, function interface{}, parameters ...interface{}) *Task {
	// Creates termination channel for the new task.
	var terminationChannel = make(chan bool)

	// Set default start time, if it was not provided.
	if startTime == nil {
		start := time.Now()
		startTime = &start
	}

	// Calculate the duration until the start time.
	waitDuration := time.Until(*startTime)

	// Create a new Task for Scheduler.
	scheduledTask := &Task{
		ID:         uuid.New().String(),
		Name:       name,
		Start:      startTime,
		Duration:   duration,
		Interval:   interval,
		stopSignal: terminationChannel,
		context:    make(map[string]interface{}),
	}

	go func() {
		// If the start time is in the future, wait until then to start the task.
		if waitDuration > 0 {
			<-time.After(waitDuration)
		}

		// Add scheduled task to the arguments of the executed function
		parameters := append([]interface{}{scheduledTask}, parameters...)

		// After waiting until the start time, execute the task once immediately.
		callFunction(function, parameters...)

		// If a duration is specified, calculate the end time from the start time.
		var endTime time.Time
		if duration != nil {
			endTime = startTime.Add(*duration)
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		defer scheduler.removeTask(scheduledTask)

		for {
			select {
			case <-terminationChannel: // If a terminationChannel signal is received, stop the task.
				return
			case tick := <-ticker.C: // On each tick, check if the end time has been reached.
				if duration != nil && tick.After(endTime) {
					return // If the current time is after the end time, stop the task.
				}
				callFunction(function, parameters...)
			}
		}
	}()

	// Add new Task to Scheduler tasks list.
	scheduler.Tasks = append(scheduler.Tasks, scheduledTask)

	return scheduledTask
}

// StopTask encapsulates task stopping sequence.
func (scheduler *Scheduler) StopTask(task *Task) {
	close(task.stopSignal)
	scheduler.removeTask(task)
}

// removeTask removes Task from the tasks list of the Scheduler.
func (scheduler *Scheduler) removeTask(scheduledTask *Task) {
	taskIndex := scheduler.FindTaskIndex(scheduledTask)
	if len(scheduler.Tasks) <= 1 {
		scheduler.Tasks = make([]*Task, 0)
	} else {
		scheduler.Tasks = append(scheduler.Tasks[:taskIndex], scheduler.Tasks[taskIndex+1:]...)
	}
}

// callFunction calls a function dynamically using reflection.
// It panics if the function call is not valid.
func callFunction(function interface{}, parameters ...interface{}) {
	functionReflection := reflect.ValueOf(function)
	if functionReflection.Kind() != reflect.Func {
		panic("provided argument is not a function")
	}

	// Prepare parameters for reflection call.
	parametersReflection := make([]reflect.Value, len(parameters))
	for index, parameter := range parameters {
		parametersReflection[index] = reflect.ValueOf(parameter)
	}

	// Call the function with the parameters.
	functionReflection.Call(parametersReflection)
}

// GetFromContext receive value for the provided key from the task context.
func (task *Task) GetFromContext(name string) interface{} {
	return task.context[name]
}

// SetToContext adds key-value pair to the task context.
func (task *Task) SetToContext(name string, value interface{}) {
	task.context[name] = value
}

// RemoveFromContext deletes value by key from the task context.
func (task *Task) RemoveFromContext(name string) {
	delete(task.context, name)
}
