package main

import "github.com/goruby/goruby/object"
import "github.com/goruby/goruby/chidorilib"
import "time"
import "os"

func main() {
	env := object.NewMainEnvironment()
	_, _ = env.Get("")

	now := time.Now()
	defer func() {
		timeNow := time.Since(now)
		f, err := os.OpenFile("dynamic-method-100.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
		}
		defer f.Close()
		if _, err := f.WriteString(timeNow.String() + "\n"); err != nil {
		}
	}()
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

		v9 := v8.Call("cust")
		programValuesList.Object["cust"] = v9

		v10 := chidorilib.Method{
			"dynamicMethod",
			"main",
			nil,
		}

		v10.Body = func(o chidorilib.Object) {

			chidorilib.IO{Puts: o.GetInstanceVariableByName("@name")}.Out()

		}

		methodHashMap["dynamicMethod"] = v10

		v11 := programValuesList.Object["cust"]
		v11.Invoke("dynamicMethod")

	}
}
