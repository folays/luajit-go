package luajit

import (
	"runtime"
	"unsafe"
)

/*
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <err.h>
#include <lua.h>
#include <lauxlib.h>

int ylua_gettop(uintptr_t L) { return lua_gettop((lua_State *)L); }
void ylua_settop(uintptr_t L, int index) { lua_settop((lua_State *)L, index); }

void ylua_insert(uintptr_t L, int index) { lua_insert((lua_State *)L, index); }

void ylua_remove(uintptr_t L, int index) { lua_remove((lua_State *)L, index); }

void ylua_pop(uintptr_t L, int n) { lua_pop((lua_State *)L, n); }

int ylua_type(uintptr_t L, int index) { return lua_type((lua_State *)L, index); }

int ylua_isfunction(uintptr_t L, int index) { return lua_isfunction((lua_State *)L, index); }

int ylua_rawequal(uintptr_t L, int index1, int index2) { return lua_rawequal((lua_State *)L, index1, index2); }

void ylua_rawgeti(uintptr_t L, int index, int n) { lua_rawgeti((lua_State *)L, index, n); }

int ylua_toboolean(uintptr_t L, int index) { return lua_toboolean((lua_State *)L, index); }

int ylua_tointeger(uintptr_t L, int index) { return lua_tointeger((lua_State *)L, index); }

lua_Number ylua_tonumber(uintptr_t L,int index) { return lua_tonumber((lua_State *)L, index); }

//#cgo noescape ylua_tolstring
const char *ylua_tolstring(uintptr_t L, int index, size_t *len) { return lua_tolstring((lua_State *)L, index, len); }

const void *ylua_topointer(uintptr_t L, int index) { return lua_topointer((lua_State *)L, index); }

void *ylua_touserdata(uintptr_t L, int index) { return lua_touserdata((lua_State *)L, index); }

//#cgo noescape ylua_getlfield
void ylua_getlfield(uintptr_t L, int index, const char *k, size_t l)
{
    lua_pushlstring((lua_State *)L, k, l);
	if (index < 0 && index > LUA_REGISTRYINDEX)
		index--;
	lua_gettable((lua_State *)L, index);
}

void ylua_pushvalue(uintptr_t L, int index) { lua_pushvalue((lua_State *)L, index); }

void ylua_pushnil(uintptr_t L) { lua_pushnil((lua_State *)L); }

void ylua_pushboolean(uintptr_t L, int b) { lua_pushboolean((lua_State *)L, b); }

void ylua_pushinteger(uintptr_t L, lua_Integer n) { lua_pushinteger((lua_State *)L, n); }

void ylua_pushnumber(uintptr_t L, double n) { lua_pushnumber((lua_State *)L, n); }

//#cgo noescape ylua_pushlstring
void ylua_pushlstring(uintptr_t L, const char *s, size_t l) { lua_pushlstring((lua_State *)L, s, l); }

void ylua_pushlightuserdata(uintptr_t L, void *p) { lua_pushlightuserdata((lua_State *)L, p); }

void ylua_pushcclosure(uintptr_t L, lua_CFunction fn, int n) { lua_pushcclosure((lua_State *)L, fn, n); }

void ylua_pushcfunction(uintptr_t L, lua_CFunction f) { lua_pushcfunction((lua_State *)L, f); }

void ylua_newtable(uintptr_t L) { lua_newtable((lua_State *)L); }

void ylua_gettable(uintptr_t L, int index) { lua_gettable((lua_State *)L, index); }

void ylua_settable(uintptr_t L, int index) { lua_settable((lua_State *)L, index); }

int ylua_getmetatable(uintptr_t L, int index) { return lua_getmetatable((lua_State *)L, index); }

int ylua_setmetatable(uintptr_t L, int index) { return lua_setmetatable((lua_State *)L, index); }

void *ylua_newuserdata(uintptr_t L, size_t size) { return lua_newuserdata((lua_State *)L, size); }

int ylua_next(uintptr_t L, int index) { return lua_next((lua_State *)L, index); }

int yluaL_ref(uintptr_t L, int t) { return luaL_ref((lua_State *)L, t); }

void yluaL_unref(uintptr_t L, int t, int ref) { return luaL_unref((lua_State *)L, t, ref); }

static const char *_ylua_bdup_cstring(const char *buf, size_t l) {
	char *buf2;

	buf2 = malloc(l + 1);
	if (!buf2)
		err(1, "luajt-go _ylua_bdup() malloc(%ld) failed", l + 1);
	memcpy(buf2, buf, l);
	buf2[l] = '\0'; // \0-"terminate" the string

	return buf2;
}

//#cgo noescape yluaL_loadbuffer
int yluaL_loadbuffer(uintptr_t L, const char *buff, size_t sz, const char *name, size_t namelen) {
	const char *name2 = _ylua_bdup_cstring(name, namelen);

	return luaL_loadbuffer((lua_State *)L, buff, sz, name2);

	free((void *)name2);
}

//#cgo noescape yluaL_loadlstring
int yluaL_loadlstring(uintptr_t L, const char *s, size_t l) {
    return yluaL_loadbuffer(L, s, l, s, l);
}

//#cgo noescape yluaL_loadlfile
int yluaL_loadlfile(uintptr_t L, const char *filename, size_t l) {
	const char *filename2 = _ylua_bdup_cstring(filename, l);

	return luaL_loadfile((lua_State *)L, filename2);

	free((void *)filename2);
}

int ylua_upvalueindex(int i) { return lua_upvalueindex(i); }

int ylua_pcall(uintptr_t L, int nargs, int nresults, int errfunc) { return lua_pcall((lua_State *)L, nargs, nresults, errfunc); }
*/
import "C"

