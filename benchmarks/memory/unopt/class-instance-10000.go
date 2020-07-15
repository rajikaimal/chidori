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

	v2 := chidorilib.Method{
		"foo",
		"Customer",
		nil,
	}

	v2.Body = func(o chidorilib.Object) {
		o.SetInstanceVariables("@address", "NYC")

		chidorilib.IO{Puts: o.GetInstanceVariableByName("@street")}.Out()

	}

	methodHashMap["foo"] = v2

	programValuesList := chidorilib.Value{}

	v3 := make(map[string]string)
	v4 := chidorilib.Class{
		ClassName:      "Customer",
		ClassVariables: v3,
		Methods:        methodHashMap,
	}

	programValuesList.Class = make(map[string]chidorilib.Class)
	programValuesList.Object = make(map[string]chidorilib.Object)
	programValuesList.Class["Customer"] = v4

	v5 := object.NewInteger(0)
	env.Set("i", v5)

	v6 := object.NewInteger(10000)
	env.Set("a", v6)
	i, _ := env.Get("i")

	a, _ := env.Get("a")

	aVal, _ := a.(*object.Integer)

	iVal, _ := i.(*object.Integer)

	for iVal.Value < aVal.Value {

		i, _ := env.Get("i")

		v7 := object.NewInteger(1)
		iVal_ := i.(*object.Integer)
		iVal = iVal_
		v8 := iVal_.Value + v7.Value
		env.Set("i", &object.Integer{Value: v8})

		v9 := programValuesList.Class["Customer"]
		v9.SetInstanceVariable("@street")

		v10 := programValuesList.Class["Customer"]

		v11 := v10.Call("cust")
		programValuesList.Object["cust"] = v11

		v12 := programValuesList.Object["cust"]
		v12.SetInstanceVariableDy(v12.Class, "@street", "dynamic var street")
		iov13 := chidorilib.IO{Puts: v12.GetInstanceVariables()}
		iov13.Out()

		v14 := chidorilib.Method{
			"dynamicMethod",
			"main",
			nil,
		}

		v14.Body = func(o chidorilib.Object) {

			chidorilib.IO{Puts: o.GetInstanceVariableByName("@a")}.Out()

		}

		methodHashMap["dynamicMethod"] = v14

		v15 := programValuesList.Object["cust"]
		v15.Invoke("dynamicMethod")

		v16 := programValuesList.Object["cust"]
		v16.Invoke("foo")

	}
}
