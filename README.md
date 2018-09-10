# Golox
This repo is a go implementation of Lox language. Lox is a dynamic scripting language designed by [Bob Nystrom](https://github.com/munificent) in his book [Crafting Interpreter](http://www.craftinginterpreters.com). Refer to [chapter 3](http://www.craftinginterpreters.com/the-lox-language.html) of the book for an overview of Lox lang.

This implementation adds some more features:

- Semicolon is not a must(you can add it if you wanted too).
- Lambda expressions(anonymous functions).
- Enhanced REPL.

## Example

```kotlin
fun fibnacii() {
    var prev = 0, cur = 0

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