funlang is a functional language made for fun.

There are three data types in funlang - numbers, arrays and fuctions.
Functions take one argument and return one value. Numbers size is
implementation defined, but it is always signed. Recommended size for
implementations is atleast int32.

Comments are written after a semicolon.

Function calls have two forms a normal function call and a chained function
call. A chained function call is always started by a normal function call.

Expression are evaluated in a chain. Arguments are passed from one expression
to another using the -> operator. Example:

10->fib->print

You can put one of more expressions into () to group them into one expression.
For example if you want to have multiple expressions in one function or use the
result of an array reference. The value of the last expression is used as the
value of the parent expression. Examples:

fun (first->expression second->expression)
([1 : 2 : 3]![rand : 3]->mod)->print


Functions can be defined by using the fun keyword or the λ symbol followed by
a function call. You can use @ to substitute for the argument.

You can name an expression using the def keyword. The identifier is substituted
at runtime. This means defined expressions can be recursive. Additionally, @
in a deffed expression substitutes to the argument passed to the identifier,
not the argument of the host.

def square @->mul
def pi 5

Arrays are a big part of funlang. They are defined by writing the starting
elements inside [], separated by colons or a new line. Arrays can be indexed
using !expr, where expr evaluates to a number. The array ref operator has
bigger precedence than the chain operator. Example:

[1 : 2 : 3]!1
([1 : 2 : 3]![@ : 2]->mul)->print

Example programs:

Recursive fibonacci:

def fib
	[
		[ @ : 2 ]->greater ; if argument is greater than 2
		fun [ [@ : 1]->sub->fib : [@ : 2]->sub->fib ]->plus ; continue the recursion
	  fun @ ; else return argument
	]->branch

Prime numbers:

def primes
	[
		[2 : @]->range ; numbers from 2 to @
		0 ; closure isn't needed
		fun [ ; returns true if the number is prime
			[2 : @!x->sqrt]->range
			@!x
			fun [@!y : @!x]->mod
		]->filter->len->not
	]->filter

Function list:

Function are written as followed: (arg|[ arg ... ])->name

Flow control:
[ value true false ]->branch       executes the true function if value isn't 0,
                                   else executes the false function

Functional tools:
[ function argument ]->call        calls the function with the argument
[ array closure [element closure]->function->element ]->map
                                   for each element, passes array containing
                                   said element and closure to the function
                                   and constructs a new array from the return
                                   values.

[ array closure [element closure value]->function->value ]->reduce
                                   for each element, passes array containing
                                   said element, the closure and the result of
                                   the last call (or 0) to the function. then
                                   returns the result of the last function

[ array closure [element closure]->function->bool ]->filter
                                   for each element passes an array containing
                                   said element and the closure to the
                                   function. if the return value is 0, the
                                   element is removed.

Array operations:
[ array index ]->delete            deletes element at index from the array
array->len                         returns the length of the array
[ min max ]->range                 returns an array starting at min (inclusive) and 
                                   ending at max (exclusive)
[ array ... ]->append              appends another arrays to the end of an array
[ array ... ]->add                 adds arguments to the end of an array 
array->car                         returns the first element of an argray
array->cdr                         returns the array without the first element

Arithmetics:
[ number ... ]->add                adds numbers together
[ number ... ]->sub                subtracts numbers from each otheri
[ number ... ]->mul                multiplies number amongts each other
[ number ... ]->div                divides numbers by each other

Logic:
[ number ... ]->greater            returns one if each number is greater than
                                   the following one
[ number ... ]->lesser             returns one if each number is lesser than
                                   the following one
[ number ... ]->equal              returns one if all numbers are equal
[ number ... ]->or                 returns one if atleast one number is not 0
[ number ... ]->and                returns one if all numbers are not 0

IO:
[ ... ]->print                     prints a string representation of all the
                                   arguments to the stdout
[ 0 ]->read                        reads one line from standard input into
                                   an array of utf32 characters
[ utf32 ... ]->prints              prints an array of utf32 as a string

Misc:
[ string ]->parse                  parser a utf32 array and returns a native
                                   data structure
