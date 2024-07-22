package locker

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

const LOCK_FILE_PATH = "~/.kdb/.storage_lock"

type FileLock struct {
	file    *os.File
	timeout time.Duration
}

func NewFileLock(timeout time.Duration) (*FileLock, error) {
	file, err := os.OpenFile(LOCK_FILE_PATH, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	return &FileLock{file: file, timeout: timeout}, nil
}

func (fl *FileLock) Lock() error {
	start := time.Now()
	for {
		err := syscall.Flock(int(fl.file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			return nil // Lock acquired successfully
		}
		if err != syscall.EWOULDBLOCK {
			return fmt.Errorf("unexpected error while acquiring lock: %v", err)
		}
		if time.Since(start) > fl.timeout {
			return fmt.Errorf("timeout while acquiring lock")
		}
		time.Sleep(50 * time.Millisecond) // Wait before trying again
	}
}

func (fl *FileLock) Unlock() error {
	return syscall.Flock(int(fl.file.Fd()), syscall.LOCK_UN)
}

func (fl *FileLock) Close() error {
	return fl.file.Close()
}
