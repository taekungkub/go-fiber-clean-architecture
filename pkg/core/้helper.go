package core

import "reflect"

func ApplyPatch(dst interface{}, src interface{}) {
	dv := reflect.ValueOf(dst).Elem()
	sv := reflect.ValueOf(src).Elem()

	for i := 0; i < sv.NumField(); i++ {
		sf := sv.Field(i)
		if sf.IsNil() {
			continue
		}
		dv.Field(i).Set(sf.Elem())
	}
}
