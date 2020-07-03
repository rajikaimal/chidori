package chidorilib

import "fmt"

type Value struct {
	Class  map[string]Class
	Object map[string]Object
	Method map[string]Method
}

type Cls interface {
	GetClassVariables() map[string]string
	SetInstanceVariables(varName string)
	Call() Object
}

type Class struct {
	ClassName      string
	ClassVariables map[string]string
	InstanceVar    map[string]int
	InstanceVarArr []string
	Methods        map[string]Method
	MethodsMap     map[string]int
	MethodsArr     []Method
}

func (c *Class) AddMethod(m Method) bool {
	if c != nil {
		methodExists := false
		for k, _ := range c.Methods {
			if k == m.Name {
				methodExists = true
			}
		}

		if !methodExists {
			c.Methods[m.Name] = m
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (c *Class) GetClassVariables() map[string]string {
	return c.ClassVariables
}

func (c *Class) SetInstanceVariables(variables map[string]string) {
	c.InstanceVar = make(map[string]int)

	//noOfVariables := len(variables)
	//variableSlice := make([]string, noOfVariables)
	var variableSlice []string

	arrIndex := 0

	for k, v := range variables {
		variableSlice = append(variableSlice, v)
		c.InstanceVar[k] = arrIndex
		arrIndex++
	}

	c.InstanceVarArr = variableSlice
}

func (c *Class) Call(objectName string) Object {
	var methodSlice []Method
	arrIndex := 0
	c.MethodsMap = make(map[string]int)

	for k, v := range c.Methods {
		methodSlice = append(methodSlice, v)
		c.MethodsMap[k] = arrIndex
		arrIndex++
	}

	c.MethodsArr = methodSlice

	obj := Object{
		Class: Class{
			ClassName:      c.ClassName,
			ClassVariables: c.ClassVariables,
			InstanceVar:    c.InstanceVar,
			InstanceVarArr: c.InstanceVarArr,
			Methods:        c.Methods,
			MethodsMap:     c.MethodsMap,
			MethodsArr:     c.MethodsArr,
		},
		ObjectName:        objectName,
		InstanceVariables: make(map[string]string),
	}

	init := obj.Methods["initialize"]
	if init.Name != "" {
		init.Body(&obj)
	}

	if obj.InstanceVariables == nil {
		obj.InstanceVariables = make(map[string]string)
	}

	return obj
}

type Obj interface {
	GetInstanceVariables() map[string]string
	GetInstanceVariableByName(name string)
	GetMethods() map[string]Method
	//SetInstanceVariables(name string, value string)
	SetInstanceVariableDy(name string, value string)
	Invoke(methodName string)
	AddMethod(m Method) bool
}

type Object struct {
	Class
	ClassName         string
	ObjectName        string
	ClassVariables    map[string]string
	InstanceVariables map[string]string
	DynamicMethods    map[string]Method
}

func (o *Object) GetInstanceVariables() map[string]string {
	instanceVariables := make(map[string]string)
	for k, v := range o.Class.InstanceVar {
		instanceVariables[k] = o.Class.InstanceVarArr[v]
	}
	return instanceVariables
}

func (o *Object) GetInstanceVariableByName(name string) string {
	variableSliceIdx := o.Class.InstanceVar[name]

	return o.Class.InstanceVarArr[variableSliceIdx]
}

func (o *Object) GetMethods() map[string]Method {
	return o.Methods
}

//func (o *Object) SetInstanceVariables(name string, value string) {
//	o.InstanceVariables[name] = value
//}
//
func (o *Object) SetInstanceVariableDy(class Class, name string, value string) {
	//not supported in optimized version
	return
}

func (o *Object) Invoke(methodName string) {
	methodIdx := o.MethodsMap[methodName]

	method := o.MethodsArr[methodIdx]
	if method.Name != "" {
		method.Body(o)
	}
}

func (o *Object) AddMethod(m Method) bool {
	if o != nil {
		methodExists := false
		for k, _ := range o.Methods {
			if k == m.Name {
				methodExists = true
			}
		}

		if !methodExists {
			o.DynamicMethods[m.Name] = m
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

type Method struct {
	Name      string
	BelongsTo string
	Body      func(o *Object)
}

type IO struct {
	Puts interface{}
}

func (io IO) Out() {
	fmt.Println(io.Puts)
}
