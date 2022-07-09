package luajit

import (
	"fmt"
	"log"
	"reflect"
)

/*
#include <stdlib.h>
#include <lua.h>

#include "function.h"
*/
import "C"

type LuaCFunction C.lua_CFunction

type goFunction struct {
	//o           any
	//method_n    int
	//method_name string
	v          reflect.Value
	flagMethod bool
}

func (L *State) function_prepare_metatable() {
	L.function_fac = L.factory_new()
}

func (L *State) function_unprepare() {
}

func _isGoFunction_a_method(v reflect.Value) bool {
	if v.Type().NumIn() == 0 {
		return false
	}

	for i := 0; i < v.Type().In(0).NumMethod(); i++ {
		var (
			method      = v.Type().In(0).Method(i)
			method_func = method.Func
		)
		if method_func.Pointer() == v.Pointer() {
			return true
		}
	}
	return false
}

func (L *State) pushGoFunction(v reflect.Value) {
	gf := &goFunction{
		v:          v,
		flagMethod: _isGoFunction_a_method(v),
	}

	L.pushUpvalueState()
	L.factory_produce(L.function_fac, gf)

	L.pushcclosure(LuaCFunction(C.function_call_c), 2)
}

//export function_call_go
func function_call_go(l LuaStatePtr) Index {
	var (
		L  = toUpvalueState(l, 1)
		gf = L.factory_checkFatal(L.function_fac, upvalueindex(2)).(*goFunction)
	)

	var (
		t         = gf.v.Type()
		inValues  []reflect.Value
		tin_index       = 0 // t.In() indexes start at 0
		lua_index Index = 1 // Lua    indexes start at 1
	)

	for tin_index = 0; tin_index < t.NumIn(); tin_index++ {
		var (
			expect_type = t.In(tin_index)
		)

		if v := L.funcEverywhereGet(expect_type, tin_index, gf.flagMethod); v != nil {
			inValues = append(inValues, reflect.ValueOf(v))
			continue // do not increase `lua_index`
		}

		{
			var (
				data  = L.ToAny(lua_index)
				value = reflect.ValueOf(data)

				value2        reflect.Value = value
				has_converted bool
				_             = has_converted
			)

			if value.IsValid() == false {
				value2 = reflect.Zero(expect_type)
			}

			for value2.Type() != expect_type {
				switch {
				case value2.CanConvert(expect_type):
					value2 = value2.Convert(expect_type)
					has_converted = true
				case value2.Kind() == reflect.Pointer:
					value2 = value2.Elem()
				default:
					err := fmt.Errorf("call [%d] : expect %v ; provide %T %v ; cannot convert",
						tin_index, expect_type, data, data)
					log.Printf("luajit.function_call_go() error : %s\n", err)
					L.SetTop(0)
					L.PushAny(err)
					return -1 // lua_error()
				}
			}

			//if silence == false {
			//	fmt.Printf("call [%d] : expect %v ; provide %T ; has_converted %t ; -> %s\n",
			//		tin_index, expect_type, data, has_converted, value2.Type())
			//}

			inValues = append(inValues, value2)
			lua_index++
		}
	}

	if len(inValues) != t.NumIn() {
		err := fmt.Errorf("\033[31mfunction_call_go\033[39m with %d args, inferred %d, expected %d", L.GetTop(), len(inValues), t.NumIn())
		log.Printf("luajit.function_call_go() error : %s\n", err)
		L.SetTop(0)
		L.PushAny(err)
		return -1 // lua_error()
	}

	{
		// We already checked that L.GetTop() == t.NumIn().
		//
		// Regarding the state of the stack :
		// - stack[0]    : does not exist (just a precision) : first index is at [1]
		// - stack[1]    : we               let it   on the stack
		// - stack[2..n] : we may choose to let them on the stack
		//
		// stack[1]          is         needed at least for XXX
		// stack[2..n] maybe be someday wanted          for some other usages.

		//L.SetTop(0)
		L.SetTop(1)

		// make (*State)._run() sanity checks happy.
		// those would be triggered if our callee do some further calls
		L._stackTops_push(1)
	}

	outValues := gf.v.Call(inValues)

	{
		L._stackTops_pop()

		// We let stack[1] on the stack before gf.V.Call(). Remove it.
		L.remove(1)
	}

	for _, out := range outValues {
		//fmt.Printf("Ret [?] : %T %v\n", out.Interface(), out)
		L.PushAny(out.Interface())
	}

	// returning L.GetTop() helps to account for Go functions which would have
	//  both (or either) returned values via Go **AND** via C L.Push*() API...
	//
	//return len(outValues)
	return L.GetTop()
}
