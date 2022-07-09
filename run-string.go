package luajit

import (
	"log"
)

func (L *State) RunString(code string) (err error) {
	if err = L.LoadString(code); err != nil {
		log.Printf("(*State).RunString(%s) ERROR : %s", code, err)
		return
	}

	if _, err = L._run(0, 0); err != nil {
		log.Printf("(*State).RunString(%s) ERROR : %s", code, err)
		return
	}

	return
}

func (L *State) RunCode(code string, args ...interface{}) (err error) {
	if err = L.LoadString(code); err != nil {
		log.Printf("(*State).RunCode(%s) ERROR : %s", code, err)
		return
	}

	L.PushMultiple(args)

	if _, err = L._run(len(args), 0); err != nil {
		log.Printf("(*State).RunCode(%s) ERROR : %s", code, err)
		return
	}

	return
}

func (L *State) RunBuffer(code []byte, chunkName string) (err error) {
	if err = L.LoadBuffer(code, chunkName); err != nil {
		log.Printf("(*State).RunBuffer(,%s) ERROR : %s", chunkName, err)
		return
	}

	if _, err = L._run(0, 0); err != nil {
		log.Printf("(*State).RunBuffer(,%s) ERROR : %s", chunkName, err)
		return
	}

	return
}
