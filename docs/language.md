# Hummus Language Documentation

## Overview
Hummus is a functional programming language with Lisp-like syntax, supporting macros, first-class functions, structs, and native plugins.

## Syntax
- **Lists**: `(a b c)`
- **Atoms**: `:keyword`
- **Strings**: `"text"`
- **Numbers**: `123`, `4.56`
- **Booleans**: `true`, `false`

## Core Forms
- **Definition**: `(def name value)`
- **Function**: `(fn [args] body)`
- **Condition**: `(if test then else)`
- **Loop**: `(for var seq body)`
- **Struct**: `(struct Name fields)`
- **Macro**: `(macro Name args body)`

## Types
- Primitive: Int, Float, String, Bool, Atom, Nothing
- Collection: List, Map
- Custom: Function, Struct, Macro, Option

## Standard Library
Modules are imported with `(use :<module>)`. Available modules include `base`, `str`, `file`, `json`, `net/http`, `log`, and more.

## Building & Running
Compile the project and run scripts using the `hummus` CLI tool. Use `hummus run <file>` to execute a file or `hummus` to start the REPL.
