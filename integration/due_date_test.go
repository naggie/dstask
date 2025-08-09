package integration

import (
	"strconv"
	"testing"
	"time"

	"github.com/naggie/dstask"
	"github.com/stretchr/testify/assert"
)

func getTestDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func getCurrentDate() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func getRelativeDate(days int) time.Time {
	now := time.Now().AddDate(0, 0, days)
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func getNextWeekday(weekday time.Weekday) time.Time {
	now := time.Now()
	daysUntil := int(weekday - now.Weekday())
	if daysUntil < 0 {
		daysUntil += 7
	}
	nextWeekday := now.AddDate(0, 0, daysUntil)
	year, month, day := nextWeekday.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func TestAddTaskWithFullDate(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task with full date", "due:2025-07-01")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Equal(t, getTestDate(2025, 7, 1), tasks[0].Due)
	assert.Equal(t, "Task with full date", tasks[0].Summary)
}

func TestAddTaskWithMonthDay(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task with month-day", "due:07-01")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	currentYear := time.Now().Year()
	assert.Equal(t, getTestDate(currentYear, 7, 1), tasks[0].Due)
	assert.Equal(t, "Task with month-day", tasks[0].Summary)
}

func TestAddTaskWithDay(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task with day only", "due:15")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	assert.Equal(t, getTestDate(currentYear, currentMonth, 15), tasks[0].Due)
	assert.Equal(t, "Task with day only", tasks[0].Summary)
}

func TestAddTaskWithToday(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task due today", "due:today")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Equal(t, getCurrentDate(), tasks[0].Due)
	assert.Equal(t, "Task due today", tasks[0].Summary)
}

func TestAddTaskWithYesterday(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task due yesterday", "due:yesterday")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Equal(t, getRelativeDate(-1), tasks[0].Due)
	assert.Equal(t, "Task due yesterday", tasks[0].Summary)
}

func TestAddTaskWithTomorrow(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task due tomorrow", "due:tomorrow")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Equal(t, getRelativeDate(1), tasks[0].Due)
	assert.Equal(t, "Task due tomorrow", tasks[0].Summary)
}

func TestAddTaskWithWeekdays(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	weekdays := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}
	weekdayTimes := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}

	for i, weekday := range weekdays {
		output, exiterr, success := program("add", "Task due "+weekday, "due:"+weekday)
		assertProgramResult(t, output, exiterr, success)

		output, exiterr, success = program("next", "due:"+weekday)
		assertProgramResult(t, output, exiterr, success)

		tasks := unmarshalTaskArray(t, output)
		expectedDate := getNextWeekday(weekdayTimes[i])
		assert.Equal(t, expectedDate, tasks[0].Due)
		assert.Equal(t, "Task due "+weekday, tasks[0].Summary)
	}
}

func TestFilterTasksByExactDate(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task 1", "due:2025-07-01")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Task 2", "due:2025-08-01")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Task 3")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "due:2025-07-01")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "Task 1", tasks[0].Summary)
	assert.Equal(t, getTestDate(2025, 7, 1), tasks[0].Due)
}

func TestFilterTasksByToday(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task due today", "due:today")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Task due tomorrow", "due:tomorrow")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "due:today")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "Task due today", tasks[0].Summary)
	assert.Equal(t, getCurrentDate(), tasks[0].Due)
}

func TestFilterTasksByOverdue(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Overdue task", "due:yesterday")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Today task", "due:today")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Future task", "due:tomorrow")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "due:overdue")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Len(t, tasks, 2)
	assert.Equal(t, "Overdue task", tasks[0].Summary)
	assert.Equal(t, "Today task", tasks[1].Summary)
}

func TestFilterTasksByThisWeekdays(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "This Monday task", "due:this-monday")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Next Monday task", "due:this-friday")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "due:this-monday")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "This Monday task", tasks[0].Summary)
}

func TestFilterTasksByNextWeekdays(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Next Monday task", "due:next-monday")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "This Monday task", "due:next-friday")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "due:next-monday")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "Next Monday task", tasks[0].Summary)
}

func TestFilterTasksDueAfter(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task 1", "due:yesterday")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Task 2", "due:today")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Task 3", "due:tomorrow")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "due.after:today")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Len(t, tasks, 2)
	assert.Equal(t, "Task 2", tasks[0].Summary)
	assert.Equal(t, "Task 3", tasks[1].Summary)
	assert.Equal(t, getCurrentDate(), tasks[0].Due)
	assert.Equal(t, getRelativeDate(1), tasks[1].Due)
}

