package luajit

import (
	"reflect"
	"strings"
)

type fnFuncsAddFilter func(name string) string

func FilterFuncsPrefix(prefix string) fnFuncsAddFilter {
	return func(name string) string {
		if strings.HasPrefix(name, prefix) {
			return strings.TrimPrefix(name, prefix)
		}
		return ""
	}
}

// FuncsAddFilter_tableIndex extends table _G[tableIndex] with all `o` functions :
//   - WITHOUT the receiver (which later the Lua callers will have to provide)
//   - condition : only consider functions having filter() returning non-empty string
//   - exported name : use filter() result
func (L *State) FuncsAddFilter_tableIndex(tableIndex Index, oType reflect.Type, filter fnFuncsAddFilter) {
	for i := 0; i < oType.NumMethod(); i++ {
		var (
			method      = oType.Method(i)
			method_func = method.Func.Interface()
			method_name = method.Name
		)

		if method_name = filter(method_name); method_name == "" {
			continue
		}

		//fmt.Printf("FuncsAddFilter_tableIndex [ %d] : %s\n", i, method_name)
		L.SetFieldAny(tableIndex, method_name, method_func)
	}
}

// FuncsAddFilter_tableName extends table _G[tableName] with all `o` functions :
//   - WITHOUT the receiver (which later the Lua callers will have to provide)
//   - condition : only consider functions having filter() returning non-empty string
//   - exported name : use filter() result
func (L *State) FuncsAddFilter_tableName(tableName string, o any, filter fnFuncsAddFilter) {
	//fmt.Printf("FuncsAddFilter_tableName _G[%s] type %T\n", tableName, o)

	L.TableGet(tableName)
	defer L.Pop(1)

	L.FuncsAddFilter_tableIndex(-1, reflect.TypeOf(o), filter)
}

// FuncsAdd_prefix extends table _G[tableName] with all `o` functions named `prefix`* :
//   - WITHOUT the receiver (which later the Lua callers will have to provide)
//   - exported name : will have "Lua_" prefix stripped
func (L *State) FuncsAdd_prefix(tableName string, o any, prefix string) {
	L.FuncsAddFilter_tableName(tableName, o, FilterFuncsPrefix(prefix))
}

// FuncsAdd extends table _G[tableName] with all `o` functions named "Lua_*" :
//   - WITHOUT the receiver (which later the Lua callers will have to provide)
//   - exported name : will have "Lua_" prefix stripped
func (L *State) FuncsAdd(tableName string, o any) {
	L.FuncsAdd_prefix(tableName, o, "Lua_")
}

// FuncAdd sets _G[tableName][funcName] = `fn` (a function)
// If you want to set _G[funcName], use L.SetGlobalAny() instead.
func (L *State) FuncAdd(tableName string, funcName string, fn any) {
	L.TableGet(tableName)
	defer L.Pop(1)

	L.SetFieldAny(-1, funcName, fn)
}
