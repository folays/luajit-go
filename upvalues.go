package luajit

import (
	"fmt"
	"runtime/cgo"
	"unsafe"
)

// TODO :
// - prepare (only once) a registry-referenced fulluserdata
// - associate a metatable
// - push this fulluserdata instead
// Benefit : being able to check that its really a *lua_State in toUpvalueState()
func (L *State) pushUpvalueState() {
	L.pushlightuserdata(unsafe.Pointer(L.h_State))
}

func toUpvalueState(l LuaStatePtr, i UpvalueIndex) (L *State) {
	var (
		index = upvalueindex(i)
	)

	switch y_lua_type(l, index) {
	case LUA_TLIGHTUSERDATA:
	default:
		fmt.Println("\033[31mWARNING\033[39m toUpvalueState() error : not a lightuserdata")
	}

	h := cgo.Handle(y_lua_touserdata(l, index))

	return h.Value().(*State)
}
