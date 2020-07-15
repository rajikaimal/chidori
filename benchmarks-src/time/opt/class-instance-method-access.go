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
	class Customer
		def initialize()
			@name = "Rajika"
			@city = "NYC"
		end
		def foo()
			puts(@name)
			puts(@city)
		end
	end
	i = 0
	a = 100
	while i < a do
		i = i + 1
		cust = Customer.new()
		cust.foo()
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
