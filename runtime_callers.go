package luajit

import (
	"path"
	"runtime"
	"strings"
)

func GetCallerSourceDir(nCallers int) string {
	_, filename, _, _ := runtime.Caller(nCallers + 1)
	return path.Dir(filename)
}

func (L *State) GetCallerPkgDir() (pkg string, pkgDir string) {
	var (
		pcs = make([]uintptr, 30)
	)
	pcs = pcs[:runtime.Callers(1, pcs)]

	frames := runtime.CallersFrames(pcs)

	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if strings.HasPrefix(frame.Function, "reflect.") {
			continue
		}
		if strings.HasPrefix(frame.Function, L.myself_pkg+".") == true &&
			strings.HasSuffix(frame.Function, L.myself_defaults) == false {
			continue
		}

		var (
			dir, file = path.Split(frame.Function)

			before, after string
			found         bool
		)

		before, _, _ = strings.Cut(file, ".")
		pkg = dir + before

		if before, after, found = strings.Cut(frame.File, "@"); found == true {
			_, after, _ = strings.Cut(path.Dir(after), "/")
			pkgDir = before + after
		} else {
			pkgDir = path.Dir(frame.File) // local module
		}

		return
	}

	return
}
