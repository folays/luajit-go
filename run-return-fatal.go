package luajit

import (
	"log"
)

func (L *State) RunReturnStringFatal(code string) (v any) {
	var (
		err error
	)
	if v, err = L.RunReturnString(code); err != nil {
		log.Fatalf("(*State).RunReturnString(%s) error : %s", code, err)
	}
	return
}

func (L *State) RunReturnBufferFatal(code []byte, chunkName string) (v any) {
	var (
		err error
	)
	if v, err = L.RunReturnBuffer(code, chunkName); err != nil {
		log.Fatalf("(*State).RunReturnBufferFatal(,%s) error : %s", code, err)
	}
	return
}

func (L *State) RunReturnCodeFatal(code string, args ...interface{}) (v any) {
	var (
		err error
	)
	if v, err = L.RunReturnCode(code, args...); err != nil {
		log.Fatalf("(*State).RunReturnCodeFatal(%s) error : %s", code, err)
	}
	return
}

func (L *State) RunReturnFuncFatal(funcName string, args ...interface{}) (v any) {
	var (
		err error
	)
	if v, err = L.RunReturnFunc(funcName, args...); err != nil {
		log.Fatalf("(*State).RunReturnFuncFatal(%s) error : %s", funcName, err)
	}
	return
}
