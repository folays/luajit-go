package luajit

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"strings"
)

func (L *State) RunEmbedFsFatal(eFS embed.FS) {
	L.RunEmbedFsPathFatal(eFS, ".")
}

func (L *State) RunEmbedFsPathFatal(eFS embed.FS, path string) {
	var (
		err error
	)

	{
		var (
			files []fs.DirEntry
		)
		if files, err = eFS.ReadDir(path); err != nil {
			log.Fatalln(err)
		}
		for _, file := range files {
			var (
				name = file.Name()
			)

			if strings.HasSuffix(name, ".lua") == false {
				continue
			}
			if path != "." {
				name = fmt.Sprintf("%s/%s", path, name)
			}

			L.RunEmbedFsFileFatal(eFS, name)
		}
	}
}

func (L *State) RunEmbedFsFileFatal(eFS embed.FS, path string) {
	if err := L.RunEmbedFsFile(eFS, path); err != nil {
		log.Fatalf("(*State).RunEmbedFsFileFatal(,%s) error : %s", path, err)
	}
}

func (L *State) RunEmbedFsFile(eFS embed.FS, path string) (err error) {
	if err = L.LoadEmbedFsFile(eFS, path); err != nil {
		//log.Printf("(*State).RunEmbedFsFile(%s) error : %s", path, err)
		return
	}
	_, err = L._run(0, 0)
	if err != nil {
		//log.Printf("(*State).RunEmbedFsFile(%s) error : %s", path, err)
		return
	}
	return
}

// LoadEmbedFsFile is Lua's loadfile() alike.
//
// It loads (embed.FS).ReadFile(), and push the chunk at the top of the stack.
//
// It also try to load & save a cached JIT-ted version from /tmp/
func (L *State) LoadEmbedFsFile(eFS embed.FS, path string) (err error) {
	var (
		data []byte

		pkg, dir = L.GetCallerPkgDir()
		info     = &chunkInfo{
			chunkPkg:  pkg,
			chunkDir:  dir,
			chunkFile: path,
		}
	)

	if data, err = eFS.ReadFile(path); err != nil {
		log.Fatalf("(*State).LoadEmbedFsFile(,%s) error : %s", path, err)
	}

	return L.LoadChunkBuffer(data, info)
}
