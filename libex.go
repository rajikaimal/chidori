package main

import "github.com/goruby/goruby/chidorilib"
import "fmt"

func main() {
	programValuesList := chidorilib.Value{}

	methodHashMap := make(map[string]chidorilib.Method)
	v0 := chidorilib.Method{
		"initialize",
		"Object",
		nil,
	}

	v0.Body = func(o chidorilib.Object) {
		fmt.Println("BODY init")
		o.SetInstanceVariables("y", "yy")
		fmt.Println(o.GetClassVariables())
		fmt.Println(o.GetInstanceVariables())
		fmt.Println("End init")
	}

	methodHashMap["initialize"] = v0

	v1 := chidorilib.Method{
		"foo",
		"Object",
		nil,
	}

	v1.Body = func(o chidorilib.Object) {
		fmt.Println(o.GetClassVariables())
		fmt.Println("ENd")
	}

	methodHashMap["foo"] = v1

	v2 := "Init"
	v3 := make(map[string]string)
	v3["x"] = v2
	v4 := chidorilib.Class{
		ClassName:      "Example",
		ClassVariables: v3,
		Methods:        methodHashMap,
	}

	programValuesList.Class = make(map[string]chidorilib.Class)
	programValuesList.Class["Example"] = v4

	v5 := programValuesList.Class["Example"]

	//fmt.Println(exampleClass.GetClassVariables())

	v6 := v5.Call()

	fmt.Println(v5.GetClassVariables())
	v6.SetInstanceVariables("name", "RAJIKA")
	fmt.Println(v6.GetInstanceVariables())

	io := chidorilib.IO{Puts: "HELLO"}
	fmt.Println("HERE")
	io.Out()
	//v7 := v6.Class.Methods["foo"]
	//v7.Body(v6)
}
