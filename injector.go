package main

import (
	"reflect"
)

type Injector struct {
	injectables map[reflect.Type]interface{}
}

func CreateInjector() *Injector {
	return &Injector{
		injectables: make(map[reflect.Type]interface{}),
	}
}

func InjectorFactory[T any](n *Injector) *T {
	return n.GetInjectable(typeOf[T]()).(*T)
}

func (n *Injector) GetInjectable(t reflect.Type) interface{} {
	_, ok := n.injectables[t]

	if !ok {
		n.Construct(t)
	}

	return n.injectables[t]
}

func (n *Injector) Construct(t reflect.Type) {
	c := reflect.New(t).Interface() // the requested type constructed as an interface

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if _, ok := f.Tag.Lookup("Inject"); ok {
			var inj reflect.Value
			if f.Type.Kind() == reflect.Pointer {
				inj = reflect.ValueOf(n.GetInjectable(f.Type.Elem()))
			} else {
				inj = reflect.ValueOf(n.GetInjectable(f.Type)) // recursively get the struct to inject based on the type
			}
			reflect.ValueOf(c).Elem().FieldByName(f.Name).Set(inj)
		}
	}

	n.injectables[t] = c
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
