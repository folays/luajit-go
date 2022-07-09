package luajit

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"unsafe"
)

/*
#include <stdlib.h>
#include <lua.h>
#include <lauxlib.h>

#include "load-file.h"
*/
import "C"

type LuaWriter C.lua_Writer

type luaFileDump struct {
	buf bytes.Buffer
}

func (L *State) _loadFile(path string) (err error) {
	if L.Y_luaL_loadfile(path) != 0 {
		return fmt.Errorf(L.errorStringPop())
	}

	return nil
}

func (L *State) LoadFile(path string) (err error) {
	var (
		pkg, dir = L.GetCallerPkgDir()
		info     = &chunkInfo{
			chunkPkg:  pkg,
			chunkDir:  dir,
			chunkFile: path,
		}
	)

	if err = L.LoadChunkFile(path, info); err != nil {
		log.Printf("(*State).LoadFile(%s) error : %s", path, err)
		return
	}

	return
}

// DumpFileBytecode dump the bytecode of the function at the top of the stack, to pathObj
func (L *State) DumpFileBytecode(pathObj string) (err error) {
	dump := &luaFileDump{}

	C.lua_dump(L.state, LuaWriter(C.file_dump_writer_c), unsafe.Pointer(dump))

	{
		var (
			file *os.File
		)

		if file, err = os.Create(pathObj); err != nil {
			log.Fatalf("(*State).DumpFileBytecode(%s) error : %s", pathObj, err)
		}
		file.ReadFrom(&dump.buf)
		file.Close()
	}

	return
}

//export file_dump_writer_go
func file_dump_writer_go(l LuaStatePtr, p unsafe.Pointer, sz C.size_t, ud unsafe.Pointer) int {
	var (
		dump = (*luaFileDump)(ud)
		buf  = C.GoBytes(p, C.int(sz))
	)

	dump.buf.Write(buf)

	return 0
}
