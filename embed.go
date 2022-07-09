package luajit

import (
	"embed"
	"errors"
	"io/fs"
	"log"
)

type embedFS embed.FS

func (L *State) embed_prepare_bridge() {
	L.BridgeAdd(embedFS{})
	L.BridgeAddShared(embedFS{}, embed.FS{})
}

func (eFS embedFS) _lua_payload_new(L *State) []any {
	var (
		pkg, dir = L.GetCallerPkgDir()
	)

	return []any{pkg, dir}
}

// Lua_readdir returns either :
// - table, nil
// - nil,   nil   : directory does not exist
// - nil,   error : any error other that ErrNotExist
func (eFS embedFS) Lua_readdir(L *State, path string) (err error) {
	var (
		entries []fs.DirEntry
	)

	if entries, err = embed.FS(eFS).ReadDir(path); err != nil {
		log.Printf("luajit.(embedFS).Lua_readdir(,%s) error : %s", path, err)
		if errors.Is(err, fs.ErrNotExist) {
			err = nil
		}
		L.pushNil()
		return
	}

	L.NewTable()
	for i, entry := range entries {
		L.SetFieldAny(-1, i+1, entry.Name())
	}

	return
}

// Lua_readfile returns either :
// - string,   nil
// - nil,      nil   : file does not exist
// - nil,      error : any error other that ErrNotExist
func (eFS embedFS) Lua_readfile(L *State, path string) (data []byte, err error) {
	if data, err = embed.FS(eFS).ReadFile(path); err != nil {
		log.Printf("luajit.(embedFS).Lua_readfile(,%s) error : %s", path, err)
		if errors.Is(err, fs.ErrNotExist) {
			err = nil
		}
		return
	}

	return
}

// Lua_loadfile returns either :
// - function, nil
// - nil,      nil   : file does not exist
// - nil,      error : any error other that ErrNotExist
func (eFS embedFS) Lua_loadfile(L *State, path string) (err error) {
	payload := L.ToPayload(1)

	var (
		data []byte

		info = &chunkInfo{
			chunkPkg:  payload[0].(string),
			chunkDir:  payload[1].(string),
			chunkFile: path,
		}
	)

	if data, err = embed.FS(eFS).ReadFile(path); err != nil {
		//log.Printf("luajit.(embedFS).Lua_loadfile(,%s) error : %s", path, err)
		if errors.Is(err, fs.ErrNotExist) {
			err = nil
		}
		L.pushNil()
		return
	}

	if err = L.LoadChunkBuffer(data, info); err != nil {
		log.Printf("luajit.(embedFS).Lua_loadfile(,%s) error : %s", path, err)
		L.pushNil()
		return
	}

	return
}