func (L *State) Y_lua_gettop() (index Index) {
	return Index(C.ylua_gettop(L.stateU))
}

func (L *State) Y_lua_settop(index Index) {
	C.ylua_settop(L.stateU, C.int(index))
}

func (L *State) Y_lua_insert(index Index) {
	C.ylua_insert(L.stateU, C.int(index))
}

func (L *State) Y_lua_remove(index Index) {
	C.ylua_remove(L.stateU, C.int(index))
}

func (L *State) Y_lua_pop(n int) {
	C.ylua_pop(L.stateU, C.int(n))
}

func y_lua_type(l LuaStatePtr, index Index) LuaType {
	lU := C.uintptr_t(uintptr(unsafe.Pointer(l)))
	return LuaType(C.ylua_type(lU, C.int(index)))
}

func (L *State) Y_lua_type(index Index) LuaType {
	return LuaType(C.ylua_type(L.stateU, C.int(index)))
}

func (L *State) Y_lua_isfunction(index Index) bool {
	return C.ylua_isfunction(L.stateU, C.int(index)) != 0
}

func (L *State) Y_lua_rawequal(index1, index2 Index) bool {
	return C.ylua_rawequal(L.stateU, C.int(index1), C.int(index2)) != 0
}

func (L *State) Y_lua_rawgeti(index Index, n Integer) {
	C.ylua_rawgeti(L.stateU, C.int(index), C.int(n))
}

func (L *State) Y_lua_toboolean(index Index) bool {
	return C.ylua_toboolean(L.stateU, C.int(index)) != 0
}

func (L *State) Y_lua_tointeger(index Index) Integer {
	return Integer(C.ylua_tointeger(L.stateU, C.int(index)))
}

func (L *State) Y_lua_tonumber(index Index) Number {
	return Number(C.ylua_tonumber(L.stateU, C.int(index)))
}

func (L *State) Y_lua_tostring(index Index) string {
	var (
		s    *C.char
		_len C.size_t
	)

	s = C.ylua_tolstring(L.stateU, C.int(index), &_len)

	// https://github.com/golang/go/issues/61361#issuecomment-1638741918
	//return C.GoStringN(s, C.int(_len))
	return string(unsafe.Slice((*byte)(unsafe.Pointer(s)), int(_len)))
}

func (L *State) Y_lua_topointer(index Index) unsafe.Pointer {
	return C.ylua_topointer(L.stateU, C.int(index))
}

func y_lua_touserdata(l LuaStatePtr, index Index) unsafe.Pointer {
	lU := C.uintptr_t(uintptr(unsafe.Pointer(l)))
	return C.ylua_touserdata(lU, C.int(index))
}

func (L *State) Y_lua_touserdata(index Index) unsafe.Pointer {
	return C.ylua_touserdata(L.stateU, C.int(index))
}

// Trick escape analysis to consider `s` as `noescape`. Invalid with a moving GC.
// XXX: Go 1.22 :
//   - directly pass (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
//   - replace the ptr->uintptr->ptr trick by a #cgo noescape <func> annotation.
//   - remove runtime.KeepAlive(s) in callers
//
// DO NOT forget to add runtime.KeepAlive(s) in callers! (which trick non-moving GC)
func _trick_cgo_escape_analysis_string(s string) (sC *C.char) {
	sU := uintptr(unsafe.Pointer(unsafe.StringData(s)))
	sC = (*C.char)(unsafe.Pointer(sU))
	return
}
func _trick_cgo_escape_analysis_bytes(buf []byte) (bC *C.char) {
	bU := uintptr(unsafe.Pointer(unsafe.SliceData(buf)))
	bC = (*C.char)(unsafe.Pointer(bU))
	return
}

