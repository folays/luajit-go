package b

import (
	"embed"
	_ "embed"
	"github.com/folays/luajit-go"
)

//go:embed b.lua
var fs embed.FS

func Test(L *luajit.State) {
	L.RunEmbedFsFatal(fs)
}
