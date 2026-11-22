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

func normalizeToLocal(t time.Time) time.Time {
    return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

func assertDateEqual(t *testing.T, expected, actual time.Time) {
    expectedLocal := normalizeToLocal(expected)
    actualLocal := normalizeToLocal(actual)
    assert.Equal(t, expectedLocal, actualLocal)
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
	assertDateEqual(t, getTestDate(2025, 7, 1), tasks[0].Due)
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
	assertDateEqual(t, getTestDate(currentYear, 7, 1), tasks[0].Due)
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
	assertDateEqual(t, getTestDate(currentYear, currentMonth, 15), tasks[0].Due)
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
	assertDateEqual(t, getCurrentDate(), tasks[0].Due)
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
	assertDateEqual(t, getRelativeDate(-1), tasks[0].Due)
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
	assertDateEqual(t, getRelativeDate(1), tasks[0].Due)
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
		assertDateEqual(t, expectedDate, tasks[0].Due)
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
	assertDateEqual(t, getTestDate(2025, 7, 1), tasks[0].Due)
	assert.Equal(t, "Task 1", tasks[0].Summary)
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
	assertDateEqual(t, getCurrentDate(), tasks[0].Due)
	assert.Equal(t, "Task due today", tasks[0].Summary)
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
	assertDateEqual(t, getCurrentDate(), tasks[0].Due)
	assertDateEqual(t, getRelativeDate(1), tasks[1].Due)
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
	assertDateEqual(t, getRelativeDate(-1), tasks[0].Due)
	assert.Equal(t, "Task 2", tasks[1].Summary)
	assertDateEqual(t, getCurrentDate(), tasks[1].Due)
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
	assertDateEqual(t, getCurrentDate(), tasks[0].Due)
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
	assertDateEqual(t, getTestDate(2025, time.June, 18), tasks[0].Due)
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
	assertDateEqual(t, getTestDate(2025, time.October, 31), tasks[0].Due)
}

func TestTemplateInDueNoneContext(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("template", "Template 1", "due:2025-10-31")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("context", "due:none")
	assertProgramResult(t, output, exiterr, success)

	// due:none in the context is a filter and keeps the due date of the template
	output, exiterr, success = program("add", "template:1", "task from template")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "--", "due:2025-10-31")
	assertProgramResult(t, output, exiterr, success)

	assert.Equal(t, []string{"task from template"}, taskSummaries(unmarshalTaskArray(t, output)))
}

func TestDueOnInDueContext(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task", "due:2025-09-01")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("context", "due:2025-09-01")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "due.on:2025-09-01")
	assertProgramResult(t, output, exiterr, success)

	assert.Equal(t, []string{"Task"}, taskSummaries(unmarshalTaskArray(t, output)))
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
	assertDateEqual(t, getTestDate(2025, time.September, 1), tasks[0].Due)
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
	assertDateEqual(t, getCurrentDate(), taskWithDue.Due)
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
	assertDateEqual(t, getCurrentDate(), resolvedTasks[0].Due)
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
	assertDateEqual(t, getCurrentDate(), tasks[0].Due)
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

func dueArg(date time.Time) string {
	return "due:" + date.Format("2006-01-02")
}

func TestFilterTasksByDueNone(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task with due", "due:today")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Task without due")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "due:none")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Equal(t, []string{"Task without due"}, taskSummaries(tasks))
}

func TestFilterTasksByDueAny(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task with due", "due:today")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Task without due")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "due:any")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Equal(t, []string{"Task with due"}, taskSummaries(tasks))
}

func TestDueNoneAsContext(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task with due", "due:today")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Task without due")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("context", "due:none")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Equal(t, []string{"Task without due"}, taskSummaries(tasks))

	// new tasks in a due:none context get no due date
	output, exiterr, success = program("add", "New task")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks = unmarshalTaskArray(t, output)
	assert.ElementsMatch(t, []string{"Task without due", "New task"}, taskSummaries(tasks))

	// conflicting due filter
	_, _, success = program("next", "due:today")
	assert.False(t, success)
}

func TestAddWithDueNoneFails(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	_, _, success := program("add", "Task", "due:none")
	assert.False(t, success)

	_, _, success = program("add", "Task", "due:any")
	assert.False(t, success)
}

func TestModifyRemovesDueDate(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "Task 1", "due:today")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("modify", "1", "due:none")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Len(t, tasks, 1)
	assert.True(t, tasks[0].Due.IsZero())

	_, _, success = program("modify", "1", "due:any")
	assert.False(t, success)
}

