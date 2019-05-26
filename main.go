package main

import (
	"fmt"

	"github.com/gopherjs/gopherwasm/js"
)

func sayHi(args []js.Value) {
	fmt.Printf("%s\n", args[0].String())
	js.Global().Call("facetEngineCallback", "one,two,three")
}

func registerCallbacks() {
	js.Global().Set("sayHi", js.NewCallback(sayHi))
}

func main() {
	c := make(chan struct{}, 0)
	fmt.Println("hello")
	registerCallbacks()
	<-c
}
