![Powered by](https://img.shields.io/badge/Powered%20By-Black%20Magic-orange.svg?longCache=true&style=flat-square) 
![Open Source Love](https://img.shields.io/badge/Open%20source-%E2%9D%A4%EF%B8%8F-brightgreen.svg?style=flat-square) 
![Gluten free](https://img.shields.io/badge/Gluten-Free-blue.svg?longCache=true&style=flat-square)


## Examples

### Variable assignment
```clojure
(def a "Hello world")
```
### Function assignment
```clojure
(def square (fn x 
  (* x x)))
```
### Anonymous function
```clojure
((fn x (* x x)) 4)
```
### Use function
```clojure
(out (square 4))
```
### Structs

```clojure
(def Animal (struct
  :name
  :age
  :race
))

(def tom (Animal
  "Tom"
  1
  :cat
))

(out (:name tom) " is a " (` (:race tom))) ; prints Tom is a cat
```

### Maps

```clojure
(def prices ({}
  (:tea :1.5$)
  (:coffee :2$)
  (:cake :3$)
))

(out "Cake costs " (` (:cake prices)))
(out "Tea costs " (` ([] :tea prices)))
```

### Macros

```clojure
(def dotimes (macro times |action| (
; || tells Hummus to not evaluate this argument but to
; literally take the AstNode as it's input parameter
  (map (range times) (fn
    (unquote action)
  ))
))

(dotimes 3 (out "Hello world"))
; Same as writing
; (out "Hello world")
; (out "Hello world")
; (out "Hello world")

(def when (macro cond action
  (quote 
    (if (unquote cond)
      (unquote action))
  )
))

(when (> 4 3) (out "A"))
; same as writing
; (if (> 4 3)
;   (out "A"))
```

### Function composition

```clojure
(def add (fn a b
  (+ a b)
))

(def square(fn x
  (* x x)
))

(def add-square-out (|> add square out))

(add-square-out 3 1) ; prints 16
```

### Examples

```clojure
(def fib (fn n
  (if (< n 2)
    1 
    (+ (fib (- n 2)) (fib (- n 1)))
  )
))

(for (range x (list "Hello" "World"))
  (out x)
)

(for true
  (out "Hello world")
)
```

## Exit the application
```clojure
(exit)
```

## Get help
```clojure
(help)
```