func TestFilterTasksDueRange(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	dates := []string{"2025-06-30", "2025-07-01", "2025-07-15", "2025-07-31", "2025-08-01"}
	for _, date := range dates {
		output, exiterr, success := program("add", "Task "+date, "due:"+date)
		assertProgramResult(t, output, exiterr, success)
	}

	output, exiterr, success := program("add", "Task without due")
	assertProgramResult(t, output, exiterr, success)

	expected := []string{"Task 2025-07-01", "Task 2025-07-15", "Task 2025-07-31"}

	output, exiterr, success = program("next", "due.after:2025-07-01", "due.before:2025-07-31")
	assertProgramResult(t, output, exiterr, success)
	assert.ElementsMatch(t, expected, taskSummaries(unmarshalTaskArray(t, output)))

	// order does not matter
	output, exiterr, success = program("next", "due.before:2025-07-31", "due.after:2025-07-01")
	assertProgramResult(t, output, exiterr, success)
	assert.ElementsMatch(t, expected, taskSummaries(unmarshalTaskArray(t, output)))

	// range as context
	output, exiterr, success = program("context", "due.after:2025-07-01", "due.before:2025-07-31")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)
	assert.ElementsMatch(t, expected, taskSummaries(unmarshalTaskArray(t, output)))
}

func TestInvalidDueRanges(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	invalid := [][]string{
		{"next", "due.before:2025-07-01", "due.before:2025-07-31"},
		{"next", "due.after:2025-07-01", "due.after:2025-07-31"},
		{"next", "due.after:2025-07-01", "due.before:2025-07-31", "due.before:2025-08-31"},
		{"next", "due:none", "due.before:2025-07-31"},
		{"next", "due:this-week", "due.before:2025-07-31"},
		{"next", "due.before:none"},
		{"add", "Task", "due.after:2025-07-01", "due.before:2025-07-31"},
		{"add", "Task", "due:this-week"},
		{"next", "due.after:2025-07-31", "due.before:2025-07-01"},
		{"next", "due.before:this-week"},
		{"modify", "1", "due:this-week"},
		{"modify", "1", "due.after:2025-07-01", "due.before:2025-07-31"},
		{"log", "Task", "due:this-week"},
		{"template", "Task", "due.after:2025-07-01", "due.before:2025-07-31"},
	}

	output, exiterr, success := program("add", "Task 1")
	assertProgramResult(t, output, exiterr, success)

	for _, args := range invalid {
		_, _, success := program(args...)
		assert.False(t, success, args)
	}
}

func startOfWeek() time.Time {
	today := getCurrentDate()
	return today.AddDate(0, 0, -((int(today.Weekday()) + 6) % 7))
}

func assertDuePeriod(t *testing.T, period string, start, end time.Time) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	for _, date := range []time.Time{start.AddDate(0, 0, -1), start, end, end.AddDate(0, 0, 1)} {
		output, exiterr, success := program("add", "Task "+date.Format("2006-01-02"), dueArg(date))
		assertProgramResult(t, output, exiterr, success)
	}

	output, exiterr, success := program("next", "due:"+period)
	assertProgramResult(t, output, exiterr, success)

	assert.ElementsMatch(t,
		[]string{"Task " + start.Format("2006-01-02"), "Task " + end.Format("2006-01-02")},
		taskSummaries(unmarshalTaskArray(t, output)))
}

func TestFilterTasksThisWeek(t *testing.T) {
	monday := startOfWeek()
	assertDuePeriod(t, "this-week", monday, monday.AddDate(0, 0, 6))
}

func TestFilterTasksNextWeek(t *testing.T) {
	monday := startOfWeek().AddDate(0, 0, 7)
	assertDuePeriod(t, "next-week", monday, monday.AddDate(0, 0, 6))
}

func TestFilterTasksThisMonth(t *testing.T) {
	today := getCurrentDate()
	first := getTestDate(today.Year(), today.Month(), 1)
	assertDuePeriod(t, "this-month", first, first.AddDate(0, 1, -1))
}

func TestFilterTasksNextMonth(t *testing.T) {
	today := getCurrentDate()
	first := getTestDate(today.Year(), today.Month()+1, 1)
	assertDuePeriod(t, "next-month", first, first.AddDate(0, 1, -1))
}

func TestDueRangeFromContextAndQuery(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	for _, date := range []string{"2025-06-30", "2025-07-15", "2025-08-01"} {
		output, exiterr, success := program("add", "Task "+date, "due:"+date)
		assertProgramResult(t, output, exiterr, success)
	}

	output, exiterr, success := program("context", "due.after:2025-07-01")
	assertProgramResult(t, output, exiterr, success)

	// due.after from the context and due.before from the query form a range
	output, exiterr, success = program("next", "due.before:2025-07-31")
	assertProgramResult(t, output, exiterr, success)
	assert.Equal(t, []string{"Task 2025-07-15"}, taskSummaries(unmarshalTaskArray(t, output)))

	// due.after and due.before on the same date form a one-day range
	_, _, success = program("next", "due.before:2025-07-01")
	assert.True(t, success)

	// same date, but another filter type is a conflict
	_, _, success = program("next", "due:2025-07-01")
	assert.False(t, success)
}

func TestAddInDueRangeContext(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("context", "due:this-week")
	assertProgramResult(t, output, exiterr, success)

	// a range is a filter only and sets no due date
	output, exiterr, success = program("add", "Task 1")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "--")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Len(t, tasks, 1)
	assert.True(t, tasks[0].Due.IsZero())
}

func TestAddWithDueInDueAnyContext(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("context", "due:any")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "Task 1", "due:tomorrow")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Len(t, tasks, 1)
	assertDateEqual(t, getRelativeDate(1), tasks[0].Due)
}