func TestFilterTasksDueBefore(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task 1", "due:yesterday")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Task 2", "due:today")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Task 3", "due:tomorrow")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "due.before:today")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Len(t, tasks, 2)
	assert.Equal(t, "Task 1", tasks[0].Summary)
	assert.Equal(t, getRelativeDate(-1), tasks[0].Due)
	assert.Equal(t, "Task 2", tasks[1].Summary)
	assert.Equal(t, getCurrentDate(), tasks[1].Due)
}

func TestFilterTasksDueOn(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task 1", "due:yesterday")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Task 2", "due:today")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Task 3", "due:tomorrow")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "due.on:today")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "Task 2", tasks[0].Summary)
	assert.Equal(t, getCurrentDate(), tasks[0].Due)
}

func TestFilterTasksDueAfterWithFullDate(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task 1", "due:2025-06-01")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Task 2", "due:2025-07-01")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Task 3", "due:2025-08-01")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "due.after:2025-06-15")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Len(t, tasks, 2)
	assert.Equal(t, "Task 2", tasks[0].Summary)
	assert.Equal(t, "Task 3", tasks[1].Summary)
}

func TestModifyCommandWithDueDates(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task 1", "due:2025-06-01")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("modify", "1", "due:2025-06-18")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "due:2025-06-18")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)

	assert.Len(t, tasks, 1)
	assert.Equal(t, getTestDate(2025, time.June, 18), tasks[0].Due)
}

func TestTemplatesWithDueDates(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("template", "Template 1", "due:2025-10-31")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "template:1", "task with due date from template")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "due:2025-10-31")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)

	assert.Len(t, tasks, 1)
	assert.Equal(t, getTestDate(2025, time.October, 31), tasks[0].Due)
}

func TestDueDatesMergeWithContext(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("context", "due:2025-09-01", "+work")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "new task with context")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "due:2025-09-01")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)

	assert.Len(t, tasks, 1)
	assert.Equal(t, getTestDate(2025, time.September, 1), tasks[0].Due)
	assert.Equal(t, "new task with context", tasks[0].Summary)
	assert.Equal(t, "work", tasks[0].Tags[0])
}

func TestNextCommandShowsDueDates(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task without due date")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Task with due date", "due:today")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Len(t, tasks, 2)

	var taskWithDue *dstask.Task
	for _, task := range tasks {
		if task.Summary == "Task with due date" {
			taskWithDue = &task
			break
		}
	}
	assert.NotNil(t, taskWithDue)
	assert.Equal(t, getCurrentDate(), taskWithDue.Due)
}

func TestShowResolvedDisplaysDueDates(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Completed task", "due:today")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	taskID := tasks[0].ID

	output, exiterr, success = program("done", strconv.Itoa(taskID))
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("show-resolved")
	assertProgramResult(t, output, exiterr, success)

	resolvedTasks := unmarshalTaskArray(t, output)
	assert.Len(t, resolvedTasks, 1)
	assert.Equal(t, "Completed task", resolvedTasks[0].Summary)
	assert.Equal(t, getCurrentDate(), resolvedTasks[0].Due)
}

func TestInvalidDateFormats(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	invalidFormats := []string{
		"due:invalid-date",
		"due:13-32",
		"due:2025-13-01",
		"due:2025-02-30",
		"due:32",
		"due:next-funday",
		"due:this-xyz",
		"due.afber:today",
	}

	failedCount := 0
	for _, format := range invalidFormats {
		_, _, success := program("add", "Task with invalid date", format)
		if !success {
			failedCount++
		}
	}
	assert.Equal(t, len(invalidFormats), failedCount)
}

func TestCaseInsensitiveDueKeywords(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	caseVariations := []string{"TODAY", "Today", "TOMORROW", "Tomorrow", "MONDAY", "Monday"}

	failedCount := 0
	for _, variation := range caseVariations {
		_, _, success := program("add", "Task with "+variation, "due:"+variation)
		if !success {
			failedCount++
		}
	}

	output, exiterr, success := program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)

	for _, task := range tasks {
		assert.NotNil(t, task.Due)
	}
	assert.Equal(t, 0, failedCount)
}

func TestCombinedDueFilters(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task 1", "due:today", "+urgent")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Task 2", "due:tomorrow", "+urgent")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Task 3", "due:today", "+normal")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "due:today", "+urgent")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "Task 1", tasks[0].Summary)
	assert.Equal(t, getCurrentDate(), tasks[0].Due)
}

func TestMultipleDueDates(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	_, _, success := program("add", "Task with multiple due dates", "due:today", "due:tomorrow")
	assert.False(t, success)
}

func TestAddMultipleDueDatesWithContext(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("context", "due:today", "+urgent")
	assertProgramResult(t, output, exiterr, success)

	_, _, success = program("add", "Task 1", "due:tomorrow")
	assert.False(t, success)
}
