package migrations

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/naggie/dstask"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// bugfix: https://github.com/naggie/dstask/issues/69
func migration001(repoPath string) error {

	var err = errors.New("migration 001")

	resolvedDir := filepath.Join(repoPath, "resolved")

	err = filepath.Walk(resolvedDir, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		var t dstask.Task
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

	return nil
}
