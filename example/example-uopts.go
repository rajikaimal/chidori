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
		def foo()
			@address = "NYC"
			puts(@address)
		end
	end

	i = 0
	a = 4

	while i < a do
		i = i + 1
		class Customer
			attr_accessor :i	
		end

		cust = Customer.new()
		cust.i = "dynamic var i"
	
		class Customer
			attr_accessor :a	
		end

		cust = Customer.new()
		cust.a = "dynamic var i"


		def cust.dynamicMethod
			puts(@a)
		end

		cust.dynamicMethod()
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
