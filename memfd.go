package memfd

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var errTooBig = errors.New("[error] memfd too large for slice")

const maxint int64 = int64(^uint(0) >> 1)

const (
	MFD_CREATE  = 319
	MFD_CLOEXEC = 0x0001
)

type MemFD struct {
	*os.File
}

func New(name string) *MemFD {
	fd, _, _ := syscall.Syscall(MFD_CREATE, uintptr(unsafe.Pointer(&name)), uintptr(MFD_CLOEXEC), 0)
	return &MemFD{
		os.NewFile(fd, name),
	}
}

func (self *MemFD) Write(bytes []byte) (int, error) {
	return syscall.Write(int(self.Fd()), bytes)
}

func (self *MemFD) Path() string {
	return fmt.Sprintf("/proc/self/fd/%d", self.Fd())
}

func (self *MemFD) Info() (os.FileInfo, error) {
	return os.Lstat(self.Path())
}

func (self *MemFD) Exec(arguments string) error {
	return syscall.Exec(self.Path(), []string{self.Name(), arguments}, nil)
}
