package main

import (
	"github.com/folays/luajit-go"
)

func main() {
	Lmonitor := luajit.NewState()
	Lmonitor.ReasonableDefaults()
	Lmonitor.RunCodeFatal("luajit_mcode.diagnose(({...})[1])", 0)

	for j := 1; j <= 1000; j++ {
		L := luajit.NewState()
		L.ReasonableDefaults()

		L.RunStringFatal(`ffi.C.malloc(256 * 4096 * 3)`) // 3 MB
		L.RunStringFatal(`ffi.C.malloc(210 * 4096)`)     //
		L.RunStringFatal(`ffi.C.malloc(4 * 7)`)          //
		L.RunStringFatal(`ffi.C.malloc(4 * 256 * 4)`)    //

		L.RunStringFatal(`ffi.C.malloc(1024 * 1024)`)

		L.RunStringFatal(`local n = 1 ; for i=1,10000 do n = ((n + 1) * 7) + 1 * 31 end ; print(n)`)

		if j%200 == 0 {
			Lmonitor.RunCodeFatal("luajit_mcode.diagnose(({...})[1])", j)
		}
	}
}
