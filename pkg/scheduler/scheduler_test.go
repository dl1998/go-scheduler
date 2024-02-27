// Package scheduler_test has tests for the scheduler package.
package scheduler

import (
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"
)

// CreateEmptyScheduler creates a new Scheduler object with default values.
func CreateEmptyScheduler() *Scheduler {
	return New()
}

// CreateSchedulerWithTasks creates a new Scheduler object with two dummy tasks.
func CreateSchedulerWithTasks() *Scheduler {
	newScheduler := New()
	newScheduler.Tasks = []*Task{
		&Task{ID: uuid.New().String(), Name: "Task 1"},
		&Task{ID: uuid.New().String(), Name: "Task 2"},
	}
	return newScheduler
}

// IsMapStringToInterface checks that provided object is of type
// map[string]interface{}.
func IsMapStringToInterface(value interface{}) bool {
	reflectionType := reflect.TypeOf(value)

	isCorrectType := false

	if reflectionType.Kind() == reflect.Map {
		isCorrectType = reflectionType.Key().Kind() == reflect.String && reflectionType.Elem().Kind() == reflect.Interface
	}

	return isCorrectType
}

// CompareMapsStringToInterface checks that two maps of type
// map[string]interface{} are equal.
func CompareMapsStringToInterface(map1 map[string]interface{}, map2 map[string]interface{}) bool {
	isEqual := true

	if len(map1) != len(map2) {
		isEqual = false
	}

	for key, _ := range map1 {
		_, ok := map2[key]

		if !ok {
			isEqual = false
			break
		}
	}

	return isEqual
}

// TestNew tests that New method returns new Scheduler object with Tasks initialized as nil.
func TestNew(t *testing.T) {
	newScheduler := CreateEmptyScheduler()
	if newScheduler.Tasks != nil {
		t.Fatalf("New Scheduler instance has incorrect value: %v.", newScheduler.Tasks)
	}
}

// TestNewTask tests that NewTask method returns a new Task object correctly
// initialized with provided variables.
func TestNewTask(t *testing.T) {
	id := "Custom ID"
	name := "Test Task"
	startTime := time.Date(2000, 01, 01, 01, 02, 03, 0, time.Local)
	duration := 10 * time.Second
	interval := time.Second
	stopSignal := make(chan bool)
	context := map[string]interface{}{
		"one":     1,
		"boolean": true,
		"string":  "Test String",
	}

	newTask := NewTask(id, name, &startTime, &duration, interval, stopSignal, context)

	if newTask.ID != id {
		t.Fatalf("Incorrect default task id. Task ID: %s", newTask.ID)
	}

	if newTask.Name != name {
		t.Fatalf("Incorrect default task name. Task Name: %s", newTask.Name)
	}

	if !newTask.Start.Equal(startTime) {
		t.Fatalf("Incorrect default task start time. Task Start: %s", newTask.Start.String())
	}

	if *newTask.Duration != duration {
		t.Fatalf("Incorrect default task duration. Task Duration: %v", newTask.Duration)
	}

	if newTask.Interval != interval {
		t.Fatalf("Incorrect default task interval. Task Interval: %v", newTask.Interval)
	}

	if newTask.stopSignal != stopSignal {
		t.Fatalf("Incorrect default task stop signal. Task Stop Signal: %v", newTask.stopSignal)
	}

	if !CompareMapsStringToInterface(newTask.context, context) {
		t.Fatalf("Incorrect default task context. Task Context: %v", newTask.context)
	}
}

// TestNewSimpleTask tests that NewSimpleTask method returns a new Task object
// initialized with default values.
func TestNewSimpleTask(t *testing.T) {
	newTask := NewSimpleTask("", time.Second)
	now := time.Now()
	acceptableTimeDelta := 1 * time.Second

	if newTask.ID == "" {
		t.Fatalf("Incorrect default task id. Task ID: %s", newTask.ID)
	}

	if newTask.Name != "" {
		t.Fatalf("Incorrect default task name. Task Name: %s", newTask.Name)
	}

	timeDifference := now.Sub(*newTask.Start)

	if timeDifference > acceptableTimeDelta {
		t.Fatalf("Incorrect default task start time. Task Start: %s", newTask.Start.String())
	}

	if newTask.Duration != nil {
		t.Fatalf("Incorrect default task duration. Task Duration: %v", newTask.Duration)
	}

	if newTask.Interval != time.Second {
		t.Fatalf("Incorrect default task interval. Task Interval: %v", newTask.Interval)
	}

	stopSignalReflectionType := reflect.TypeOf(newTask.stopSignal)

	if stopSignalReflectionType.Kind() != reflect.Chan || stopSignalReflectionType.Elem().Kind() != reflect.Bool {
		t.Fatalf("Incorrect default task stop signal. Task Stop Signal: %v", newTask.stopSignal)
	}

	if !IsMapStringToInterface(newTask.context) {
		t.Fatalf("Incorrect default task context. Task Context: %v", newTask.context)
	}
}

