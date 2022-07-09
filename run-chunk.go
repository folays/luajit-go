package luajit

import (
	"crypto/sha256"
	"unsafe"
)

func (L *State) _chunkCodeInfo(code string) (codeByte []byte, info *chunkInfo) {
	codeByte = unsafe.Slice(unsafe.StringData(code), len(code))

	var (
		sum256    = sha256.Sum256(codeByte)
		chunkName = unsafe.String(&sum256[0], len(sum256))

		pkg, dir = L.GetCallerPkgDir()
	)

	info = &chunkInfo{
		chunkPkg:  pkg,
		chunkDir:  dir,
		chunkFile: chunkName,
	}

	return
}

func (L *State) RunChunkString(code string) (err error) {
	var (
		codeByte, info = L._chunkCodeInfo(code)
	)

	if err = L.LoadChunkBuffer(codeByte, info); err != nil {
		//log.Printf("(*State).RunChunkString(%s) error : %s", code, err)
		return
	}

	if _, err = L._run(0, 0); err != nil {
		//log.Printf("(*State).RunChunkString(%s) ERROR : %s", code, err)
		return
	}

	return
}

func (L *State) RunChunkCode(code string, args ...interface{}) (err error) {
	var (
		codeByte, info = L._chunkCodeInfo(code)
	)

	if err = L.LoadChunkBuffer(codeByte, info); err != nil {
		//log.Printf("(*State).RunChunkCode(%s) error : %s", code, err)
		return
	}

	L.PushMultiple(args)
	if _, err = L._run(len(args), 0); err != nil {
		//log.Printf("(*State).RunChunkCode(%s) ERROR : %s", code, err)
		return
	}

	return
}

func (L *State) RunChunkBuffer(code []byte, chunkName string) (err error) {
	var (
		pkg, dir = L.GetCallerPkgDir()
		info     = &chunkInfo{
			chunkPkg:  pkg,
			chunkDir:  dir,
			chunkFile: chunkName,
		}
	)

	if err = L.LoadChunkBuffer(code, info); err != nil {
		//log.Printf("(*State).RunChunkBuffer(,%s) error : %s", chunkName, err)
		return
	}

	if _, err = L._run(0, 0); err != nil {
		//log.Printf("(*State).RunChunkBuffer(,%s) ERROR : %s", chunkName, err)
		return
	}

	return
}
