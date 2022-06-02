package context

import (
	"reflect"

	"github.com/sirupsen/logrus"
)

type Injector struct {
	injectables map[reflect.Type]any
}

func createInjector() *Injector {
	return &Injector{
		injectables: make(map[reflect.Type]any),
	}
}

func Factory[T any](n *Injector) *T {
	return n.getInjectable(typeOf[T]()).(*T)
}

func (n *Injector) AddInjectable(i any) {
	t := reflect.TypeOf(i)
	if _, ok := n.injectables[t]; !ok {
		Log.WithField("type", t.String()).Debug("Added Injectable")
		n.injectables[reflect.TypeOf(i)] = i
	} else {
		Log.WithField("type", t.String()).Fatal("An injectable with this type already exists")
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
	Log.WithField("type", t.String()).Debug("Constructing new instance for injection")
	c := reflect.New(t).Interface() // the requested type constructed as an interface

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if _, ok := f.Tag.Lookup("Inject"); ok {
			Log.WithFields(logrus.Fields{
				"field": f.Name,
				"type":  f.Type.String(),
			}).Debug("Found injectable field, getting injectable instance...")
			inj := reflect.ValueOf(n.getInjectable(f.Type))

			// destination field in the created struct
			d := reflect.ValueOf(c).Elem().Field(i)
			if d.IsValid() && d.CanSet() {
				d.Set(inj)
			} else {
				Log.WithFields(logrus.Fields{
					"type":  t.Name(),
					"field": f.Name,
				}).Fatal("Could not inject into field within type, is it valid and exported?")
			}
		}
	}

	n.injectables[t] = c
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
