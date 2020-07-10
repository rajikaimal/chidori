package chidorilib

import "fmt"

type Value struct {
	Class  map[string]Class
	Object map[string]Object
	Method map[string]Method
}

type Cls interface {
	GetClassVariables() map[string]string
	SetInstanceVariable(varName string)
	Call() Object
}

type Class struct {
	ClassName      string
	ClassVariables map[string]string
	InstanceVar    map[string]string
	Methods        map[string]Method
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

func (c *Class) SetInstanceVariable(variableName string) {
	instanceHashMap := c.InstanceVar
	if instanceHashMap == nil {
		c.InstanceVar = make(map[string]string)
	}

	c.InstanceVar[variableName] = ""
}

func (c *Class) Call(objectName string) Object {
	obj := Object{
		Class: Class{
			ClassName:      c.ClassName,
			ClassVariables: c.ClassVariables,
			InstanceVar:    c.InstanceVar,
			Methods:        c.Methods,
		},
		ObjectName:        objectName,
		InstanceVariables: make(map[string]string),
	}

	init := obj.Methods["initialize"]
	if init.Name != "" {
		init.Body(obj)
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
	SetInstanceVariables(name string, value string)
	SetInstanceVariableDy(name string, value string)
	Invoke()
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
	return o.InstanceVariables
}

func (o *Object) GetInstanceVariableByName(name string) string {
	if _, ok := o.InstanceVariables[name]; ok {
		return o.InstanceVariables[name]
	}

	return ""
}

func (o *Object) GetMethods() map[string]Method {
	return o.Methods
}

func (o *Object) SetInstanceVariables(name string, value string) {
	o.InstanceVariables[name] = value
}

func (o *Object) SetInstanceVariableDy(class Class, name string, value string) {
	o.InstanceVariables[name] = value
}

func (o *Object) Invoke(methodName string) {
	method := o.Methods[methodName]
	if method.Name != "" {
		method.Body(*o)
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
	Body      func(o Object)
}

type IO struct {
	Puts interface{}
}

func (io IO) Out() {
	fmt.Println(io.Puts)
}
