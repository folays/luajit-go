package luajit

import (
	"log"
)

// beware : ToString() could modify value, which can confuses next call to lua_next().
func (L *State) _table_iterate(tableIndex Index, fn func()) {
	tableIndex = L.abs_index(tableIndex)
	L.pushNil()
	for L.Y_lua_next(tableIndex) != 0 {
		fn()
		L.Pop(1) // remove value
	}
	return
}

// TableGet pushes (or create) _G[tableName] on the stack, which must be a table.
func (L *State) TableGet(tableName any) {
	L.GetGlobal(tableName)
	switch L.Type(-1) {
	case LUA_TTABLE:
		return
	case LUA_TNIL:
		L.Pop(1)
		L.NewTable()
		L.pushvalue(-1)        // keep a backup to return it
		L.SetGlobal(tableName) // will Pop the table
		return
	default:
		log.Fatalf("TableGet() : already not a table")
	}
}

func (L *State) _table_get_fatal(tableName any) {
	L.GetGlobal(tableName) // table
	if L.IsTable(-1) == false {
		log.Fatalf("_table_get_fatal(%s) : not a table", tableName)
	}
}

// TableGetField pushes `_G[tableName][k]` onto the stack.
func (L *State) TableGetField(tableName any, k any) {
	L._table_get_fatal(tableName) // table
	L.GetField(-1, k)             // table,value
	L.remove(-2)                  // value
}

// TableSetAny does `_G[tableName][k] = v` where `k` is at top.
// Does not exist in Lua.
// As it would in Lua, pops the top of the stack (`k`).
func (L *State) TableSetAny(tableName, v any) {
	L._table_get_fatal(tableName) // table
	L.insert(-2)                  // table,key
	L.SetTableAny(-2, v)          // table
	L.Pop(1)                      // <empty>
}

// TableSetField does `_G[tableName][k] = v` where `v` is at top.
// As in Lua, pops the top of the stack (`v`).
func (L *State) TableSetField(tableName, k any) {
	L._table_get_fatal(tableName) // table
	L.insert(-2)                  // table,value
	L.SetField(-2, k)             // table
	L.Pop(1)                      // <empty>
}

// TableSetFieldAny does `_G[tableName][k] = v`
// Does not modify the stack.
func (L *State) TableSetFieldAny(tableName, k any, v any) {
	L._table_get_fatal(tableName) // table
	L.SetFieldAny(-1, k, v)       // table
	L.Pop(1)                      // <empty>
}

func (L *State) TableGetKeysAny() (ret []any) {
	L._table_iterate(-1, func() {
		ret = append(ret, L.ToAny(-1))
	})

	return
}

func (L *State) TableGetKeysString() (ret []string) {
	L._table_iterate(-1, func() {
		ret = append(ret, L.ToStringSafe(-1))
	})

	return
}
