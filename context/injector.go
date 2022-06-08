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

func (n *Injector) AddInjectable(c any) {
	t := dereferencePointer(reflect.TypeOf(c))

	if _, ok := n.injectables[t]; !ok {
		Log.WithField("type", t.String()).Debug("Added Injectable")
		n.injectables[t] = n.inject(c, t)
	} else {
		Log.WithField("type", t.String()).Fatal("An injectable with this type already exists")
	}
}

func (n *Injector) getInjectable(t reflect.Type) any {
	if _, ok := n.injectables[t]; !ok {
		n.construct(t)
	}

	return n.injectables[t]
}

func (n *Injector) construct(t reflect.Type) {
	Log.WithField("type", t.String()).Debug("Constructing new instance for injection")
	c := reflect.New(t).Interface() // the requested type constructed as an interface
	n.injectables[t] = n.inject(c, t)
}

func (n *Injector) inject(c any, t reflect.Type) any {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if _, ok := f.Tag.Lookup("Inject"); ok {
			Log.WithFields(logrus.Fields{
				"field": f.Name,
				"type":  f.Type.String(),
			}).Debug("Found injectable field, injecting...")
			inj := reflect.ValueOf(n.getInjectable(dereferencePointer(f.Type)))

			// destination field in the created struct
			d := reflect.ValueOf(c).Elem().Field(i)
			if d.IsValid() && d.CanSet() {
				d.Set(inj)
			} else {
				Log.WithFields(logrus.Fields{
					"type":  t.Name(),
					"field": f.Name,
				}).Fatal("Could not inject into field within this type, is it valid and exported?")
			}
		}
	}
	return c
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func dereferencePointer(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Pointer {
		return t.Elem()
	} else {
		return t
	}
}
