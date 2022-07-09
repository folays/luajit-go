package main

import (
	"fmt"
	"github.com/folays/luajit-go"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	{
		fmt.Printf("PROFILE START\n")
		start := time.Now()
		defer func() {
			fmt.Printf("LUAJIT_PROFILE : end after %v\n", time.Since(start))
		}()
	}
	{
		file, _ := os.Create("/tmp/luajit.pprof")
		defer file.Close()
		pprof.StartCPUProfile(file)
		defer pprof.StopCPUProfile()
	}

	//if true {
	//	for i := 0; i < 100000*7; i++ {
	//		os.Stat("/tmp/luajit_bytecode/main_defaults.ljbc")
	//	}
	//	return
	//}

	for i := 0; i < 100000; i++ {
		var (
			L *luajit.State
		)

		L = luajit.NewState()
		L.ReasonableDefaults()
		//L.Close()
	}
}
