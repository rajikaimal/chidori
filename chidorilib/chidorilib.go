package chidorilib

import "fmt"

type ProgramValues struct {
	Classes map[string]Class
}

type Value struct {
	Class  map[string]Class
	Object map[string]Object
	Method map[string]Method
}

type Cls interface {
	GetClassVariables() map[string]string
	Call() Object
}

type Class struct {
	ClassName      string
	ClassVariables map[string]string
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

func (c *Class) Call() Object {
	obj := Object{
		Class: Class{
			ClassName:      c.ClassName,
			ClassVariables: c.ClassVariables,
			Methods:        c.Methods,
		},
		InstanceVariables: make(map[string]string),
	}

	init := obj.Methods["initialize"]
	init.Body(obj)

	return obj
}

type Obj interface {
	GetInstanceVariables() map[string]string
	SetInstanceVariables(name string, value string)
}

type Object struct {
	Class
	ClassName         string
	ClassVariables    map[string]string
	InstanceVariables map[string]string
}

func (o *Object) GetInstanceVariables() map[string]string {
	return o.InstanceVariables
}

func (o *Object) SetInstanceVariables(name string, value string) {
	o.InstanceVariables[name] = value
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
