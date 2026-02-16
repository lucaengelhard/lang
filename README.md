# Syntax

## Variables

- Type can be inferred by the result of the assigned expression or explicit
- All variables are immutable by default
```rust
let a = 1; 
let b: int = 2;

let mut c = 3;
c += 1;
```
## Functions

- Return types are inferred or explicit
- Functions are values -> anonymous functions
```rust
fn foo(a: int, b: string) {
  ...
}

fn bar(...) -> ReturnType {
  ...
}

let baz = (...) => {
  ...
};

let barbaz = (...) -> ReturnType {
  ...
};
```
- Functions args definable as positional or named arguments
```rust
fn foo(a: int, b: string) {
  ...
}

foo(1, "abc");

foo(1, b: "abc");

foo(b: "abc", a: 1);

!! INVALID: f(b: "abc", 1);

```
- Instead of named args, a Record (key-value interface) can be passed
- Instead of positional args, an Indexable (any object indexable by a uint) can be passed
- Records and Indexables only "cover" the args they define, the rest can be bassed normally
```rust

fn foo(a: int, b: string, c: float) {
  ...
}

let dict = {
  a: 1,
  b: "abc",
  c: 2.2
};
foo(dict);

let arr = [1, "abc", 2.2];
foo(arr);

let partial_dict = {a: 2, c: 2.2};
foo(partial_dict, b: "def");

let partial_array = [1, "def"];
foo(partial_array, 2.2);

```

- Functions callable as members of type of first argument
```rust
fn foo(x: int) {
  ...
}

let a: int = 3;

a.foo();
```
- By default every arg is passed by value -> Reference/Pointer needs to be explicitly stated
```rust
fn foo(x: int) {
  ...
}

fn bar(y: *int) {
  ...
}

let a: int = 2;

foo(a);

bar(&a);
```

- Functions can be generic
- Generics can be restricted by an interface

```rust
fn foo<T>(x: T) {
  ...
}

fn bar<T satisfies Baz>(x: T) {
  ...
}
```

## Conditionals

- if statements don't evaluate to a value, but inline evaluates to result of expression (Statement if one of the branches is a block)
```rust
if (a < b) {
  ...
} else {
  ...
}

if (...) {
  ...
} else if {
  ...
}

let baz = if (...) ... else ...;
```
## Switch

- Switch statements can match on values/expressions, structs or interfaces.
- When matching on structs or interfaces without key properties, they can also be destructured (and the properties can also be matched by value)
```rust

interface Foo {
  ...
}

struct Bar {
  a: int;
  b: sring;
}


switch a {
  1     => {...},
  Foo   => {...},
  Bar   => {...},
  a < 2 => {...},
}

switch x {
  Bar{a: 1},            // Only matches if a == 1;
  Bar{a, b}   => {...}, // a and be can be used as variables
}
```
- default case defined by just capturing the variable (it's always considered true) 
```rust
switch x {
  ...
  a => {} 
}
```
- Switch statements (like if statements) return a value if there is no block?
```rust
switch x {
  ... => 
}
```

## Loops
- Iterating through an Iterable (Arrays, Lists, key-value-pairs of hashmaps) also possible

```rust
for (let mut i = 0; i < 10; i++) {
  ...
}

for (let mut i = 0; i < 10; i++) {
  ...
  continue
}

for (let mut i = 0; i < 10; i++) {
  ...
  break
}

let a: Iterable<int> = [1, 2, 3, 4];

for (el in a) {
  ...
}

for (el, index in a) {
  ...
}

while (...) {
  ...
}

```

## Types, Structs, Interfaces, Enums

- Interfaces are evaluated loosely -> They only need to be satisified not matched exactly (values can have more properties)
```rust
interface Addable {
  add();
}

fn foo(x: Addable) {
  ...
}

let a:int = 2;

foo(a);
```

- Structs only have properties, no method (though properties can be functions)
- All properties need to have a value (default values need to be initialized) 

```rust
struct Bar {
  a: int;
  b: sring;
  c: (x: int, y: int) -> int;  
}

let x = Bar{
  a: 1,
  b: "abc",
  c: (x, y) {x + y} // if anonymous function, types of arguments can be inferred
};

struct Baz {
  a: int = 0;
}

let y = Baz{}; // Baz{a: 0}
```
- Enum keys are zero indexed uints by default, but can be initialized as explicit values
- When A value is a number, the following values are that value + 1 if not otherwise defined 
```rust
enum Bar {                // Enum<int>
  VALUE,
  ANOTHERVALUE
}

enum Bar {
  VALUE = 2,
  ANOTHERVALUE            // Bar.ANOTHERVALUE evaluates to 3
}

enum Baz {                // Enum<string>
  VALUE = "abc",
  ANOTHERVALUE = "def",
}

enum Foo {                // Enum<string, int>
  VALUE = "abc",
  ANOTHERVALUE,           // Bar.ANOTHERVALUE evaluates to 1? (index in enum)
}
```

- Interfaces and Structs can be generic
```rust
interface Foo<T> {
  bar() -> T;
}

struct Baz<T> {
  a: T;
}

interface Bar<T satisfies Foo> {
  ...
}
```

- typeof and satisifies can also be used in normal code
- typeof can be used for stricted type equality
```rust
interface Foo {
  ...
}

if (bar satisfies Foo) {
  ...
}

let a = 1;

let t = typeof a    // t = int 
```