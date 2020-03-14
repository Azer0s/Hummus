![Powered by](https://img.shields.io/badge/Powered%20By-Black%20Magic-orange.svg?longCache=true&style=flat-square) 
![Open Source Love](https://img.shields.io/badge/Open%20source-%E2%9D%A4%EF%B8%8F-brightgreen.svg?style=flat-square) 
![Gluten free](https://img.shields.io/badge/Gluten-Free-blue.svg?longCache=true&style=flat-square)

Hummus is a LISPish language written in Go

## Commands

### Variable assignment
```lisp
(def a "Hello world")
```
### Function assignment
```lisp
(def square (fn x 
	(* x x)))
```
### Anonymous function
```lisp
((fn x (* x x)) 4)
```
### Use function
```lisp
(out (square 4))
```
### Examples
```lisp
(def fib (fn n
	(if (< n 2)
		1 
    (+ (fib (- n 2)) (fib (- n 1))))
))
```

## Exit the application
```lisp
(exit)
```

## Get help
```lisp
(help)
```