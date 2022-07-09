package luajit

// GetTable pushes `t[k]`, where `t` is at given index, `k` at top.
// As in Lua, pops the top the stack (the key).
func (L *State) GetTable(index Index) {
	L.Y_lua_gettable(index)
}

// GetField pushes `t[k]`
// Does not modify the stack (beyond pushing the result)
func (L *State) GetField(index Index, k any) {
	index = L.abs_index(index)
	L.PushAny(k)
	L.Y_lua_gettable(index)
}

// SetTable does `t[k] = v`
// As in Lua, pops both the key and the value at the top of the stack.
func (L *State) SetTable(index Index) {
	L.Y_lua_settable(index)
}

// SetTableAny does `t[k] = v` where `t` is at given index, `k` at top.
// Does not exist in Lua.
// As it would in Lua, pops the top of the stack (`k`).
func (L *State) SetTableAny(index Index, v any) {
	index = L.abs_index(index)
	L.PushOne(v) // pushes value. caller did already push the key.
	L.Y_lua_settable(index)
}

// SetField does `t[k] = v` where `t` is at given index, `v` at top.
// As in Lua, pops the top of the stack (`v`).
func (L *State) SetField(index Index, k any) {
	index = L.abs_index(index)
	L.PushOne(k)
	L.insert(-2) // move key before value
	L.Y_lua_settable(index)
}

// SetFieldAny does `t[k] = v` where `t` is at given index.
// Does not modify the stack.
func (L *State) SetFieldAny(index Index, k any, v any) {
	index = L.abs_index(index)
	L.PushOne(k)
	L.PushOne(v)
	L.Y_lua_settable(index)
}
