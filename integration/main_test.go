package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/naggie/dstask"
	"gotest.tools/assert"
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	gitConfigWarnIfUnset("user.name")
	gitConfigWarnIfUnset("user.email")

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

func makeDstaskRepo(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := ioutil.TempDir("", "dstask")
	if err != nil {
		t.Fatal()
	}
	// Initialize with a random branch to ensure dstask does not rely on
	// a particular default branch name.
	randomBranch := makeRandomString("branch_", 6)
	_ = randomBranch
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

func makeRandomString(prefix string, length uint) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return fmt.Sprintf("%s%s", prefix, string(b))
}

func gitConfigWarnIfUnset(param string) {
	cmd := exec.Command("git", "config", "--get", param)
	err := cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() == 1 {
			log.Printf("WARNING: git config %s is unset. This will likely cause test failures.", param)
		}
	}
}
