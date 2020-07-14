package main

import "github.com/goruby/goruby/object"
import "github.com/goruby/goruby/chidorilib"
import "time"
import "os"
import "github.com/pkg/profile"

func main() {
	env := object.NewMainEnvironment()
	_, _ = env.Get("")

	var resultv1 []object.RubyObject

	now := time.Now()
	defer func() {
		timeNow := time.Since(now)
		f, err := os.OpenFile("time.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
		}
		defer f.Close()
		if _, err := f.WriteString(timeNow.String() + "\n"); err != nil {
		}
	}()

	resultv1 = append(resultv1, &object.String{Value: "Starlink"})

	resultv1 = append(resultv1, &object.String{Value: "Falcon9"})

	resultv1 = append(resultv1, &object.String{Value: "Falcon9"})

	resultv1 = append(resultv1, &object.String{Value: "Falcon9"})

	resultv1 = append(resultv1, &object.String{Value: "Falcon9"})

	resultv1 = append(resultv1, &object.String{Value: "Falcon9"})

	resultv1 = append(resultv1, &object.String{Value: "Falcon9"})

	v2 := &object.Array{Elements: resultv1}
	chidorilib.IO{Puts: v2.Inspect()}.Out()

	myArr := &object.Array{Elements: resultv1}
	chidorilib.IO{Puts: myArr.Elements[1].Inspect()}.Out()
}
