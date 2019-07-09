// Copyright (c) 2012, Suryandaru Triandana <ztteng@gmail.com>
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd

package storage

import (
	"os"
	"syscall"
)

type unixFileLock struct {
	f *os.File
}

func (fl *unixFileLock) release() error {
	if err := setFileLock(fl.f, false, false); err != nil {
		return err
	}
	return fl.f.Close()
}

func newFileLock(path string, readOnly bool) (fl fileLock, err error) {
	var flag int
	if readOnly {
		flag = os.O_RDONLY
	} else {
		flag = os.O_RDWR
	}
	f, err := os.OpenFile(path, flag, 0)
	if os.IsNotExist(err) {
		f, err = os.OpenFile(path, flag|os.O_CREATE, 0644)
	}
	if err != nil {
		return
	}
	err = setFileLock(f, readOnly, true)
	if err != nil {
		f.Close()
		return
	}
	fl = &unixFileLock{f: f}
	return
}

//LOCK_SH 建立共享锁定。多个进程可同时对同一个文件作共享锁定。
//LOCK_EX 建立互斥锁定。一个文件同时只有一个互斥锁定。
//LOCK_UN 解除文件锁定状态
//LOCK_NB 无法建立锁定时，此操作可不被阻断，马上返回进程。通常与LOCK_SH或LOCK_EX 做OR(|)组合
func setFileLock(f *os.File, readOnly, lock bool) error {
	how := syscall.LOCK_UN
	if lock {
		if readOnly {
			how = syscall.LOCK_SH
		} else {
			how = syscall.LOCK_EX
		}
	}
	return syscall.Flock(int(f.Fd()), how|syscall.LOCK_NB)
}

func rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func isErrInvalid(err error) bool {
	if err == os.ErrInvalid {
		return true
	}

	// Go >= 1.8 returns *os.PathError instead
	if patherr, ok := err.(*os.PathError); ok && patherr.Err == syscall.EINVAL {
		return true
	}
	return false
}

//fsync()负责将参数fd 所指的文件数据, 由系统缓冲区写回磁盘, 以确保数据同步
func syncDir(name string) error {
	// As per fsync manpage, Linux seems to expect fsync on directory, however
	// some system don't support this, so we will ignore syscall.EINVAL.
	//
	// From fsync(2):
	//   Calling fsync() does not necessarily ensure that the entry in the
	//   directory containing the file has also reached disk. For that an
	//   explicit fsync() on a file descriptor for the directory is also needed.
	//打开文件后记得 直接defer 关闭文件文件
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := f.Sync(); err != nil && !isErrInvalid(err) {
		return err
	}
	return nil
}
