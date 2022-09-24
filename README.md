# eval
支持函数和判断等功能的表达式计算器.

## Installation
```sh
go get -u github.com/NilCent/eval
```

## quick start
```go
package main

import (
	"github.com/NilCent/eval"
	"fmt"
)

func main() {
	i := eval.New()
	i.Do(`let fun = fn(x) {
		if (x > 10) {
			return x * 3;
		} else {
			return x * 5;
		}
	}`)
	a, _ := i.EvalInt("fun(5)")
	fmt.Println(a)
}
```