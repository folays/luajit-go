package luajit

// Module_preload : does package.preload[moduleName] = fnLuaopen
func (L *State) Module_preload(moduleName string, fnLuaopen LuaCFunction) (err error) {
	//err = L.RunCode("package.preload[({...})[1]] = ({...})[2]", moduleName, fnLuaopen)

	L.TableGetField("package", "preload")
	defer L.Pop(1)

	L.SetFieldAny(-1, moduleName, fnLuaopen)

	return
}
