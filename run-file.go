package luajit

import (
	"log"
)

func (L *State) RunFile(path string) (err error) {
	if err = L.LoadFile(path); err != nil {
		log.Printf("(*State).RunFile(%s) ERROR : %s", path, err)
		return
	}

	if _, err = L._run(0, 0); err != nil {
		log.Printf("(*State).RunFile(%s) ERROR : %s", path, err)
		return
	}

	return
}
