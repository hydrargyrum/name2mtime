// SPDX-License-Identifier: WTFPL

package main

import "os"
import "os/exec"
import "path"
import "testing"
import "time"

func TestApp(t *testing.T) {
	names := []string{
		"IMG_20220327_123456.jpg",
		"foo-2022-03-27-12-34-56.zip",
	}

	expected := time.Date(2022, 3, 27, 12, 34, 56, 0, time.Local)

	for _, filename := range names {
		filepath := path.Join(t.TempDir(), filename)
		fp, err := os.Create(filepath)
		defer fp.Close()

		if err != nil {
			t.Errorf("could not create %q", filepath)
			break
		}

		cmd := exec.Command("./name2mtime", filepath)
		if err := cmd.Run(); err != nil {
			t.Errorf("failed to run %s", cmd)
			break
		}

		stat, err := os.Stat(filepath)
		if err != nil {
			t.Errorf("failed to stat %q", filepath)
			break
		}

		got := stat.ModTime()
		if got != expected {
			t.Errorf("for %q, expected %s, got %s", filepath, expected, got)
			break
		}
	}
}
