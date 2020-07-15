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

	resultv1 := make([]object.RubyObject, 2)

	resultv1[0] = &object.String{Value: "Starlink"}

	resultv1[1] = &object.String{Value: "Falcon9"}

	v2 := &object.Array{Elements: resultv1}
	chidorilib.IO{Puts: v2.Inspect()}.Out()
	env.Set("myArr", v2)

	v3 := object.NewInteger(0)
	env.Set("i", v3)

	v4 := object.NewInteger(100)
	env.Set("a", v4)
	i, _ := env.Get("i")

	a, _ := env.Get("a")

	aVal, _ := a.(*object.Integer)

	iVal, _ := i.(*object.Integer)

	for iVal.Value < aVal.Value {

		i, _ := env.Get("i")

		v5 := object.NewInteger(1)
		iVal_ := i.(*object.Integer)
		iVal = iVal_
		v6 := iVal_.Value + v5.Value
		env.Set("i", &object.Integer{Value: v6})

		myArr := &object.Array{Elements: resultv1}
		chidorilib.IO{Puts: myArr.Elements[1].Inspect()}.Out()
	}
}
