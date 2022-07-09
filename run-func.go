package luajit

import (
	"fmt"
	"log"
)

/*
#include <stdlib.h>
#include <stdio.h>
*/
import "C"

func (L *State) run_func_prepare() {
	L.OpenLib("debug", LuaOpen_debug)
	L.TableGetField("debug", "traceback")
	L.error_ref = RefFunction(L.RegistryRef())
}

func (L *State) RunFunc(funcName string, args ...interface{}) (err error) {
	_, err = L._runFuncWithResults(0, funcName, args)
	return
}

// RunFuncWithResults caller should call RunClearResults() after, if nbResults > 0
func (L *State) RunFuncWithResults(nbResults int, funcName string, args ...interface{}) (retResults int, err error) {
	retResults, err = L._runFuncWithResults(nbResults, funcName, args)
	return
}

func (L *State) _runFuncWithResults(nbResults int, funcName string, args []interface{}) (retResults int, err error) {
	L.Y_lua_getfield(GlobalsIndex, funcName)
	if L.Y_lua_isfunction(-1) == false {
		fmt.Printf("\033[38;2;128;128;128mERROR:\033[39;22m %s is not a function\n", funcName)
		err = fmt.Errorf("_runFuncWithResults() %s : not a function", funcName)
		L.Pop(1) // remove function (which is not a function...)
		return
	}
	nbArgs := L.PushMultiple(args)
	retResults, err = L._run(nbArgs, nbResults)
	if err != nil {
		//fmt.Printf("\033[38;2;128;128;128mERROR: %s\033[39;22m\n", err)
		return
	}
	return
}

// Run executes code in stack
func (L *State) _run(nbArgs int, nbResults int) (retResults int, err error) {
	var (
		top                = L.GetTop() - Index(1+nbArgs) // f,[args] : top - (f + nbArgs)
		indexErrFunc Index = 0
	)

	switch {
	case len(L.stackTops) == 0 && top == 0:
	case len(L.stackTops) >= 1 && top == L.stackTops[len(L.stackTops)-1]:
	default:
		// This check could be lowered to a warning, or altogether be removed.
		// It could be legit if we had previous data on the Lua stack,
		//  .e.g. because we got there deep from inside another Lua func() call...
		log.Fatalf("\033[31mERROR\033[39m (*State)._run() : top != 0 (%d) ; previous call missing RunClearResults()\n", top)
	}

	{
		if Ref(L.error_ref) != RefNil {
			L.RegistryGet(Ref(L.error_ref)) // f,[args],traceback
			L.insert(top + 1)               // traceback,f,[args]
			indexErrFunc = top + 1
		}

		if L.Y_lua_pcall(nbArgs, nbResults, indexErrFunc) != 0 {
			C.fflush(C.stdout)
			C.fflush(C.stderr)
			err = fmt.Errorf("%s", L.errorStringPop())
			//fmt.Printf("\033[35mpcall error\033[39m: %s\n", err)
		}

		if Ref(L.error_ref) != RefNil {
			L.Y_lua_remove(top + 1) // remove `traceback'
		}
	}

	//C.fflush(C.stdout)

	retResults = int(L.GetTop() - top)
	if nbResults > 0 {
		L._stackTops_push(top)
	}

	return
}

// RunClearResults reset stack top, to the top before the call.
func (L *State) RunClearResults() {
	top := L._stackTops_pop()
	L.SetTop(top)
}

func (L *State) _stackTops_push(top Index) {
	L.stackTops = append(L.stackTops, top)
}

func (L *State) _stackTops_pop() (top Index) {
	top = L.stackTops[len(L.stackTops)-1]          // fetch last
	L.stackTops = L.stackTops[:len(L.stackTops)-1] // shrink
	return
}
