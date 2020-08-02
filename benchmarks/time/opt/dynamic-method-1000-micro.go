package main

import "github.com/goruby/goruby/object"
import "github.com/goruby/goruby/chidorilib"
import "time"
import "os"

func main() {
	env := object.NewMainEnvironment()
	_, _ = env.Get("")

	methodHashMap := make(map[string]chidorilib.Method)

	v1 := chidorilib.Method{
		"initialize",
		"Customer",
		nil,
	}

	v1.Body = func(o *chidorilib.Object) {
		instanceVars := make(map[string]string)

		instanceVars["@name"] = "Rajika"

		o.SetInstanceVariables(instanceVars)
	}

	methodHashMap["initialize"] = v1

	v2 := chidorilib.Method{
		"foo",
		"Customer",
		nil,
	}

	v2.Body = func(o *chidorilib.Object) {

		chidorilib.IO{Puts: o.GetInstanceVariableByName("@name")}.Out()

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

	v6 := object.NewInteger(1000)
	env.Set("a", v6)
	i, _ := env.Get("i")

	a, _ := env.Get("a")

	aVal, _ := a.(*object.Integer)

	iVal, _ := i.(*object.Integer)
	now := time.Now()

	for iVal.Value < aVal.Value {

		i, _ := env.Get("i")

		v7 := object.NewInteger(1)
		iVal_ := i.(*object.Integer)
		iVal = iVal_
		v8 := iVal_.Value + v7.Value
		env.Set("i", &object.Integer{Value: v8})

		v9 := programValuesList.Class["Customer"]

		v10 := v9.Call("cust")
		programValuesList.Object["cust"] = v10

		v11 := programValuesList.Object["cust"]
		v11.Invoke("foo")
	}

	defer func() {
		timeNow := time.Since(now)
		f, err := os.OpenFile("dynamic-method-1000-micro.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
		}
		defer f.Close()
		if _, err := f.WriteString(timeNow.String() + "\n"); err != nil {
		}
	}()
}
