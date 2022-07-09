package luajit

import (
	"fmt"
	"unsafe"
)

func (L *State) _getGlobalAny(v any) {
	L.PushOne(v)
	L.Y_lua_gettable(GlobalsIndex)
}

// GetGlobal pushes _G[v] on the stack
func (L *State) GetGlobal(v any) {
	L._getGlobalAny(v)
}

func (L *State) GetGlobalAny(v any) any {
	L._getGlobalAny(v)
	defer L.Pop(1)
	return L.ToAny(-1)
}

func (L *State) GetGlobalBoolean(v any) bool {
	L._getGlobalAny(v)
	defer L.Pop(1)
	return L.ToBool(-1)
}

func (L *State) GetGlobalString(v any) string {
	L._getGlobalAny(v)
	defer L.Pop(1)
	return L.ToString(-1)
}

func (L *State) GetGlobalPointer(v any) (ptr unsafe.Pointer) {
	L._getGlobalAny(v)
	defer L.Pop(1)
	defer func() {
		if ptr == nil {
			fmt.Println("\033[31mWARNING\033[39m: GetGlobalPointer returned nil pointer")
		}
	}()
	return L.ToPointer(-1)
}

func (L *State) GetGlobalTable(v any) *Table {
	L._getGlobalAny(v)
	defer L.Pop(1)
	return L.ToTable(-1)
}

// SetGlobal set _G[k] = topOfTheStack and pops the stack.
func (L *State) SetGlobal(k any) {
	L.PushOne(k)
	L.insert(-2) // move key before value
	L.Y_lua_settable(GlobalsIndex)
}

func (L *State) SetGlobalAny(k any, v any) {
	L.PushOne(k)
	L.PushOne(v)
	L.Y_lua_settable(GlobalsIndex)
}

func (L *State) SetGlobalString(name string, value string) {
	L.SetGlobalAny(name, value)
}
