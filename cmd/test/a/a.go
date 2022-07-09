package a

import (
	"embed"
	"github.com/folays/luajit-go"
)

//go:embed a1.lua
var fs1 embed.FS

//go:embed a2.lua
var fs2 embed.FS

func Push(L *luajit.State) {
	L.RunEmbedFsFatal(fs1)
	L.SetGlobalAny("embed_a2", fs2)
}
