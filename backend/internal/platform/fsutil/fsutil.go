// Package fsutil holds filesystem helpers.
package fsutil

import (
	"fmt"
	"os"
	"path/filepath"
)

// TempSuffix ends the name of every temp file WriteFileAtomic creates, so
// leftovers from a crash can be recognized and removed.
const TempSuffix = ".tmp"

// WriteFileAtomic writes data to path through a temp file in the same
// directory followed by a rename, so readers and crashes never observe a
// half-written file. The rename is atomic on the same filesystem.
func WriteFileAtomic(path string, data []byte, perm os.FileMode) (err error) {
	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".*"+TempSuffix)
	if err != nil {
		return fmt.Errorf("create temp file for %s: %w", path, err)
	}
	tmpName := tmp.Name()
	defer func() {
		if err != nil {
			_ = tmp.Close()        // may already be closed; the original error matters more
			_ = os.Remove(tmpName) // best effort; leftovers are removed at startup
		}
	}()

	if _, err = tmp.Write(data); err != nil {
		return fmt.Errorf("write temp file for %s: %w", path, err)
	}
	if err = tmp.Sync(); err != nil {
		return fmt.Errorf("sync temp file for %s: %w", path, err)
	}
	if err = tmp.Close(); err != nil {
		return fmt.Errorf("close temp file for %s: %w", path, err)
	}
	if err = os.Chmod(tmpName, perm); err != nil {
		return fmt.Errorf("chmod temp file for %s: %w", path, err)
	}
	if err = os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("rename temp file to %s: %w", path, err)
	}
	return nil
}
