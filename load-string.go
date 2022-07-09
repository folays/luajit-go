package luajit

import (
	"fmt"
)

// LoadString loads code into lua stack
func (L *State) LoadString(str string) error {
	if L.Y_luaL_loadstring(str) != 0 {
		return fmt.Errorf(L.errorStringPop())
	}

	return nil
}

func (L *State) LoadBuffer(buf []byte, chunkName string) error {
	if L.Y_luaL_loadbuffer(buf, chunkName) != 0 {
		return fmt.Errorf(L.errorStringPop())
	}

	return nil
}
