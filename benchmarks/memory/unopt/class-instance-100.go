package main

import "github.com/goruby/goruby/object"
import "github.com/goruby/goruby/chidorilib"
import "github.com/pkg/profile"
import "runtime"

func main() {
	runtime.SetCPUProfileRate(10000)
	defer profile.Start(profile.MemProfile).Stop()
	
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
		v11.SetInstanceVariableDy(v11.Class, "@street", "5th street")
		iov12 := chidorilib.IO{Puts: v11.GetInstanceVariables()}
		iov12.Out()

		v13 := chidorilib.Method{
			"dynamicMethod",
			"main",
			nil,
		}

		v13.Body = func(o chidorilib.Object) {

			chidorilib.IO{Puts: o.GetInstanceVariableByName("@street")}.Out()

		}

		methodHashMap["dynamicMethod"] = v13

		v14 := programValuesList.Object["cust"]
		v14.Invoke("dynamicMethod")

	}
}
