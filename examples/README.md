# Examples

This package contains examples how to use Scheduler.

## basic

Demonstrates how to schedule the simplest task possible. 
It schedules task that prints "Hello, World!" during 10 second with interval of 1 second.

## interrupt

Demonstrates how to schedule task that prints counter value, where each execution it increases its value by 1. 
After scheduling it waits 5 seconds and interrupts task execution.

## taskcontext

Demonstrates how to schedule a task with data that are shared between executions in the Task.context.
It schedules a task that performs function which reads counter from Task.context,
if counter is not found in the task context it set counter as 0.
After printing "[${counter}] Hello, World!", it increases counter and save it to Task.context.

## withparameters

Demonstrates how to pass argument(s) to the function that will be scheduled.