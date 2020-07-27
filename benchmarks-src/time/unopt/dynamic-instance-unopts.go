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
		end
	end
	i = 0
	a = 100
	while i < a do
		i = i + 1
		class Customer
			attr_accessor :street
		end
		cust = Customer.new()
		cust.street = "dynamic var street"
		def cust.dynamicMethod
			puts(@street)
		end
		cust.dynamicMethod()
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
