// Copyright 1999-2016. Parallels IP Holdings GmbH.

package main

import (
	"syscall"
	"unsafe"
)

func getFreeDiskSpaceInPath(path string) (uint64, error) {
	h := syscall.MustLoadDLL("kernel32.dll")
	c := h.MustFindProc("GetDiskFreeSpaceExW")

	var freeBytes int64
	var totalBytes int64
	var availBytes int64

	r1, _, err := c.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&freeBytes)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&availBytes)),
	)

	if r1 == 0 {
		return 0, err
	}

	return uint64(freeBytes), nil
}
