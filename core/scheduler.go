package core

import (
	"codnect.io/chrono"
	"context"
	"time"
)

type Task struct {
	Function func(ctx context.Context, args ...interface{}) interface{} // Task function
	Args     []interface{}
}

// ScheduleTask schedules a task with its arguments.
func ScheduleTask(task Task, scheduler chrono.TaskScheduler, startTime time.Time) (chrono.ScheduledTask, error) {
	return scheduler.Schedule(func(ctx context.Context) {
		// Execute the task function with its arguments
		task.Function(ctx, task.Args...)
	}, chrono.WithTime(startTime))
}

// chrono.WithStartTime(now.Year(), now.Month(), now.Day(), 18, 45, 0),
//chrono.WithLocation("America/New_York"))

func ScheduleAtFixedRate(task Task, scheduler chrono.TaskScheduler, period time.Duration) (chrono.ScheduledTask, error) {
	return scheduler.ScheduleAtFixedRate(func(ctx context.Context) {
		// Execute the task function with its arguments
		task.Function(ctx, task.Args...)
	}, period,
	)
}

func ScheduleAtFixedDelay(task Task, scheduler chrono.TaskScheduler, period time.Duration) (chrono.ScheduledTask, error) {
	return scheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		// Execute the task function with its arguments
		task.Function(ctx, task.Args...)
	}, period,
	)
}

//chrono.WithLocation("America/New_York")

func ScheduleWithCronExp(task Task, scheduler chrono.TaskScheduler, expression string) (chrono.ScheduledTask, error) {
	return scheduler.ScheduleWithCron(func(ctx context.Context) {
		// Execute the task function with its arguments
		task.Function(ctx, task.Args...)
	}, expression,
	)
}

func cancelTask(task chrono.ScheduledTask) {
	task.Cancel()
}

// takes in 2 arguments first:int and second:int
