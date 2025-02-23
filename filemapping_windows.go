package gettext

import (
	"errors"
	"io/fs"
	"os"
	"syscall"
	"unsafe"
)

// Adapted from https://github.com/golang/exp/blob/master/mmap/mmap_windows.go

func (m *fileMapping) tryMap(f fs.File, size int64) error {
	of, ok := f.(*os.File)
	if !ok {
		return errors.New("virtual filesystem doesn't support mmap")
	}
	low, high := uint32(size), uint32(size>>32)
	fmap, err := syscall.CreateFileMapping(syscall.Handle(of.Fd()), nil, syscall.PAGE_READONLY, high, low, nil)
	if err != nil {
		return err
	}
	defer syscall.CloseHandle(fmap)

	ptr, err := syscall.MapViewOfFile(fmap, syscall.FILE_MAP_READ|syscall.FILE_MAP_COPY, 0, 0, uintptr(size))
	if err != nil {
		return err
	}
	m.data = (*[1<<31 - 1]byte)(unsafe.Pointer(ptr))[:size]
	m.isMapped = true
	return nil
}

func (m *fileMapping) closeMapping() error {
	return syscall.UnmapViewOfFile(uintptr(unsafe.Pointer(&m.data[0])))
}
