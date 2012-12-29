// Copyright 2012 Kevin Gillette. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package futil

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// Atomic behaves like os.File, except that it normally uses a
// temporary file until Close is called, after which the file will be
// atomically renamed to the value of Dest.
type Atomic struct {
	*os.File
	Dest string
}

// Close closes the underlying file descriptor and atomically renames
// the working copy in place of Dest.
func (f *Atomic) Close() error {
	// windows might error out on moving an open file
	err := f.File.Close()
	if err := os.Rename(f.Name(), f.Dest); err != nil {
		return err
	}
	return err
}

// Destroy unlinks the temporary copy.
func (f *Atomic) Destroy() error {
	return os.Remove(f.Name())
}

// CreateAtomic will create an empty temporary file which will be
// renamed to name when closed.
func CreateAtomic(name string) (*Atomic, error) {
	dir, file := filepath.Split(name)
	if dir == "" {
		dir = "."
	}
	f, err := ioutil.TempFile(dir, "."+file+"~")
	if err != nil {
		return nil, err
	}
	return &Atomic{f, name}, nil
}

// OpenAtomic creates a temporary copy of name which, when closed,
// will be renamed back over name.
func OpenAtomic(name string) (*Atomic, error) {
	f, err := CreateAtomic(name)
	if err != nil {
		return nil, err
	}
	if err = CopyFile(f.Name(), name); err != nil {
		f.File.Close()
		f.Destroy()
		return nil, err
	}
	return f, nil
}
