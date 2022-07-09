package luajit

import (
	"log"
)

/*
#include <stdlib.h>
#include <lua.h>
*/
import "C"

type LuaType int

const (
	LUA_TNONE          LuaType = C.LUA_TNONE
	LUA_TNIL           LuaType = C.LUA_TNIL
	LUA_TBOOLEAN       LuaType = C.LUA_TBOOLEAN
	LUA_TLIGHTUSERDATA LuaType = C.LUA_TLIGHTUSERDATA
	LUA_TNUMBER        LuaType = C.LUA_TNUMBER
	LUA_TSTRING        LuaType = C.LUA_TSTRING
	LUA_TTABLE         LuaType = C.LUA_TTABLE
	LUA_TFUNCTION      LuaType = C.LUA_TFUNCTION
	LUA_TUSERDATA      LuaType = C.LUA_TUSERDATA
	LUA_TTHREAD        LuaType = C.LUA_TTHREAD
	LUA_TCDATA         LuaType = 10
)

func (typ LuaType) String() string {
	switch typ {
	case LUA_TNONE:
		return "none"
	case LUA_TNIL:
		return "nil"
	case LUA_TNUMBER:
		return "number"
	case LUA_TBOOLEAN:
		return "boolean"
	case LUA_TSTRING:
		return "string"
	case LUA_TTABLE:
		return "table"
	case LUA_TFUNCTION:
		return "function"
	case LUA_TUSERDATA:
		return "userdata"
	case LUA_TTHREAD:
		return "thread"
	case LUA_TLIGHTUSERDATA:
		return "lightuserdata"
	case LUA_TCDATA:
		return "cdata"
	default:
		log.Fatalf("luajit.(LuaType).String() error : unhandled LuaType %d", typ)
	}
	return "?"
}

func (L *State) Type(index Index) LuaType {
	return L.Y_lua_type(index)
}

func (L *State) TypeName(index Index) string {
	return L.Type(index).String()
}

func (L *State) IsType(index Index, typ LuaType) bool {
	return L.Type(index) == typ
}

func (L *State) IsNil(index Index) bool {
	return L.IsType(index, LUA_TNIL)
}

func (L *State) IsTable(index Index) bool {
	return L.IsType(index, LUA_TTABLE)
}
