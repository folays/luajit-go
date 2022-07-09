package main

import (
	"embed"
	"fmt"
	"github.com/folays/luajit-go"
	"math"
	"os"
	"reflect"
)

//go:embed *.lua
var fs embed.FS

func main() {
	var (
		L *luajit.State
	)

	L = luajit.NewState()
	L.ReasonableDefaults()

	L.RunEmbedFsFatal(fs)

	var (
		datas = []any{
			int(math.MaxInt64),
			int64(math.MaxInt64),
			int(math.MinInt64),
			int64(math.MinInt64),

			uint(math.MaxUint64),
			uint64(math.MaxUint64),
		}
	)

	for _, data := range datas {
		L.SetGlobalAny("a", data)
		L.RunFuncFatal("check")

		ret := L.GetGlobalAny("a")
		ret2 := reflect.ValueOf(ret).Convert(reflect.TypeOf(data)).Interface()

		if ret2 != data || fmt.Sprint(ret) != fmt.Sprint(data) {
			fmt.Printf("\033[31mNOT OK\033[39m\n")
			os.Exit(1)
		}
	}
}
