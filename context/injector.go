package context

import (
	"reflect"

	"github.com/sirupsen/logrus"
)

type Injector struct {
	log         *logrus.Logger
	injectables map[reflect.Type]any
}

func createInjector(log *logrus.Logger) *Injector {
	return &Injector{
		log:         log,
		injectables: make(map[reflect.Type]any),
	}
}

func Factory[T any](n *Injector) *T {
	return n.getInjectable(typeOf[T]()).(*T)
}

func (n *Injector) AddInjectable(i any) {
	t := reflect.TypeOf(i)
	if _, ok := n.injectables[t]; !ok {
		n.log.WithField("type", t.String()).Debug("Added Injectable")
		n.injectables[reflect.TypeOf(i)] = i
	} else {
		n.log.WithField("type", t.String()).Fatal("An injectable with this type already exists")
	}
}

func (n *Injector) getInjectable(t reflect.Type) any {
	_, ok := n.injectables[t]

	if !ok {
		n.construct(t)
	}

	return n.injectables[t]
}

func (n *Injector) construct(t reflect.Type) {
	c := reflect.New(t).Interface() // the requested type constructed as an interface

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if _, ok := f.Tag.Lookup("Inject"); ok {
			var inj reflect.Value
			if f.Type.Kind() == reflect.Pointer {
				inj = reflect.ValueOf(n.getInjectable(f.Type.Elem()))
			} else {
				inj = reflect.ValueOf(n.getInjectable(f.Type)) // recursively get the struct to inject based on the type
			}
			reflect.ValueOf(c).Elem().FieldByName(f.Name).Set(inj)
		}
	}

	n.injectables[t] = c
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
