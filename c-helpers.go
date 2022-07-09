package luajit

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"unsafe"
)

/*
#include <stdlib.h>
#include <lua.h>
*/
import "C"

func (L *State) Load_C_Helpers() {
	L.FuncAdd("luajit", "exists", L.helper_file_exists)
	L.FuncAdd("luajit", "readfile", L.helper_readfile)
	L.FuncAdd("luajit", "loadfile", L.helper_loadfile)
	L.FuncAdd("luajit", "stacktrace", L._helper_stacktrace)
	L.FuncAdd("luajit", "set_error_handler", L._helper_set_error_handler)
}

func (L *State) _helper_stacktrace() []byte {
	return debug.Stack()
}

func (L *State) _helper_set_error_handler(ref RefFunction) {
	L.error_ref = ref
}

func (L *State) _helper_path_context(level int) (pkg string, pkgDir string, pkgFile string, err error) {
	var (
		ar          C.lua_Debug
		what        = (*C.char)(unsafe.Pointer(unsafe.SliceData([]byte("S\x00"))))
		pathContext string
		fields      []string
	)

	if C.lua_getstack(L.state, C.int(level), &ar) != 1 {
		return "", "", "", fmt.Errorf("lua_getstack(L, %d) error", level)
	}

	if C.lua_getinfo(L.state, what, &ar) == 0 {
		return "", "", "", fmt.Errorf("lua_getinfo(L, S, &ar) error")
	}

	pathContext = C.GoString(ar.source)

	if strings.HasPrefix(pathContext, "@") == false {
		return "", "", "", fmt.Errorf("(*State)._helper_path_context() error : unhandled pathContext %s", pathContext)
	}

	if fields = strings.SplitN(pathContext, ":", 3); len(fields) < 3 {
		return "", "", "", fmt.Errorf("(*State)._helper_path_context() error : incorrect pathContext %s", pathContext)
	}

	return fields[0][1:], fields[1], fields[2], nil
}

func (L *State) helper_file_exists(path string, level int) bool {
	var (
		_, pkgDir, pkgFile, err = L._helper_path_context(level)
	)

	if err != nil {
		log.Fatalf("(*State).helper_file_exists(%s) error : %s", path, err)
		return false
	}

	newPath := filepath.Join(pkgDir, pkgFile, "..", path)

	if _, err = os.Stat(newPath); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	}

	log.Fatalf("(*State).helper_file_exists(%s) error : %s", newPath, err)
	return false
}

// helper_readfile returns either :
// - string,   nil
// - nil,      nil   : file does not exist
// - nil,      error : any error other that ErrNotExist
func (L *State) helper_readfile(path string, level int) (data []byte, err error) {
	var (
		pkgDir, PkgFile string
	)
	_, pkgDir, PkgFile, err = L._helper_path_context(level)

	if err != nil {
		log.Fatalf("(*State).helper_readfile(%s) error : %s", path, err)
		return
	}

	newPath := filepath.Join(pkgDir, PkgFile, "..", path)

	if data, err = os.ReadFile(newPath); err != nil {
		log.Printf("(*State).helper_readfile(%s) error : %s", newPath, err)
		if errors.Is(err, fs.ErrNotExist) {
			err = nil
		}
		return
	}

	return
}

// helper_loadfile returns either :
// - function,   nil
// - nil,        nil   : file does not exist
// - nil,        error : any error other that ErrNotExist
func (L *State) helper_loadfile(path string, level int) (err error) {
	var (
		pkg, pkgDir, PkgFile string
		data                 []byte
	)
	pkg, pkgDir, PkgFile, err = L._helper_path_context(level)

	if err != nil {
		log.Fatalf("(*State).helper_loadfile(%s) error : %s", path, err)
		L.pushNil()
		return
	}

	newPath := filepath.Join(pkgDir, PkgFile, "..", path)

	if data, err = os.ReadFile(newPath); err != nil {
		//log.Printf("(*State).helper_loadfile(%s) error : %s", newPath, err)
		if errors.Is(err, fs.ErrNotExist) {
			err = nil
		}
		L.pushNil()
		return
	}

	var (
		info = &chunkInfo{
			chunkPkg:  pkg,
			chunkDir:  pkgDir,
			chunkFile: filepath.Join(PkgFile, "..", path),
		}
	)

	if err = L.LoadChunkBuffer(data, info); err != nil {
		log.Printf("(*State).helper_loadfile(%s) error : %s", path, err)
		L.pushNil()
		return
	}

	return
}
