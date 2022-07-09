package luajit

import (
	"reflect"
)

// MethodsAddFilter_tableIndex extends table _G[tableIndex] with all `o` methods :
//   - condition : only consider functions having filter() returning non-empty string
//   - exported name : use filter() result
func (L *State) MethodsAddFilter_tableIndex(tableIndex Index, oValue reflect.Value, filter fnFuncsAddFilter) {
	for i := 0; i < oValue.NumMethod(); i++ {
		var (
			method_func = oValue.Method(i).Interface()
			method_name = oValue.Type().Method(i).Name
		)

		if method_name = filter(method_name); method_name == "" {
			continue
		}

		//fmt.Printf("MethodsAddFilter_tableIndex [ %d] : %s\n", i, method_name)
		L.SetFieldAny(tableIndex, method_name, method_func)
	}
}

// MethodsAddFilter_tableName extends table _G[tableName] with all `o` methods :
//   - condition : only consider functions having filter() returning non-empty string
//   - exported name : use filter() result
func (L *State) MethodsAddFilter_tableName(tableName string, o any, filter fnFuncsAddFilter) {
	//fmt.Printf("MethodsAddFilter_tableName _G[%s] type %T\n", tableName, o)

	L.TableGet(tableName)
	defer L.Pop(1)

	L.MethodsAddFilter_tableIndex(-1, reflect.ValueOf(o), filter)
}

// MethodsAdd_prefix extends table _G[tableName] with all `o` methods named `prefix`* :
//   - exported name : will have "Lua_" prefix stripped
func (L *State) MethodsAdd_prefix(tableName string, o any, prefix string) {
	L.MethodsAddFilter_tableName(tableName, o, FilterFuncsPrefix(prefix))
}

// MethodsAdd extends table _G[tableName] with all `o` methods named "Lua_*" :
//   - exported name : will have "Lua_" prefix stripped
func (L *State) MethodsAdd(tableName string, o any) {
	L.MethodsAdd_prefix(tableName, o, "Lua_")
}
