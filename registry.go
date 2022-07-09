package luajit

/*
#include <lua.h>
#include <lauxlib.h>
*/
import "C"

type Ref C.int
type RefFunction Ref

const (
	RefNil Ref = C.LUA_REFNIL
	NoRef  Ref = C.LUA_NOREF
)

// RegistryRef remove the top
func (L *State) RegistryRef() (ref Ref) {
	return L.Y_luaL_ref(RegistryIndex)
}

// RegistryRefIndex do *not* remove the index
func (L *State) RegistryRefIndex(index Index) Ref {
	L.Y_lua_pushvalue(index)
	return L.RegistryRef()
}

func (L *State) RegistryUnref(ref Ref) {
	L.Y_luaL_unref(RegistryIndex, ref)
}

// RegistryGet push LUA_REGISTRY_INDEX[ref] onto the stack
func (L *State) RegistryGet(ref Ref) {
	L.Y_lua_rawgeti(RegistryIndex, Integer(ref))
}
