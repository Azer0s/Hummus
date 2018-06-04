![Powered by](https://img.shields.io/badge/Powered%20By-Black%20Magic%20and%20lambda%20calculus-orange.svg?longCache=true&style=flat-square)<br>
![Open Source Love](https://img.shields.io/badge/Open%20source-%E2%9D%A4%EF%B8%8F-brightgreen.svg?style=flat-square)<br>
![Gluten free](https://img.shields.io/badge/Gluten-Free-blue.svg?longCache=true&style=flat-square)
## Commands

### Variable assignment
```r
name:=value
```
### Function assignment
```r
name:=(arguments - comma separated).(process)
```
### Anonymous function
```r
(arguments - comma separated).(process).(values - comma separated)
```
### Use function
```r
name(arguments - comma seperated)
```
### Examples
```r
x:=2
y:=(x).(x*x)
z:=y(y(x))

a:=(x,y).(x)
b:=a(true,false)

k:=(i).(i ? y(10):y(12))
(x).(x ? Hello : World).(true)
```

## Open file
> LBox.exe test.lbox

## Exit the application
> exit

## Clear the console
> console

## Get help
> help

## Enable recursion
> rec
