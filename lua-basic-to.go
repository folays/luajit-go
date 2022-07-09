package luajit

import (
	"log"
	"unsafe"
)

// ToBool can modify value
func (L *State) ToBool(index Index) bool {
	return L.Y_lua_toboolean(index)
}

// ToInt can modify value
func (L *State) ToInt(index Index) int {
	return int(L.Y_lua_tointeger(index))
}

func (L *State) ToIntSafe(index Index) (n int) {
	L.pushvalue(index)
	defer L.Pop(1)
	return L.ToInt(-1)
}

// ToFloat64 can modify value
func (L *State) ToFloat64(index Index) (n float64) {
	return float64(L.Y_lua_tonumber(index))
}

// ToString can modify value
func (L *State) ToString(index Index) (s string) {
	return L.Y_lua_tostring(index)
}

// ToStringSafe will dup value if needed before ToString()
func (L *State) ToStringSafe(index Index) (s string) {
	L.pushvalue(index)
	defer L.Pop(1)
	return L.ToString(-1)
}

func (L *State) _toCString(index Index) (s string) {
	return L.Y_lua_tostring(index)
}

func (L *State) ToPointer(index Index) unsafe.Pointer {
	return L.Y_lua_topointer(index)
}

func (L *State) ToUserdata_raw(index Index) unsafe.Pointer {
	return L.Y_lua_touserdata(index)
}

func (L *State) ToAny(index Index) any {
	switch typ := L.Type(index); typ {
	case LUA_TNIL:
		return nil
	case LUA_TNUMBER:
		return L.ToFloat64(index)
	case LUA_TBOOLEAN:
		return L.ToBool(index)
	case LUA_TSTRING:
		return L.ToString(index)
	case LUA_TLIGHTUSERDATA:
		return L.ToPointer(index)
	case LUA_TUSERDATA:
		return L.ToUserdata(index)
	//case LUA_TTABLE:
	case LUA_TFUNCTION:
		return RefFunction(L.RegistryRefIndex(index))
	//case LUA_TTHREAD:
	case LUA_TCDATA:
		return L.ToCtype(index)
	default:
		log.Fatalf("(*State).ToAny() error : unhandled lua_type %d %s", typ, typ)
	}
	return nil
}
