// Copyright (c) 2013, Rob Thornton
// All rights reserved.
// This software is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package token

import "fmt"

type Error struct {
	pos Pos
	msg string
}

type File struct {
	base  Pos     // Base Pos of file
	errs  []Error // List of errors which have occured in processing
	lines []int   // Location of each line ending ('\n')
	name  string  // Filename
	size  int     // Length of file
}

// In the future, will take a FileSet as an argument
func NewFile(name, str string) *File {
	f := new(File)
	f.base = 1 // will be retrieved by FileSet.Base() in future
	f.name = name
	f.size = len(str)
	return f
}

func (f *File) AddLine(off int) {
	f.lines = append(f.lines, off)
}

func (f *File) AddError(p Pos, args ...interface{}) {
	if f.ValidPos(p) {
		f.errs = append(f.errs, Error{p, fmt.Sprint(args...)})
	} else {
		panic("Invalid Position!") // this a little extreme?
	}
}

func (f *File) Base() Pos {
	return f.base
}

func (f *File) NumErrors() int {
	return len(f.errs)
}

func (f *File) PrintError(e Error) {
	var i, line, column int
	for i = 0; i < len(f.lines); i++ {
		if int(e.pos) < f.lines[i]+1 {
			break
		}
	}
	line = i + 1
	if i == 0 {
		column = int(e.pos)
	} else {
		column = int(e.pos) - (f.lines[i-1] + 1)
	}
	if len(f.name) > 0 {
		fmt.Println(f.name, "- Line:", line, "Column:", column, "-", e.msg)
	} else {
		fmt.Println("Line:", line, "Column:", column, "-", e.msg)
	}
}

func (f *File) PrintErrors() {
	for _, err := range f.errs {
		f.PrintError(err)
	}
}

func (f *File) Size() int {
	return f.size
}

func (f *File) ValidPos(p Pos) bool {
	return p >= f.base && p < f.base+Pos(f.size)
}

type Pos int

const NoPos Pos = 0

func (p Pos) IsValid() bool {
	return p > NoPos
}
