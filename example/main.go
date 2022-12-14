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
	_, err := i.EvalInt("fun(asd)")
	if err != nil {
		fmt.Println(err)
	}
	a, _ := i.EvalInt("fun(5)")
	fmt.Println(a)
}