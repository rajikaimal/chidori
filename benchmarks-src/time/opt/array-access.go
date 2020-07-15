package main

import (
	"fmt"
	"github.com/goruby/goruby/compiler"
	"github.com/goruby/goruby/object"
)

func main() {
	//input := `
	//names = ["FOO", "BAR"]
	//`
	input := `
	myArr = ["Starlink", "Falcon9"]

	i = 0
	a = 100

	while i < a do
		i = i + 1

		myArr[1]
	end
		`
	c := compiler.New()

	out, err := c.Compile("", input)

	if err != nil {
		fmt.Println("Compile errors: ", err)
	}

	fmt.Println(out)
	res, _ := out.(object.RubyClassObject)
	fmt.Println(res)
}
