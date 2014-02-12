
1 - Introduction
================

Calc is a toy programming language. It's original design was to stand as a
simple language spec by which to build a compiler and/or interpreter around.
It is meant as a learning tool for compiler design and little to no thought
has been given to its suitability for any other purpose.

It uses Lisp/Scheme-like syntax. It is not an implementation of any language
despite any similarities it may bare. The only numerical type it currently
has are Number's (int32 or int64 depending on architecture). Floating point
numbers are not currently supported but may be implemented at a later date,
as may be other numerical representations outside of decimal. Calc also has
a String type to represent character strings. Strings are not currently
well supported but it is planned to have at least comparison and concatenation
implemented at some future date.

Calc is very simple and lacks many, many features found in most modern
languages. At present it consists of just an interpreter and is thereby
just an interpreted language. It is planned to eventually implement a
compiler back-end, too.

2 - Installation
================

If you want to tinker and use the language it is best to clone a copy of the
project from it's github page. You will need a copy of Go 1.x installed to
compile the project.

If you wish to simply use the interpreter, for some odd reason, you could
probably use the 'go get' tool.

3 - Details
===========

Calc uses LISP-like S-Expressions, or reverse polish notation. Each expression
is encapsulated by brackets '()'. An expression must start with, and have at
most one of, an operator or method. It may contain zero or more arguments.
Each element must be separated by an empty space.

Currently implemented:

  * Basic mathematical operations: + - * / %
  * Logical: and or
  * Comparison: = <> < <= > >=
  * Assignment: set
	* Branching: if switch-case
	* Methods: define
	* Basic IO: print

An example:

(+ 3 2)

Methods may either be a built-in method or a user defined method. First an
example with no arguments:

(print)

...which prints a blank line to standard out, and an example with multiple
arguments:

(print "Hello world" "!")

... which prints "Hello world!" to standard out.

Operators and the print method take an arbitrary number of arguments but
most other builtin methods and user defined methods take an exact number of
arguments. Supplying the incorrect number of arguments to this methods will
result in a parsing error.

For working examples, check out the scripts sub directory which, currently,
has a fibonacci and a factorial example. There is also a test script which
you can read through with more example code. Uncomment some sections to
produce errors.


4 - More Information
====================

There is still a fair amount that needs to be implemented. As previously
mentioned, floating-point numbers and other number representations will
likely be added. Type assertions should be added. The ability to create
data structures is desired, too. Packages and importing are also planned.

There are things about Calc which the author does not like. One, it is not
too strictly typed. You can do bizarre things like have a function return
either a Number or a String which is a design flaw inherent to most (all?)
dynamically typed languages. The author prefers a type system which is both
strong and static.

Type assertions must be implemented as a result: Number? String?

Calc 2.0, therefore, may implement a stronger type system. A function
declaration may take the form of:

(define (func-name:int arg1 arg2:int arg3:string) (...))

This would define a function named func-name which returns an int, accepts
two arguments of type int and a third argument of type string.

Type inference can be used during instantiation of variables or it's type
could be clearly defined:

(set a 1)
(set a:int 1)


5 - License
===========
This project falls under a Simplified BSD License. You can find a copy of this
license in the file LICENSE in the root directory of this project. If one is
not included a copy may also be found at:

http://opensource.org/licenses/BSD-2-Clause
