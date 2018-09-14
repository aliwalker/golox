# Golox
This repo is a go implementation of Lox language. Lox is a dynamic scripting language designed by [Bob Nystrom](https://github.com/munificent) in his book [Crafting Interpreter](http://www.craftinginterpreters.com). Lox is really like JavaScript. Refer to [chapter 3](http://www.craftinginterpreters.com/the-lox-language.html) of the book for an overview of Lox lang.

This implementation adds some more features:

- [x] Semicolon is not a must. :-)
- [x] Lambda expressions(anonymous functions).
- [x] Support break statement from loops.
- [ ] Support getters/setters, static methods for classes.
- [ ] Support `for ... of` statement.
- [ ] Support Arrays, Maps.
- [ ] Enhanced REPL.

## Example

```kotlin
fun fibnacii() {
    var prev = 0, cur = 0

    // lambda expression.
    return () -> {
        if (cur == 0) {
            cur = 1
            return 0
        }

        var res = cur
        cur += prev
        prev = res

        return res
    }
}

var next = fibnacii()
for (var i = 0; i < 10; i += 1) {
    print next()
}
```