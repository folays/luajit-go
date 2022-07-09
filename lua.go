package luajit

import (
	"reflect"
	"runtime/cgo"
	"unsafe"
)

/*
#cgo pkg-config: luajit
#cgo CFLAGS: -g1
#include <stdlib.h>
#include <lua.h>
#include <lauxlib.h>
*/
import "C"

// LuaStatePtr is a type to represent `struct lua_State`
type LuaStatePtr *C.struct_lua_State

// State stores lua state
type State struct {
	state     LuaStatePtr
	stateU    C.uintptr_t
	stackTops []Index // stores the lua_gettop() before each call (minus: f and nbArgs).

	myself_pkg      string
	myself_defaults string

	h_State cgo.Handle

	error_ref RefFunction

	state_fac *factory  // *State itself (factory metatable)
	state_u   *userdata // *State itself (userdata)

	all_mt_by_ref map[Ref]*metatable            // by Lua's registry reference
	all_mt_by_ptr map[unsafe.Pointer]*metatable // by Lua's lua_topointer()

	function_fac *factory

	function_everywhere map[reflect.Type]*funcEverywhere

	all_br_by_type_bridge map[reflect.Type]*Bridge
	all_br_by_type_shared map[reflect.Type]*Bridge
}

// NewState creates new Lua state
func NewState() (L *State) {
	L = &State{
		state:      C.luaL_newstate(),
		myself_pkg: reflect.TypeOf(State{}).PkgPath(),
		error_ref:  RefFunction(RefNil),
	}
	L.stateU = C.uintptr_t(uintptr(unsafe.Pointer(L.state)))
	L.run_func_prepare()

	//L.OpenLibs()

	L.h_State = cgo.NewHandle(L)

	L.metatable_prepare()

	L.state_fac = L.factory_new()
	L.state_u = L.factory_produceRef(L.state_fac, L)

	L.function_prepare_metatable()

	L.function_prepare_everywhere()

	L.bridge_prepare()

	// default bridges
	L.embed_prepare_bridge()

	return
}

// Close destroys lua state
func (L *State) Close() {
	L.function_unprepare()

	L.h_State.Delete()
	L.h_State = 0

	C.lua_close(L.state)
}

func (L *State) errorStringPop() string {
	s := L.Y_lua_tostring(-1)
	L.Pop(1) // remove the error
	return s
}

func (L *State) CollectGarbage() {
	L.RunString("collectgarbage()")
}

func (L *State) GC() {
	L.CollectGarbage()
}
