package luajit

import (
	"unsafe"
)

/*
#include <stdlib.h>
#include <lua.h>
*/
import "C"

type Index C.int

type UpvalueIndex C.int

type Integer C.int
type Number C.double

const (
	RegistryIndex = Index(C.LUA_REGISTRYINDEX)
	GlobalsIndex  = Index(C.LUA_GLOBALSINDEX)
)

func (L *State) GetTop() Index {
	return L.Y_lua_gettop()
}

func (L *State) SetTop(index Index) {
	L.Y_lua_settop(index)
}

func (L *State) Pop(n int) {
	L.Y_lua_pop(n)
}

func (L *State) pushvalue(index Index) {
	L.Y_lua_pushvalue(index)
}

func (L *State) pushlightuserdata(ptr unsafe.Pointer) {
	L.Y_lua_pushlightuserdata(ptr)
}

func (L *State) newuserdata(size int) (ptr unsafe.Pointer) {
	return L.Y_lua_newuserdata(size)
}

func (L *State) newuserdata_Type(uType any) (ptr unsafe.Pointer) {
	return L.Y_lua_newuserdata(int(unsafe.Sizeof(uType)))
}

func (L *State) pushcclosure(fn LuaCFunction, nUpvalues int) {
	L.Y_lua_pushcclosure(fn, nUpvalues)
}

func (L *State) remove(index Index) {
	L.Y_lua_remove(index)
}

func (L *State) insert(index Index) {
	L.Y_lua_insert(index)
}

// convert from relative to absolute
func (L *State) abs_index(index Index) Index {
	if index > 0 || index <= RegistryIndex {
		return index
	}
	return Index(L.GetTop()) + index + 1
}

func upvalueindex(i UpvalueIndex) Index {
	return y_lua_upvalueindex(i)
}

func (L *State) NewTable() {
	L.Y_lua_newtable()
}
