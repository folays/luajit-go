package luajit

import (
	"log"
	"reflect"
)

type funcEverywhere struct {
	v         any
	index_max int
}

func (L *State) function_prepare_everywhere() {
	L.function_everywhere = make(map[reflect.Type]*funcEverywhere)

	L.FuncEverywherePass(L, 0)
}

// FuncEverywherePass records TypeOf(v) and will pass it around in
// every Go functions needing one.
func (L *State) FuncEverywherePass(v any, type_in_index_max int) {
	var (
		typeOf = reflect.TypeOf(v)
	)

	switch {
	case typeOf.Kind() == reflect.Pointer && typeOf.Elem().Kind() == reflect.Struct:
	default:
		log.Fatalf("(*State).FuncEverywherePass(%T) : only `*struct` accepted")
	}

	L.function_everywhere[typeOf] = &funcEverywhere{
		v:         v,
		index_max: type_in_index_max,
	}
}

func (L *State) funcEverywhereGet(typeOf reflect.Type, type_in_index int, flagMethod bool) (v any) {
	if flagMethod {
		type_in_index -= 1
	}

	if everywhere := L.function_everywhere[typeOf]; everywhere != nil && type_in_index <= everywhere.index_max {
		return everywhere.v
	}

	return nil
}
