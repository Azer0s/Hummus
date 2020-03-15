![Powered by](https://img.shields.io/badge/Powered%20By-Black%20Magic-orange.svg?longCache=true&style=flat-square) 
![Open Source Love](https://img.shields.io/badge/Open%20source-%E2%9D%A4%EF%B8%8F-brightgreen.svg?style=flat-square) 
![Gluten free](https://img.shields.io/badge/Gluten-Free-blue.svg?longCache=true&style=flat-square)



## Commands

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

### Macros

```clojure
(def dotimes (macro times action (
  (quote
    (for (range times)
      (unquote action))
  )
)))

(dotimes 5 (out "Hello world")) ; prints Hello world 5 times

(def when (macro cond action (
  (quote 
    (if (unquote cond)
      (unquote action))
  )
)))

(when (> 4 3) (out "A")) ; prints A
(when (> 3 4) (out "A")) ; does not print
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
