package luajit

import (
	"fmt"
	"log"
	"reflect"
	"unsafe"
)

func (L *State) PushMultiple(args []any) (nbArgs int) {
	for _, arg := range args {
		nbArgs += L.PushAny(arg)
	}
	return
}

func (L *State) PushOne(arg any) {
	nbArgs := L.PushAny(arg)
	switch {
	case nbArgs < 1:
		log.Fatalf("(*State).PushOne() : sub- PushAny() returned < 1")
	case nbArgs == 1:
	case nbArgs > 1:
		L.Pop(nbArgs - 1)
	}
}

func (L *State) PushAny(arg any) (nbArgs int) {
	{
		switch v := arg.(type) {
		case nil:
			L.Y_lua_pushnil()
		case string:
			L.Y_lua_pushstring(v)
		case []byte:
			if v == nil {
				L.Y_lua_pushnil()
			} else {
				L.Y_lua_pushstring(string(v)) // relies on `v` to noescape
			}
		case int:
			L._push_int64(int64(v))
		case uint:
			L._push_uint64(uint64(v))
		case int8:
			L.Y_lua_pushinteger(Integer(v))
		case int16:
			L.Y_lua_pushinteger(Integer(v))
		case int32:
			L.Y_lua_pushinteger(Integer(v))
		case int64:
			L._push_int64(v)
		case uint8:
			L.Y_lua_pushinteger(Integer(v))
		case uint16:
			L.Y_lua_pushinteger(Integer(v))
		case uint32:
			L.Y_lua_pushinteger(Integer(v))
		case uint64:
			L._push_uint64(v)
		case float32:
			L.Y_lua_pushnumber(Number(v))
		case float64:
			L.Y_lua_pushnumber(Number(v))
		case bool:
			L.Y_lua_pushboolean(v)
		case unsafe.Pointer:
			L.pushlightuserdata(v)
		case LuaCFunction:
			L.Y_lua_pushcfunction(v)
		case Ref:
			L.RegistryGet(v)
		case []any:
			nbArgs += -1 + L.PushMultiple(v)
		case error:
			L.PushOne(v.Error())
		default:
			switch reflect.TypeOf(v).Kind() {
			case reflect.Func:
				L.pushGoFunction(reflect.ValueOf(v))
			case reflect.Struct:
				L.pushGoStruct(reflect.ValueOf(v))
			case reflect.Pointer:
				L.pushGoPointer(reflect.ValueOf(v))
			case reflect.Slice:
				L.pushSlice(reflect.ValueOf(v))
			default:
				fmt.Printf("\033[31mERROR\033[39m PushAny(): unknown arg type %T\n", arg)
				L.Y_lua_pushnil()
				panic(fmt.Sprintf("PushAny(): unknown arg type %T", arg))
			}
		}
	}
	nbArgs++
	return
}

func (L *State) _push_int64(v int64) {
	switch v > -(1<<50) && v < +(1<<50) {
	case true:
		L.Y_lua_pushinteger(Integer(v))
	case false:
		L._push_cdata_int_str(fmt.Sprintf("return %dLL", v))
	}
}

func (L *State) _push_uint64(v uint64) {
	switch v < +(1 << 50) {
	case true:
		L.Y_lua_pushinteger(Integer(v))
	case false:
		L._push_cdata_int_str(fmt.Sprintf("return %dULL", v))
	}
}

func (L *State) pushNil() {
	L.Y_lua_pushnil()
}

func (L *State) pushSlice(v reflect.Value) {
	if v.IsNil() {
		L.pushNil()
		return
	}

	L.NewTable()
	for i := 0; i < v.Len(); i++ {
		L.SetFieldAny(-1, i+1, v.Index(i).Interface())
	}
}
