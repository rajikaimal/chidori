# chidori

This project is based on the wonderful work done by [GoRuby](https://github.com/goruby/goruby) project

This is the implementation to support the work presented in my [MSc Thesis](https://www.researchgate.net/publication/343685520_Chidori_-_Ruby_performance_analysis)

GoRuby, an implementation of Ruby written in Go. This repository contains compiler for subset of GoRuby.

```
$ go build chidori.go
$ ./chidori compiled <rubyfile.rb>
$ ./chidori run compiled.go
```

## Example

Ruby source code

```rb
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
	cust.street = "5th Street"
	def cust.getStreet
			puts(@street)
	end
  cust.getStreet()
end
```

Transpiled Go code

```go
package main

import "github.com/goruby/goruby/object"
import "github.com/goruby/goruby/chidorilib"

func main() {
	env := object.NewMainEnvironment()
	_, _ = env.Get("")

	methodHashMap := make(map[string]chidorilib.Method)

	v1 := chidorilib.Method{
		"initialize",
		"Customer",
		nil,
	}

	v1.Body = func(o chidorilib.Object) {
		o.SetInstanceVariables("@name", "Rajika")
	}

	methodHashMap["initialize"] = v1

	programValuesList := chidorilib.Value{}

	v2 := make(map[string]string)
	v3 := chidorilib.Class{
		ClassName:      "Customer",
		ClassVariables: v2,
		Methods:        methodHashMap,
	}

	programValuesList.Class = make(map[string]chidorilib.Class)
	programValuesList.Object = make(map[string]chidorilib.Object)
	programValuesList.Class["Customer"] = v3

	v4 := object.NewInteger(0)
	env.Set("i", v4)

	v5 := object.NewInteger(100)
	env.Set("a", v5)
	i, _ := env.Get("i")

	a, _ := env.Get("a")

	aVal, _ := a.(*object.Integer)

	iVal, _ := i.(*object.Integer)

	for iVal.Value < aVal.Value {

		i, _ := env.Get("i")

		v6 := object.NewInteger(1)
		iVal_ := i.(*object.Integer)
		iVal = iVal_
		v7 := iVal_.Value + v6.Value
		env.Set("i", &object.Integer{Value: v7})

		v8 := programValuesList.Class["Customer"]
		v8.SetInstanceVariable("@street")

		v9 := programValuesList.Class["Customer"]

		v10 := v9.Call("cust")
		programValuesList.Object["cust"] = v10

		v11 := programValuesList.Object["cust"]
		v11.SetInstanceVariableDy(v11.Class, "@street", "5th Street")
		iov12 := chidorilib.IO{Puts: v11.GetInstanceVariables()}
		iov12.Out()

		v13 := chidorilib.Method{
			"getStreet",
			"main",
			nil,
		}

		v13.Body = func(o chidorilib.Object) {

			chidorilib.IO{Puts: o.GetInstanceVariableByName("@street")}.Out()

		}

		methodHashMap["getStreet"] = v13

		v14 := programValuesList.Object["cust"]
		v14.Invoke("getStreet")

	}
}
```

MIT
