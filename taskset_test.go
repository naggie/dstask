package dstask

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var TestTasks []*Task = []*Task{
	{
		UUID:    "01234567-8901-2345-6789-101112131415",
		Summary: "This is Test Task 1",
		ID:      1,
	},
	{
		UUID:    "ee4b08e6-2bfe-4265-a358-a3fdb7a27cfb",
		Summary: "This is Test Task 2",
		ID:      2,
	},
	{
		UUID:    "efaa5e7f-8cb0-420c-94c9-98142bd33721",
		Summary: "This is Test Task 3",
		ID:      3,
	},
}
var TestTs *TaskSet = &TaskSet{
	tasks:       TestTasks,
	tasksByID:   make(map[int]*Task),
	tasksByUUID: make(map[string]*Task),
}

func TestSearchForUUID(t *testing.T) {
	for _, task := range TestTasks {
		TestTs.tasksByUUID[task.UUID] = task
		TestTs.tasksByID[task.ID] = task
	}

	type testcase struct {
		input          string
		outputExpected string
		errorExpected  bool
	}

	var tests = []testcase{
		{
			input:          "012",
			outputExpected: "01234567-8901-2345-6789-101112131415",
			errorExpected:  false,
		},
		{
			input:          "ee",
			outputExpected: "ee4b08e6-2bfe-4265-a358-a3fdb7a27cfb",
			errorExpected:  false,
		},
		{
			input:          "ef",
			outputExpected: "efaa5e7f-8cb0-420c-94c9-98142bd33721",
			errorExpected:  false,
		},

		{
			input:          "1e",
			outputExpected: "",
			errorExpected:  true,
		},
		{
			input:          "e",
			outputExpected: "",
			errorExpected:  true,
		},
	}
	for _, tc := range tests {

		output, err := TestTs.SearchForUUID(tc.input)
		//t.Logf("Got %s and error: %v", output, err)
		assert := assert.New(t)
		assert.Equal(tc.outputExpected, output)
		if tc.errorExpected {
			assert.NotNil(err)
		} else {
			assert.Nil(err)
		}

	}
}
func TestMustGetByUUID(t *testing.T) {
	for _, task := range TestTasks {
		TestTs.tasksByUUID[task.UUID] = task
		TestTs.tasksByID[task.ID] = task
	}

	type testcase struct {
		input          string
		outputExpected Task
		errorExpected  bool
	}

	var tests = []testcase{
		{ // Testing partial UUID in Task Set
			input:          "012",
			outputExpected: *TestTasks[0], //Return First Task
			errorExpected:  false,
		},
		{ // Testing partial UUID not in Task Set
			input:          "0012",
			outputExpected: Task{}, //Return First Task
			errorExpected:  true,
		},
		{ // Valid UUID in TaskSet
			input:          "ee4b08e6-2bfe-4265-a358-a3fdb7a27cfb",
			outputExpected: *TestTasks[1], //Return second task.
			errorExpected:  false,
		},
		{ // Changing the last character, valid UUID, but not in TaskSet
			input:          "ee4b08e6-2bfe-4265-a358-a3fdb7a27cf1",
			outputExpected: Task{},
			errorExpected:  true,
		},
	}
	for _, tc := range tests {

		output, err := TestTs.MustGetByUUID(tc.input)
		//t.Logf("Got %s and error: %v \n", output, err)
		assert := assert.New(t)
		assert.Equal(tc.outputExpected, output)
		if tc.errorExpected {
			assert.NotNil(err)
		} else {
			assert.Nil(err)
		}

	}
}
func TestMustGetByID(t *testing.T) {
	for _, task := range TestTasks {
		TestTs.tasksByUUID[task.UUID] = task
		TestTs.tasksByID[task.ID] = task
	}

	type testcase struct {
		input          int
		outputExpected Task
		errorExpected  bool
	}

	var tests = []testcase{
		{ // Choosing ID in Task Set
			input:          1,
			outputExpected: *TestTasks[0], //Return First Task
			errorExpected:  false,
		},
		{ // Choosing ID in Task Set
			input:          2,
			outputExpected: *TestTasks[1], //Return Second Task
			errorExpected:  false,
		},
		{ // Choosing ID in Task Set
			input:          3,
			outputExpected: *TestTasks[2], //Return Third Task
			errorExpected:  false,
		},
		{ // Choose ID not in task set.
			input:          4,
			outputExpected: Task{},
			errorExpected:  true,
		},
	}
	for _, tc := range tests {

		output, err := TestTs.MustGetByID(tc.input)
		//t.Logf("Got %s and error: %v \n", output, err)
		assert := assert.New(t)
		assert.Equal(tc.outputExpected, output)
		if tc.errorExpected {
			assert.NotNil(err)
		} else {
			assert.Nil(err)
		}

	}
}
func TestMustGetTask(t *testing.T) {
	for _, task := range TestTasks {
		TestTs.tasksByUUID[task.UUID] = task
		TestTs.tasksByID[task.ID] = task
	}

	type testcase struct {
		input          interface{}
		outputExpected Task
		errorExpected  bool
	}

	var tests = []testcase{
		{ // Choosing ID in Task Set
			input:          1,
			outputExpected: *TestTasks[0], //Return First Task
			errorExpected:  false,
		},
		{ // Choose ID not in task set.
			input:          4,
			outputExpected: Task{},
			errorExpected:  true,
		},
		{ // Choose Invalid ID type not in task set.
			input:          4.43,
			outputExpected: Task{},
			errorExpected:  true,
		},
		{ // Choose Invalid ID type not in task set.
			input:          []string{"one", "two"},
			outputExpected: Task{},
			errorExpected:  true,
		},
	}
	for _, tc := range tests {

		output, err := TestTs.MustGetTask(tc.input)
		//t.Logf("Got %s and error: %v \n", output, err)
		assert := assert.New(t)
		assert.Equal(tc.outputExpected, output)
		if tc.errorExpected {
			assert.NotNil(err)
		} else {
			assert.Nil(err)
		}

	}
}
