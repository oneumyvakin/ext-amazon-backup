// Copyright 1999-2016. Parallels IP Holdings GmbH.

package main

import "syscall"

func getFreeDiskSpaceInPath(path string) (uint64, error) {
	var stat syscall.Statfs_t

	err := syscall.Statfs(path, &stat)
	if err != nil {
		return 0, err
	}

	return stat.Bavail * uint64(stat.Bsize), err
}
