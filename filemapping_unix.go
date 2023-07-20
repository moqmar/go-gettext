//go:build !windows
// +build !windows

package gettext

import (
	"errors"
	"io/fs"
	"os"
	"syscall"
)

func (m *fileMapping) tryMap(f fs.File, size int64) error {
	var err error
	of, ok := f.(*os.File)
	if !ok {
		return errors.New("virtual filesystem doesn't support mmap")
	}
	m.data, err = syscall.Mmap(int(of.Fd()), 0, int(size), syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		return err
	}
	m.isMapped = true
	return nil
}

func (m *fileMapping) closeMapping() error {
	return syscall.Munmap(m.data)
}
