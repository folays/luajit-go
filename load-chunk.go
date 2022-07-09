package luajit

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type chunkInfo struct {
	chunkPkg  string
	chunkDir  string
	chunkFile string
}

const bytecode_path = "/tmp/luajit_bytecode"

var flag_bytecode = false

func init() {
	if _, flag_bytecode = os.LookupEnv("BYTECODE"); flag_bytecode {
		os.Mkdir(bytecode_path, 0755)
	}
}

func (info *chunkInfo) signature() string {
	return fmt.Sprintf("@%s:%s:%s", info.chunkPkg, info.chunkDir, info.chunkFile)
}

func (info *chunkInfo) signBytes() []byte {
	return []byte(info.signature())
}

func (info *chunkInfo) pathObj() (pathObj string) {
	pathObj = info.chunkPkg + "/" + info.chunkFile
	pathObj = strings.ReplaceAll(pathObj, "/", "_")
	pathObj = strings.TrimSuffix(pathObj, ".lua")
	pathObj = bytecode_path + "/" + pathObj + ".ljbc"
	return
}

func (L *State) LoadChunkFile(pathSrc string, info *chunkInfo) (err error) {
	if flag_bytecode == false { // skip bytecode
		if data, err := os.ReadFile(pathSrc); err != nil {
			return err
		} else {
			return L.LoadBuffer(data, info.signature())
		}
	}

	var (
		statSrc, _ = os.Stat(pathSrc)

		pathObj    = info.pathObj()
		statObj, _ = os.Stat(pathObj)

		data []byte
	)

	if statObj != nil && statObj.ModTime().After(statSrc.ModTime()) {
		return L._loadFile(pathObj)
	}

	if data, err = os.ReadFile(pathSrc); err != nil {
		return
	}

	if err = L.LoadBuffer(data, info.signature()); err != nil {
		log.Printf("(*State).LoadChunkFile(%s) error : %s", pathSrc, err)
	}

	if errDump := L.DumpFileBytecode(pathObj); errDump != nil {
		return
	}

	return
}

func (L *State) LoadChunkBuffer(code []byte, info *chunkInfo) (err error) {
	if flag_bytecode == false { // skip bytecode
		return L.LoadBuffer(code, info.signature())
	}

	var (
		sumData = bytes.Join([][]byte{code, info.signBytes()}, []byte(""))
		sum256  = sha256.Sum256(sumData)
		sum29   = binary.LittleEndian.Uint32(sum256[:]) >> 3

		pathObj    = info.pathObj()
		statObj, _ = os.Stat(pathObj)
	)

	if statObj != nil && uint32(statObj.ModTime().Nanosecond()) == sum29 {
		return L._loadFile(pathObj)
	}

	if err = L.LoadBuffer(code, info.signature()); err != nil {
		log.Printf("(*State).LoadChunkBuffer(,%s) error : %s", info.signature(), err)
		return
	}

	if errDump := L.DumpFileBytecode(pathObj); errDump != nil {
		return
	} else {
		mtime := time.Unix(time.Now().Unix()-3*86400, int64(sum29))
		os.Chtimes(pathObj, mtime, mtime)
	}

	return
}
