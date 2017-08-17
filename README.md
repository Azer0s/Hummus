![Powered by](https://img.shields.io/badge/Powered%20by-Black%20magic%20and%20lambda%20calculus-orange.svg)

[![Open Source Love](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)](https://github.com/ellerbrock/open-source-badges/)

![Gluten free](http://forthebadge.com/images/badges/gluten-free.svg)
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
