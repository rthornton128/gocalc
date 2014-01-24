Contents
========

1 Introduction
2 Installation
3 License

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

Operators include: + (add), - (subtract), * (multiply), / (divide) and
% (remainder). An example

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

*There currently exists a bug with user methods where this is not true. At
present, user methods will discard extra arguments, or fail (sometimes
silently), if there are not enough. This will be fixed in the future.*


4 - License
===========
This project falls under a Simplified BSD License. You can find a copy of this
license in the file LICENSE in the root directory of this project. If one is
not included a copy may also be found at:

http://opensource.org/licenses/BSD-2-Clause
