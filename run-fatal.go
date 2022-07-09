package luajit

import (
	"log"
)

func (L *State) RunStringFatal(code string) {
	if err := L.RunString(code); err != nil {
		log.Fatalf("(*State).RunStringFatal(%s) error : %s", code, err)
	}
}

func (L *State) RunCodeFatal(code string, args ...interface{}) {
	if err := L.RunCode(code, args...); err != nil {
		log.Fatalf("(*State).RunCodeFatal(%s) error : %s", code, err)
	}
}

func (L *State) RunBufferFatal(code []byte, chunkName string) {
	if err := L.RunBuffer(code, chunkName); err != nil {
		log.Fatalf("(*State).RunBufferFatal(,%s) error : %s", chunkName, err)
	}
}

func (L *State) RunFuncFatal(funcName string, args ...any) {
	if err := L.RunFunc(funcName, args...); err != nil {
		log.Fatalf("(*State).RunFuncFatal(%s) error : %s", funcName, err)
	}
}

// RunFuncWithResultsFatal caller should call RunClearResults() after, if nbResults > 0
func (L *State) RunFuncWithResultsFatal(nbResults int, funcName string, args ...any) (retResults int) {
	var (
		err error
	)
	if retResults, err = L.RunFuncWithResults(nbResults, funcName, args...); err != nil {
		log.Fatalf("(*State).RunFuncWithResultsFatal(,%s) error : %s", funcName, err)
	}
	return
}

func (L *State) RunFileFatal(path string) {
	if err := L.RunFile(path); err != nil {
		log.Fatalf("(*State).RunFileFatal(%s) error : %s", path, err)
	}
}

func (L *State) RunChunkStringFatal(code string) {
	if err := L.RunChunkString(code); err != nil {
		log.Fatalf("(*State).RunChunkStringFatal(%s) error : %s", code, err)
	}
}

func (L *State) RunChunkCodeFatal(code string, args ...any) {
	if err := L.RunChunkCode(code, args...); err != nil {
		log.Fatalf("(*State).RunChunkCodeFatal(%s) error : %s", code, err)
	}
}

func (L *State) RunChunkBufferFatal(data []byte, chunkName string) {
	if err := L.RunChunkBuffer(data, chunkName); err != nil {
		log.Fatalf("(*State).RunChunkBufferFatal(%s) error : %s", chunkName, err)
	}
}
