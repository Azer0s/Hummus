### Naming

Some functions from the std lib are prefixed, others aren't. A function is prefixed if
its functionality is not a core function of the language. If a function is not prefixed,
the function name should never be overwritten because the functionality of it is so
grave that it should never be masked.

### Commenting

Comments start with the function name followed by a short description of the function.
In the `in` list, the input variables are described. In the `out` section, the return
value is described.
<br>
E.g.:

```clojure
;; str/nth Returns a character of a string
;; in:
;; * arg .. the lookup string
;; * idx .. the index of the character
;; out: the character
```