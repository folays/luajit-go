package luajit

import (
	"log"
)

/*
#include <stdlib.h>
#include <lua.h>
#include <lualib.h>
#include <lauxlib.h>
*/
import "C"

var (
	LuaOpen_base    = LuaCFunction(C.luaopen_base)
	LuaOpen_package = LuaCFunction(C.luaopen_package)
	LuaOpen_string  = LuaCFunction(C.luaopen_string)
	LuaOpen_table   = LuaCFunction(C.luaopen_table)
	LuaOpen_math    = LuaCFunction(C.luaopen_math)
	LuaOpen_io      = LuaCFunction(C.luaopen_io)
	LuaOpen_os      = LuaCFunction(C.luaopen_os)
	LuaOpen_debug   = LuaCFunction(C.luaopen_debug)

	//
	//LuaOpen_coroutine = LuaCFunction(C.luaopen_coroutine)

	LuaOpen_bit           = LuaCFunction(C.luaopen_bit)
	LuaOpen_jit           = LuaCFunction(C.luaopen_jit)
	LuaOpen_ffi           = LuaCFunction(C.luaopen_ffi)
	LuaOpen_string_buffer = LuaCFunction(C.luaopen_string_buffer)
)

func (L *State) OpenLib(moduleName string, fn LuaCFunction) {
	var (
		err error
	)

	//return L.RunCodeFatal(`({...})[1](({...})[2])`, fn, moduleName)

	if err = L.RunClosure(fn, moduleName); err != nil {
		log.Fatalf("(*State).OpenLib(%s) error : %s", moduleName, err)
	}
}

// OpenLibs loads lua libraries
func (L *State) OpenLibs() {
	C.luaL_openlibs(L.state)
}
