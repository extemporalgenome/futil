// Copyright 2012 Kevin Gillette. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package futil

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestAtomic(t *testing.T) {
	name := "delete-this"
	data := "X"
	f, err := CreateAtomic(name)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(name)
	defer f.Destroy()
	defer f.File.Close()
	info1, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("temporary name:", f.Name())
	f.WriteString(data)
	err = f.Close()
	if err != nil {
		t.Fatal(err)
	}
	info2, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(info1, info2) {
		t.Fatal("file mismatch")
	}
	f, err = OpenAtomic(name)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Destroy()
	defer f.File.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	} else if string(b) != data {
		t.Fatalf("data mismatch: expected %q but found %q", data, b)
	}
}
