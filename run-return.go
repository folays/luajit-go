package luajit

import (
	"log"
)

// RunReturnString RunReturn* family : very niche functions.
//
// Example: How to pass a value from Go, and retrieve it after being processed by some Lua code:
//
//	value, err := L.RunReturnCode(`return do_something_with(({...})[1])`, someUnprocessedValue)
func (L *State) RunReturnString(code string) (v any, err error) {
	return L.RunReturnCode(code)
}

func (L *State) RunReturnBuffer(code []byte, chunkName string) (v any, err error) {
	if err = L.LoadBuffer(code, chunkName); err != nil {
		log.Printf("(*State).RunReturnBuffer(,%s) ERROR : %s", chunkName, err)
		return
	}
	_, err = L._run(0, 1)
	if err != nil {
		log.Printf("(*State).RunReturnBuffer(,%s) ERROR : %s", chunkName, err)
		return
	}
	defer L.RunClearResults()
	return L.ToAny(-1), err
}

func (L *State) RunReturnCode(code string, args ...interface{}) (v any, err error) {
	if err = L.LoadString(code); err != nil {
		log.Printf("(*State).RunReturnCode(%s) ERROR : %s", code, err)
		return
	}
	L.PushMultiple(args)
	if _, err = L._run(len(args), 1); err != nil {
		log.Printf("(*State).RunReturnCode(%s) ERROR : %s", code, err)
		return
	}
	defer L.RunClearResults()
	return L.ToAny(-1), err
}

func (L *State) RunReturnFunc(funcName string, args ...interface{}) (v any, err error) {
	if _, err = L._runFuncWithResults(1, funcName, args); err != nil {
		log.Printf("(*State).RunReturnFunc(%s) ERROR : %s", funcName, err)
		return
	}
	defer L.RunClearResults()
	return L.ToAny(-1), err
}
