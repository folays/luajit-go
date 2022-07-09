package luajit

import (
	"embed"
	"os"
	"runtime"
)

//go:embed defaults.lua
var luaDefaults []byte

//go:embed lua-modules
//go:embed lua-helpers
var luaFS embed.FS

var silence = false

func init() {
	_, silence = os.LookupEnv("SILENCE")
}

func (L *State) ReasonableDefaults() (err error) {
	{
		pc, _, _, _ := runtime.Caller(0)
		L.myself_defaults = runtime.FuncForPC(pc).Name()
	}

	//L.OpenLibs()

	L.SetGlobalAny("luajit_embed", luaFS)

	//L.OpenLib("debug", LuaOpen_debug)
	L.OpenLib("", LuaOpen_base)
	L.OpenLib("package", LuaOpen_package)
	L.OpenLib("table", LuaOpen_table)
	L.OpenLib("string", LuaOpen_string)
	L.Module_preload("math", LuaOpen_math)
	L.Module_preload("io", LuaOpen_io)
	L.Module_preload("os", LuaOpen_os)
	L.Module_preload("bit", LuaOpen_bit)
	L.Module_preload("jit", LuaOpen_jit)
	L.Module_preload("ffi", LuaOpen_ffi)
	L.Module_preload("string.buffer", LuaOpen_string_buffer)

	L.RunChunkBufferFatal(luaDefaults, "defaults.lua")

	L.Load_C_Helpers()
	L.RunEmbedFsPathFatal(luaFS, "lua-helpers")

	L.SetGlobalString("LUAJIT_SOURCE_DIR", GetCallerSourceDir(0))
	L.SetGlobalString("LUAJIT_CALLER_DIR", GetCallerSourceDir(1))

	return
}