// TestTask_String tests that Task.String method return data in the correct
// format.
func TestTask_String(t *testing.T) {
	id := "Custom ID"
	name := "Test Task"
	startTime := time.Date(2000, 01, 01, 01, 02, 03, 0, time.Local)
	duration := 10 * time.Second
	interval := time.Second
	stopSignal := make(chan bool)
	context := map[string]interface{}{
		"one":     1,
		"boolean": true,
		"string":  "Test String",
	}

	newTask := NewTask(id, name, &startTime, &duration, interval, stopSignal, context)

	expectedString := fmt.Sprintf(
		"ID: %s\nName: %s\nStart: %s\nDuration: %s\nInterval: %s\nContext: %v",
		newTask.ID,
		newTask.Name,
		newTask.Start.String(),
		newTask.Duration.String(),
		newTask.Interval.String(),
		newTask.context,
	)

	if newTask.String() != expectedString {
		t.Fatalf("Incorrect string for the task object. Expected: %s. Actual: %s.", expectedString, newTask.String())
	}
}

// TestScheduler_FindTaskIndex tests that Scheduler.FindTaskIndex method returns
// correct index for the provided Task object.
func TestScheduler_FindTaskIndex(t *testing.T) {
	newScheduler := CreateSchedulerWithTasks()
	for index, task := range newScheduler.Tasks {
		foundIndex := newScheduler.FindTaskIndex(task)
		if foundIndex != index {
			t.Fatalf("Task with ID: \"%s\" was found under incorrect index: %d, correct index is %d.", task.ID, foundIndex, index)
		}
	}
}

// TestScheduler_FindTaskIndex_EmptyList tests that Scheduler.FindTaskIndex
// method returns -1, if Task object was not found on the list of
// Scheduler.Tasks.
func TestScheduler_FindTaskIndex_EmptyList(t *testing.T) {
	newScheduler := CreateEmptyScheduler()
	newTask := &Task{ID: uuid.New().String()}
	foundIndex := newScheduler.FindTaskIndex(newTask)
	if foundIndex != -1 {
		t.Fatalf("Task with ID: \"%s\" was found under index: %d, on the empty list.", newTask.ID, foundIndex)
	}
}

// TestScheduler_FindTaskByName tests that Scheduler.FindTaskByName method
// returns Task object based on the provided task name. It iterates through all
// Scheduler.Tasks and tries to find this Task using its name.
func TestScheduler_FindTaskByName(t *testing.T) {
	newScheduler := CreateSchedulerWithTasks()
	for _, task := range newScheduler.Tasks {
		foundTask := newScheduler.FindTaskByName(task.Name)
		if task != foundTask {
			t.Fatalf("Task with Name: \"%s\" was not found, instead it returned %v.", task.Name, foundTask)
		}
	}
}

// TestScheduler_FindTaskByName_NotExist tests that Scheduler.FindTaskByName method
// returns nil when Task with provided task name doesn't exist.
func TestScheduler_FindTaskByName_NotExist(t *testing.T) {
	newScheduler := CreateSchedulerWithTasks()
	foundTask := newScheduler.FindTaskByName("Incorrect Task Name")
	if foundTask != nil {
		t.Fatalf("Found not existing Task: %v, using FindTaskByName.", foundTask)
	}
}

// TestScheduler_FindTaskByID tests that Scheduler.FindTaskByID method returns
// Task object based on the provided task id. It iterates through all
// Scheduler.Tasks and tries to find this Task using its id.
func TestScheduler_FindTaskByID(t *testing.T) {
	newScheduler := CreateSchedulerWithTasks()
	for _, task := range newScheduler.Tasks {
		foundTask := newScheduler.FindTaskByID(task.ID)
		if task != foundTask {
			t.Fatalf("Task with ID: \"%s\" was not found, instead it returned %v.", task.ID, foundTask)
		}
	}
}

