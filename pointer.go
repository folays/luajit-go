package luajit

import (
	"log"
	"reflect"
)

func (L *State) pushGoPointer(v reflect.Value) {
	if v.Kind() != reflect.Pointer {
		log.Fatalf("(*State).pushGoPointer() error : .Kind() is not Pointer")
	}

	var (
		vFinal reflect.Value
	)

	for vFinal = v; vFinal.Kind() == reflect.Pointer; vFinal = vFinal.Elem() {
	}

	switch vFinal.Kind() {
	case reflect.Struct:
		L.pushGoStruct(v)
	default:
		log.Fatalf("(*State).pushGoPointer() error : [.Elem...].Kind() is not a Struct")
	}
}
