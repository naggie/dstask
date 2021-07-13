package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/naggie/dstask"
	"gotest.tools/assert"
)

func TestMain(m *testing.M) {
	if err := compile(); err != nil {
		log.Fatalf("compile error: %v", err)
	}
	cleanup := func() {
		if err := os.Remove("dstask"); err != nil {
			log.Panic("could not remove integration test binary")
		}
	}

	exitCode := m.Run()
	cleanup()
	os.Exit(exitCode)
}

func compile() error {
	// We expect to execute in the ./integration directory, and we will output
	// our test binary there.
	cmd := exec.Command("go", "build", "-mod=vendor", "-o", "./dstask", "../cmd/dstask/main.go")
	return cmd.Run()
}

// Create a callable closure that will run our test binary against a
// particular repository path. Any variables set in the environment will be
// passed to the test subprocess.
func testCmd(repoPath string) func(args ...string) ([]byte, *exec.ExitError, bool) {
	return func(args ...string) ([]byte, *exec.ExitError, bool) {
		cmd := exec.Command("./dstask", args...)
		env := os.Environ()
		cmd.Env = append(env, fmt.Sprintf("DSTASK_GIT_REPO=%s", repoPath))
		output, err := cmd.Output()
		exitErr, ok := err.(*exec.ExitError)
		if !ok {
			if err != nil {
				return output, nil, false
			}
			return output, nil, true
		}
		return output, exitErr, exitErr.Success()
	}
}

// Sets an environment variable, and returns a callable closure to unset it.
func setEnv(key, value string) func() {
	if err := os.Setenv(key, value); err != nil {
		panic(err)
	}
	return func() {
		if err := os.Unsetenv(key); err != nil {
			panic(err)
		}
	}
}

func logFailure(t *testing.T, output []byte, exiterr *exec.ExitError) {
	t.Helper()
	t.Logf("stdout: %s", string(output))
	t.Logf("stderr: %v", string(exiterr.Stderr))
}

func unmarshalTaskArray(t *testing.T, data []byte) []dstask.Task {
	t.Helper()
	var tasks []dstask.Task
	err := json.Unmarshal(data, &tasks)
	assert.NilError(t, err)
	return tasks
}

func unmarshalProjectArray(t *testing.T, data []byte) []dstask.Project {
	t.Helper()
	var projects []dstask.Project
	err := json.Unmarshal(data, &projects)
	assert.NilError(t, err)
	return projects
}

func makeDstaskRepo(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := ioutil.TempDir("", "dstask")
	if err != nil {
		t.Fatal()
	}
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatal()
	}

	cleanup := func() {
		os.RemoveAll(dir)
	}
	return dir, cleanup
}

func assertProgramResult(t *testing.T, output []byte, exiterr *exec.ExitError, successExpected bool) {
	t.Helper()
	if exiterr != nil || !successExpected {
		logFailure(t, output, exiterr)
		t.Fatalf("%v", exiterr)
	}
}
