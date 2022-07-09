package luajit

import (
	"fmt"
	"log"
	"runtime/cgo"
	"unsafe"
)

/*
#include <stdlib.h>
#include <lua.h>
*/
import "C"

// TODO : finalizer if last reference goes away
type userdata struct {
	ref Ref
}

type userdata_region [1]C.uintptr_t

// create a new userdata (which won't be left on the stack, but instead registry-ed)
func (L *State) userdata_newRef(mt *metatable, values ...any) (u *userdata) {
	u = &userdata{}

	L.userdata_new(mt, values...)
	u.ref = L.RegistryRef()

	return
}

// create a new userdata, pushes it on the stack
func (L *State) userdata_new(mt *metatable, values ...any) {
	L.metatable_constructor(mt, values...)
}

// push the userdata onto the stack.
func (L *State) userdata_push(u *userdata) {
	L.RegistryGet(u.ref)
}

func (L *State) userdata_check(index Index, mt *metatable) (v any, err error) {
	index = L.abs_index(index)

	switch L.Type(index) {
	case LUA_TUSERDATA:
	default:
		fmt.Println("\033[31mWARNING\033[39m: userdata_check() error : not a fulluserdata")
		return
	}

	if L.Y_lua_getmetatable(index) == false {
		log.Println("(*State).userdata_check() error : no metatable")
		return nil, fmt.Errorf("no metatable")
	}
	L.metatable_push(mt)
	defer L.Pop(2)

	switch L.Y_lua_rawequal(-2, -1) {
	case true:
		v = L._userdata_to_value(index, mt)
		return v, nil
	default: // false
		return nil, fmt.Errorf("metatable does not match")
	}
}

func (L *State) userdata_checkFatal(index Index, mt *metatable) (v any) {
	var (
		err error
	)

	if v, err = L.userdata_check(index, mt); err != nil {
		log.Fatalf("(*State).userdata_checkFatal() error : %s", err)
	}

	return
}
func (L *State) _userdata_new(mt *metatable, values ...any) {
	h := cgo.NewHandle(values)

	udata := (*userdata_region)(L.newuserdata_Type(userdata_region{}))
	udata[0] = C.uintptr_t(h)
}

func (L *State) _userdata_gc(mt *metatable, ptr unsafe.Pointer) {
	{
		var (
			udata = (*userdata_region)(ptr)
			h     = cgo.Handle(udata[0])
			v     = h.Value()
		)

		fmt.Printf("\033[35mUSERDATA COLLECT\033[39m on %T %v\n", v, v)

		h.Delete()
	}
}

// BEWARE : caller have FULL responsibility of checking that `index` holds one of OUR `userdata`
func (L *State) _userdata_to_values(index Index, mt *metatable) []any {
	var (
		udata = (*userdata_region)(L.ToUserdata_raw(index))
		h     = cgo.Handle(udata[0])
		v     = h.Value()
	)
	return v.([]any)
}

// BEWARE : caller have FULL responsibility of checking that `index` holds one of OUR `userdata`
func (L *State) _userdata_to_value(index Index, mt *metatable) any {
	return L._userdata_to_values(index, mt)[0]
}

func (L *State) ToUserdata(index Index) (v any) {
	var (
		mt = L.metatable_get_by_index(index)
	)

	v = L._userdata_to_value(index, mt)

	return
}

func (L *State) ToPayload(index Index) (values []any) {
	var (
		mt = L.metatable_get_by_index(index)
	)

	values = L._userdata_to_values(index, mt)
	if len(values) >= 1 {
		values = values[1:] // skip #1 value (the userdata itself)
	}

	return
}