func (L *State) Y_lua_getfield(index Index, key string) {
	C.ylua_getlfield(L.stateU, C.int(index), _trick_cgo_escape_analysis_string(key), C.ulong(len(key)))
	runtime.KeepAlive(key)
}

func (L *State) Y_lua_pushvalue(index Index) {
	C.ylua_pushvalue(L.stateU, C.int(index))
}

func (L *State) Y_lua_pushnil() {
	C.ylua_pushnil(L.stateU)
}

func (L *State) Y_lua_pushboolean(b bool) {
	switch b {
	case false:
		C.ylua_pushboolean(L.stateU, C.int(0))
	case true:
		C.ylua_pushboolean(L.stateU, C.int(1))
	}
}

func (L *State) Y_lua_pushinteger(n Integer) {
	C.ylua_pushinteger(L.stateU, C.long(n))
}

func (L *State) Y_lua_pushnumber(n Number) {
	C.ylua_pushnumber(L.stateU, C.double(n))
}

func (L *State) Y_lua_pushstring(s string) {
	C.ylua_pushlstring(L.stateU, _trick_cgo_escape_analysis_string(s), C.ulong(len(s)))
	runtime.KeepAlive(s)
}

func (L *State) Y_lua_pushlightuserdata(ptr unsafe.Pointer) {
	C.ylua_pushlightuserdata(L.stateU, ptr)
}

func (L *State) Y_lua_pushcclosure(fn LuaCFunction, nbUpvalues int) {
	C.ylua_pushcclosure(L.stateU, fn, C.int(nbUpvalues))
}

func (L *State) Y_lua_pushcfunction(fn LuaCFunction) {
	C.ylua_pushcfunction(L.stateU, fn)
}

func (L *State) Y_lua_newtable() {
	C.ylua_newtable(L.stateU)
}

func (L *State) Y_lua_gettable(index Index) {
	C.ylua_gettable(L.stateU, C.int(index))
}

func (L *State) Y_lua_settable(index Index) {
	C.ylua_settable(L.stateU, C.int(index))
}

func (L *State) Y_lua_getmetatable(index Index) bool {
	return C.ylua_getmetatable(L.stateU, C.int(index)) != 0
}

func (L *State) Y_lua_setmetatable(index Index) {
	C.ylua_setmetatable(L.stateU, C.int(index))
}

func (L *State) Y_lua_newuserdata(size int) unsafe.Pointer {
	return C.ylua_newuserdata(L.stateU, C.ulong(size))
}

func (L *State) Y_lua_next(index Index) Index {
	return Index(C.ylua_next(L.stateU, C.int(index)))
}

func (L *State) Y_luaL_ref(t Index) Ref {
	return Ref(C.yluaL_ref(L.stateU, C.int(t)))
}

func (L *State) Y_luaL_unref(t Index, ref Ref) {
	C.yluaL_unref(L.stateU, C.int(t), C.int(ref))
}

func (L *State) Y_luaL_loadbuffer(buf []byte, chunkName string) int {
	defer runtime.KeepAlive(buf)
	defer runtime.KeepAlive(chunkName)
	return int(C.yluaL_loadbuffer(L.stateU, _trick_cgo_escape_analysis_bytes(buf), C.ulong(len(buf)), _trick_cgo_escape_analysis_string(chunkName), C.ulong(len(chunkName))))
}

func (L *State) Y_luaL_loadstring(s string) int {
	defer runtime.KeepAlive(s)
	return int(C.yluaL_loadlstring(L.stateU, _trick_cgo_escape_analysis_string(s), C.ulong(len(s))))
}

func (L *State) Y_luaL_loadfile(path string) int {
	defer runtime.KeepAlive(path)
	return int(C.yluaL_loadlfile(L.stateU, _trick_cgo_escape_analysis_string(path), C.ulong(len(path))))
}

func y_lua_upvalueindex(i UpvalueIndex) Index {
	return GlobalsIndex - Index(i)
	//return Index(C.ylua_upvalueindex(C.int(i)))
}

func (L *State) Y_lua_pcall(nbArgs int, nbResults int, indexErrFunc Index) int {
	return int(C.ylua_pcall(L.stateU, C.int(nbArgs), C.int(nbResults), C.int(indexErrFunc)))
}
