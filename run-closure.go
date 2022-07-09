package luajit

import (
	"fmt"
	"log"
)

func (L *State) RunClosure(fn LuaCFunction, args ...interface{}) (err error) {
	_, err = L._runClosureWithResults(0, fn, args)
	return

}

func (L *State) _runClosureWithResults(nbResults int, fn LuaCFunction, args []any) (retResults int, err error) {
	L.pushcclosure(fn, 0)
	if L.Type(-1) != LUA_TFUNCTION {
		log.Fatalf("(*State)._runClosureWithResults() : not a function")
	}
	nbArgs := L.PushMultiple(args)
	retResults, err = L._run(nbArgs, nbResults)
	if err != nil {
		fmt.Printf("\033[38;2;128;128;128mERROR: %s\033[39;22m\n", err)
		return
	}
	return
}
