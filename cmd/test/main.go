package main

import (
	"github.com/folays/luajit-go"
	"github.com/folays/luajit-go/cmd/test/a"
	"github.com/folays/luajit-go/cmd/test/b"
)

func main() {
	var (
		L *luajit.State
	)

	L = luajit.NewState()
	L.ReasonableDefaults()

	a.Push(L)
	b.Test(L)

	L.SetGlobalAny("fn", func() {})
	L.RunString("fn = nil ; collectgarbage()")

	L.RunString("print(debug.traceback())")
}
