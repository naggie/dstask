package dstask

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// RunMigrations executes migrations.
func RunMigrations(repoPath string) error {
	return migration001(repoPath)
}

// bugfix: https://github.com/naggie/dstask/issues/69
func migration001(repoPath string) error {

	var err = errors.New("migration 001")

	resolvedDir := filepath.Join(repoPath, "resolved")

	// Perform the migration changes on disk.
	err = filepath.Walk(resolvedDir, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		var t Task
		err = yaml.Unmarshal(data, &t)
		if err != nil {
			return errors.Wrapf(err, "malformed task file: %s", info.Name())
		}

		var zeroTime time.Time
		if t.Resolved == zeroTime {
			t.Resolved = info.ModTime()
			fixed, err := yaml.Marshal(&t)
			if err != nil {
				return errors.Wrapf(err, "could not marshal task %s", t.UUID)
			}
			err = ioutil.WriteFile(path, fixed, info.Mode())
			if err != nil {
				return errors.Wrapf(err, "write file: %s", path)
			}

		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "error from filepath.Walk")
	}

	// If no errors occurred, commit the result.
	commitMsg := `Database migration 001

This is a fix for issue https://github.com/naggie/dstask/issues/69

Tasks marked resolved were not getting their Resolved timestamp set.
Resolved timestamps were incorrectly set to Go's time.Time zero value.
This migration makes the assumption that the filesystem modtime is
the resolved time, and sets it accordingly.
`

	MustGitCommit(repoPath, commitMsg)

	return nil
}
