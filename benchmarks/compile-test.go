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
		o.SetInstanceVariables("@name", "RA")
	}

	methodHashMap["initialize"] = v1

	v2 := chidorilib.Method{
		"a",
		"Customer",
		nil,
	}

	v2.Body = func(o chidorilib.Object) {

		chidorilib.IO{Puts: o.GetInstanceVariableByName("@name")}.Out()

	}

	methodHashMap["a"] = v2

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

	v5 := programValuesList.Class["Customer"]

	v6 := v5.Call("cust")
	programValuesList.Object["cust"] = v6

	v7 := programValuesList.Object["cust"]
	v7.Invoke("a")
}
