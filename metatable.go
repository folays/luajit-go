package luajit

import (
	"fmt"
	"log"
	"unsafe"
)

/*
#include <stdlib.h>
#include <lua.h>

#include "metatable.h"
*/
import "C"

type metatable struct {
	ref Ref

	cb_new metatableCbNew
	cb_gc  metatableCbGc
}

type metatableCbNew func(mt *metatable, values ...any)
type metatableCbGc func(mt *metatable, ptr unsafe.Pointer)

func (L *State) metatable_prepare() {
	L.all_mt_by_ref = make(map[Ref]*metatable)
	L.all_mt_by_ptr = make(map[unsafe.Pointer]*metatable)
}

// create a new metatable (which won't be left on the stack, but instead registry-ed)
func (L *State) metatable_newRef() (mt *metatable) {
	mt = &metatable{}

	L.NewTable()
	mt.ref = L.RegistryRef()

	L.metatable_setFieldClosure(mt, "__gc", LuaCFunction(C.metatable_gc_c))

	L.all_mt_by_ref[mt.ref] = mt
	L.all_mt_by_ptr[L.metatable_ptr(mt)] = mt

	return
}

// return the address of Lua metatable
func (L *State) metatable_ptr(mt *metatable) unsafe.Pointer {
	L.RegistryGet(mt.ref)
	defer L.Pop(1)

	return L.ToPointer(-1)
}

func (L *State) metatable_get_by_ptr(ptr unsafe.Pointer) *metatable {
	return L.all_mt_by_ptr[ptr]
}

// push the metatable onto the stack.
func (L *State) metatable_push(mt *metatable) {
	L.RegistryGet(mt.ref)
}

func (L *State) metatable_cb_set_new(mt *metatable, cb_fn_new metatableCbNew) {
	mt.cb_new = cb_fn_new
}

func (L *State) metatable_cb_set_gc(mt *metatable, cb_fn_gc metatableCbGc) {
	mt.cb_gc = cb_fn_gc
}

// SetField does `mt[event] = metamethod` where `metamethod` is at the top of the stack.
// As in Lua, pops the top of the stack (`metamethod`).
func (L *State) metatable_setField(mt *metatable, event string) {
	L.RegistryGet(mt.ref)
	defer L.Pop(1)

	L.insert(-2) // move table before value

	L.SetField(-2, event)
}

// SetFieldAny does `mt[event] = metamethod`.
// Does not modify the stack.
func (L *State) metatable_setFieldAny(mt *metatable, event string, metamethod any) {
	L.RegistryGet(mt.ref)
	defer L.Pop(1)

	L.SetFieldAny(-1, event, metamethod)
}

// SetFieldAny does `mt[event] = metamethod`.
//
//	Does not modify the stack.
//	Will always have #1     upvalue  as `*State`
//	Will append      #2...n upvalues as given parameters
func (L *State) metatable_setFieldClosure(mt *metatable, event string, metamethod LuaCFunction, upvalues ...any) {
	L.RegistryGet(mt.ref)
	defer L.Pop(1)

	L.pushUpvalueState()
	nbArgs := L.PushMultiple(upvalues)

	L.pushcclosure(metamethod, 1+nbArgs)

	L.SetField(-2, event)
}

func (L *State) metatable_get_by_index(index Index) (mt *metatable) {
	if L.IsType(index, LUA_TUSERDATA) == false {
		log.Fatalf("(*State).metatable_get_by_index() error : not a fulluserdata")
	}

	if L.Y_lua_getmetatable(index) == false {
		log.Fatalf("(*State).metatable_get_by_index() error : no metatable")
	}
	defer L.Pop(1)

	if mt = L.metatable_get_by_ptr(L.ToPointer(-1)); mt == nil {
		log.Fatalf("(*State).metatable_get_by_index() error : unknown metatable")
	}

	return
}

func (L *State) metatable_constructor(mt *metatable, values ...any) {
	if mt.cb_new == nil {
		log.Fatalf("\033[31mFATAL\033[39m (*State).metatable_constructor(,%T) : metatable has no __new callback\n", values)
	}

	mt.cb_new(mt, values...)

	L.metatable_push(mt)
	L.Y_lua_setmetatable(-2)
}

//export metatable_gc_go
func metatable_gc_go(l LuaStatePtr) int {
	var (
		L   = toUpvalueState(l, 1)
		mt  = L.metatable_get_by_index(1)
		ptr = L.ToUserdata_raw(1)
	)

	switch mt.cb_gc != nil {
	case false:
		fmt.Printf("\033[31mWARNING\033[39m luajit.metatable_gc_go() : metatable has no __gc callback\n")
	case true:
		mt.cb_gc(mt, ptr)
	}

	return 0
}
