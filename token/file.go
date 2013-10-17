package token

import "fmt"

type Error struct {
	pos Pos
	msg string
}

type File struct {
	base  Pos     // Base Pos of file
	errs  []Error // List of errors which have occured in processing
	lines []Pos   // Location of each line ending ('\n')
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

func (f *File) AddLine(p Pos) {
	if f.ValidPos(p) {
		f.lines = append(f.lines, p)
	}
}

func (f *File) AddError(p Pos, args ...interface{}) {
	if f.ValidPos(p) {
		f.errs = append(f.errs, Error{p, fmt.Sprint(args...)})
	}
}

func (f *File) NumErrors() int {
	return len(f.errs)
}

func (f *File) PrintError(e Error) {
	line, column := 1, int(e.pos)
	for i, p := range f.lines {
		//fmt.Println(e.pos, "vs", p)
		if e.pos < p {
			break
		}
		line = i + 1
		column = int(p-e.pos) + 1
	}
	//fmt.Println("e.pos:", e.pos, "f.lines:", f.lines)
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
