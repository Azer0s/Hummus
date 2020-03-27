![Powered by](https://img.shields.io/badge/Powered%20By-Black%20Magic-orange.svg?longCache=true&style=flat-square) 
![Open Source Love](https://img.shields.io/badge/Open%20source-%E2%9D%A4%EF%B8%8F-brightgreen.svg?style=flat-square) 
![Gluten free](https://img.shields.io/badge/Gluten-Free-blue.svg?longCache=true&style=flat-square)

## Getting started

`makefile` coming soon. Meanwhile, this will have to do. 🤷🏻‍♂️

```bash
git clone https://github.com/Azer0s/Hummus.git
cd Hummus
chmod +x ./scripts/prepare_release.sh
./scripts/prepare_release.sh
go build -o bin/hummus
./bin/hummus
```

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
  (map (range times) (list :fn
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
    (list :if (unquote cond)
      (unquote action))
  )
))

(when (> 4 3) (out "A"))
; same as writing
; (if (> 4 3)
;   (out "A"))
```


### Actor model

```clojure
(def ping (fn
  (for true
    (def msg (receive))
    (def type (nth 0 msg))
    (def sender (nth 1 msg))

    (if (= type :pong)
      ((fn
        (out "PONG")
        (sync/sleep 1 :s)
        (send sender (list :ping self))
      ))
      (out "Invalid message!")
    )
  )
))

(def pong (fn
  (for true
    (def msg (receive))
    (def type (nth 0 msg))
    (def sender (nth 1 msg))

    (if (= type :ping)
      ((fn
        (out "PING")
        (sync/sleep 1 :s)
        (send sender (list :pong self))
      ))
      (out "Invalid message!")
    )
  )
))

(def pingPid (spawn ping))
(def pongPid (spawn pong))

(send pongPid (list :ping pingPid))
(in)
```

### Map ⇄ Filter ⇄ Reduce

```clojure
(def pilots (list
  ({}
    (:id 2)
    (:name "Wedge Antilles")
    (:faction "Rebels")
  )
  ({}
    (:id 8)
    (:name "Ciena Ree")
    (:faction "Empire")
  )
  ({}
    (:id 40)
    (:name "Iden Versio")
    (:faction "Empire")
  )
  ({}
    (:id 66)
    (:name "Thane Kyrell")
    (:faction "Rebels")
  )
))

(each
  (map pilots (fn x
    (str/concat (` (:name x)) " => " (` (:faction x)))
  ))
(fn x
  (out x)
))

(each
  (filter pilots (fn x
    (= (:faction x) "Empire")
  ))
(fn x
  (out (:name x))
))

(out (reduce pilots (fn x acc
  (if (= (:faction x) "Empire")
    (+ acc 1)
    acc
  )
) 0))
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

(pipe/do pilots
  (map.. (fn x
    (str/concat (` (:name x)) " => " (` (:faction x)))
  ))
  (each.. (fn x
    (out x)
  ))
)

(pipe/do pilots
  (filter.. (fn x
    (= (:faction x) "Empire")
  ))
  (map.. (fn x
    (:name x)
  ))
  (each.. (fn x
    (out x)
  ))
)

(pipe/do pilots
  (reduce.. (fn x acc
    (if (= (:faction x) "Empire")
      (+ acc 1)
      acc
    )
  ) 0)
  out
)
```

### HTTP Server

```clojure
(http/handle "/" (fn req
  "<h1> Hello </h1>"
))

(http/handle "/test" (fn req
  "<h1> Test </h1>"
))

(http/serve ":8080")
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
