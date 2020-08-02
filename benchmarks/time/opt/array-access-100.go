package main

import "github.com/goruby/goruby/object"
import "github.com/goruby/goruby/chidorilib"
import "time"
import "os"
import "runtime"
import "fmt"

func monitorRuntime() {
	m := &runtime.MemStats{}
	f, err := os.Create(fmt.Sprintf("mmem.csv"))
	if err != nil {
		panic(err)
	}
	f.WriteString("Time;Allocated;Total Allocated; System Memory;Num Gc;Heap Allocated;Heap System;Heap Objects;Heap Released;\n")
	for {
		runtime.ReadMemStats(m)
		f.WriteString(fmt.Sprintf("%s;%d;%d;%d;%d;%d;%d;%d;%d;\n", 1, m.Alloc, m.TotalAlloc, m.Sys, m.NumGC, m.HeapAlloc, m.HeapSys, m.HeapObjects, m.HeapReleased))
		time.Sleep(5 * time.Second)
	}
}

func main() {
	monitorRuntime()
	env := object.NewMainEnvironment()
	_, _ = env.Get("")

	now := time.Now()
	defer func() {
		timeNow := time.Since(now)
		f, err := os.OpenFile("array-access-100.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
		}
		defer f.Close()
		if _, err := f.WriteString(timeNow.String() + "\n"); err != nil {
		}
	}()
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