// TestScheduler_FindTaskByID_NotExist tests that Scheduler.FindTaskByID method
// returns nil when Task with provided task id doesn't exist.
func TestScheduler_FindTaskByID_NotExist(t *testing.T) {
	newScheduler := CreateSchedulerWithTasks()
	foundTask := newScheduler.FindTaskByID("Incorrect Task ID")
	if foundTask != nil {
		t.Fatalf("Found not existing Task: %v, using FindTaskByID.", foundTask)
	}
}

// TestScheduler_ScheduleTask tests that Scheduler.ScheduleTask method correctly
// schedules Task and executes it. In this test it uses dummy function that
// increases counter value every execution, at the end it checks that within
// specified Task.Duration it was executed correct number of times.
func TestScheduler_ScheduleTask(t *testing.T) {
	counter := 0
	var testFunction = func(task *Task) {
		counter += 1
	}

	taskName := "Test Task"
	durationSeconds := 5
	duration := time.Duration(durationSeconds) * time.Second

	newScheduler := CreateEmptyScheduler()
	newTask := newScheduler.ScheduleTask(taskName, nil, &duration, 1*time.Second, testFunction)

	time.Sleep(duration)

	if counter != durationSeconds {
		t.Fatalf("Task has been scheduled for %v with %v interval, but it was executed only %v times.", newTask.Duration, newTask.Interval, counter)
	}
}

// TestScheduler_StopTask tests that Scheduler.StopTask method correctly stops
// execution of the scheduled Task. It sends stop signal for the provided task,
// this triggers task termination and clean-up by Scheduler.
func TestScheduler_StopTask(t *testing.T) {
	counter := 0
	var testFunction = func(task *Task) {
		counter += 1
	}

	taskName := "Test Task"
	durationSeconds := 5
	duration := time.Duration(durationSeconds) * time.Second

	newScheduler := CreateEmptyScheduler()
	newTask := newScheduler.ScheduleTask(taskName, nil, nil, 1*time.Second, testFunction)

	time.Sleep(duration)

	err := newScheduler.StopTask(newTask)

	foundTask := false

	for _, task := range newScheduler.Tasks {
		if task == newTask {
			foundTask = true
		}
	}

	if foundTask && err != nil {
		t.Fatalf("Task \"%s\" with id \"%s\" has not been stopped.", newTask.Name, newTask.ID)
	}
}

// TestScheduler_StopTask_NotExist tests that Scheduler.StopTask method correctly
// handles situation where provided task doesn't exist.
func TestScheduler_StopTask_NotExist(t *testing.T) {
	taskName := "Test Task"
	intervalSeconds := 5
	interval := time.Duration(intervalSeconds) * time.Second

	newScheduler := CreateEmptyScheduler()
	newTask := NewSimpleTask(taskName, interval)

	err := newScheduler.StopTask(newTask)
	expected := fmt.Sprintf("task with id: %s cannot be stopped, because it was not found", newTask.ID)

	if err.Error() != expected {
		t.Fatalf("Error has not been thrown for missing task with id %s.", newTask.ID)
	}
}

// TestTask_GetFromContext tests that Task.GetFromContext method returns correct
// value from the Task.context based on the provided key.
func TestTask_GetFromContext(t *testing.T) {
	task := NewTask("", "", nil, nil, time.Second, nil, nil)

	key := "key"
	value := "value"

	task.context[key] = value

	assignedValue := task.GetFromContext(key)

	if assignedValue != value {
		t.Fatalf("Task value has been set incorrectly. Expected: %v. Actual: %v.", value, assignedValue)
	}
}

// TestTask_SetToContext tests that Task.SetToContext method correctly sets a
// value for the provided key in the Task.context.
func TestTask_SetToContext(t *testing.T) {
	task := NewTask("", "", nil, nil, time.Second, nil, nil)

	key := "key"
	value := "value"

	task.SetToContext(key, value)

	assignedValue := task.context[key]

	if assignedValue != value {
		t.Fatalf("Task value has been set incorrectly. Expected: %v. Actual: %v.", value, assignedValue)
	}
}

// TestTask_RemoveFromContext tests that Task.RemoveFromContext method removes
// key-value pair from the context for the provided key.
func TestTask_RemoveFromContext(t *testing.T) {
	task := NewTask("", "", nil, nil, time.Second, nil, nil)

	key := "key"
	value := "value"

	task.context[key] = value

	task.RemoveFromContext(key)

	contextValue := task.context[key]

	if contextValue != nil {
		t.Fatalf("Task value has not been removed from the context. Value: %v.", contextValue)
	}
}
