# LuaJIT-Go

Build instructions at the bottom. By that I mean requirements :
- having LuaJIT in `pkg-config`...
- having gcc or whatever so that `Cgo` can compile C code (obvisouly I link with LuaJIT library...) 

Beside that, this module only has to be `import`ed, like a "normal" module.

## How To Use

```
import (
	"github.com/folays/luajit-go"
)

func main() {
    var L *luajit.State = luajit.NewState()
    L.ReasonableDefaults()
    
    do some stuff with L...
}
```

### + Load some `embedded` (`embed.FS`) `.lua` file ?

```
//go:embed *.lua
//go:embed subdir/*.lua
var luaFS embed.FS

L.RunEmbedFsFatal(luaFS)
L.RunEmbedFsPathFatal(luaFS, "subdir")
```

Note : Lots of calls have a non-`Fatal()`-suffixed version, which returns an `error`.

There is also some helpers to load+run some `embedded` Go's `string` or `[]byte`.

### + Set / Get some Lua variables? (data is a `any`)

```
L.SetGlobalAny("myvar", data)

myvar := L.GetGlobal*("myvar")
myvar := L.GetGlobalAny("myvar").(typecast) // because it returns a `any`
myvar := L.GetGlobalBoolean("myvar")        // returns a `bool`
```

### + Run a func or some code `string` ?

```
L.RunFuncFatal("some_function") // run a func named ...

L.RunStringFatal(`print("hello world")`) // eval and run
```

Note : each of those above can take any number of args. (`...[]any`)

### + Run an anonymous func, and pass it some arg(s) ? // 

```
L.RunCodeFatal("some_function(({...})[1])", args ...[]any)
```

Note : Each of the `args` (`...[]any`) is `lua_push*()`'ed, as a suitable Lua type.

Each arg can be retrieved by `({...})[n]` with `n >= 1 && < len(args)`.

Note : Above the ~~girl~~ function has no name, and is no-one. The function is a 1st-class variable, which is then `pcall()`'ed with the `...[]any`.
You could prettify by prepending the code with some `local a, b = select(1, ...)` with some length caveats (if any of the arg is `nil`), or with some `local args = {...}`.

## Helpers

### Basic helpers exported to Lua

A new `luajit.*` Lua module is exposed to Lua, specific to my module. Do not confure with `jit.*` (LuaJIT's one).

- `luajit.exists()` : relative `path black-magic fileops`.
- `luajit.readfile()` : ^ same
- `luajit.loadfile()` : ^ same


- `luajit.stacktrace()` : return the Go stack trace. Used by `0b-tools-debug.lua`. To "concat" stack of Lua+Go.
- `luajit.set_error_handler()` : same as above. I mean, to set the error handler for `pcall()`, including `Run.*()` from Go.

The unspoken `path black-magic fileops` is the fact that those will operate relative to the caller. The magicness comes from the fact that if the Lua caller was defined :
- defined in local filesystem path (on your computer) : it will be relative to that.
- defined in a embed.FS : it will be relative to that. (I think, I don't remember)

It means that users of this module can, during development, have some local files (never pushed to git),
and `luajit.*()` will mostly find them in the relative path you expect it to do.

Of course, for published modules, you should only use Lua methods exposed on `embed.FS` objects.

### Export a Go function to Lua

```
L.FuncAdd("calls", "go_mem_usage", go_mem_usage)

func go_mem_usage() string {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	return fmt.Sprintf("%.3f MB", float64(stats.HeapAlloc)/1e6)
}
```

The black-magicness of the module will take care of passed-in and returned Go arguments.

### Export multiple Go functions to Lua, **EASILY**

```
main() {
    ui.L.MethodsAdd("calls", &calls{})
}

type mycalls struct{}

func (mycall *mycalls) Lua_Test(a int, b int, c float32) (int, string) {
	fmt.Printf("HEY TEST HAS BEEN CALLED ! %v %v %v\n", a, b, c)

	return 42, "quack quack"
}
```

All `mycall`(`type struct`)'s methods prefixed with `Lua_` will be available in Lua as `calls.*`

All return arguments will be returned to Lua... : `local ret, msg = calls.test()`

### Export a Go object (like a `struct` or whatever) **EASILY** as an userdata+methods to Lua

See here how all Go's `embed.FS` are exposed usefully to Lua : https://github.com/folays/luajit-go/blob/master/embed.go

```
// needed to be able to "overload" embed.FS with some additional Lua_*() Go methods
type embedFS embed.FS

func (L *State) embed_prepare_bridge() {
	L.BridgeAdd(embedFS{})                   // register *YOUR* Go type
	L.BridgeAddShared(embedFS{}, embed.FS{}) // alias *YOUR* Go type, with the Go base type
}

[...]

func (eFS embedFS) Lua_loadfile(L *State, path string) (err error) { [...] }
```

Here, registering `embed.FS` would not be so much useful, because you can't add Go methods on it.

The above let you :
- tell Lua that `embed.FS` can be assumed to being equivalent to *YOUR* embedFS nearly-alias type
- add a Lua `:loadfile(path)` to any exposed-to-Lua embed.FS object  

**ATTENTION REQUIRED** : That's one of the **MAIN FACILITATOR** of `LuaJIT-Go`.

Please re-read again : You can directly expose your Go objects, to Lua, as userdata objects,
**WITHOUT ANY** boilerplate code.

You just need to have some `Lua_()` methods on your object, which takes some Go args, and returns some Go args. You don't even have to involve Lua in any form here.

### Pass an embed.FS to Lua

You can pass any Go's `embed.FS` as an argument to Lua, which will give you an useful Lua object.

### Automatically pass some `Go` type to every Go function needing one...

```
func (L *State) FuncEverywherePass(v any, type_in_index_max int) { ...
```

The intention here is that, if you have some `Go type` that nearly 99% of your exposed Go functions will need,
it can be cumbersome to pass it everytime.

`LuaJIT-Go` itself calls `L.FuncEverywherePass(L, 0)` so that any Go function called from Lua
needing a `*luajit.State` in the first (`0`) arg, Lua code can omit it.
And as a consequence, you can also omit to expose it to Lua from Go.

That's only usefull if you can tell that in your whole `Lua`'s "state" (scope = `lua_State *`),
you will **only** have **ONE** instance of this type ever created (somehow like a singleton).

## Included Lua modules

- `json` : https://github.com/rxi/json.lua
- `reflect` : https://github.com/corsix/ffi-reflect / http://corsix.github.io/ffi-reflect/

You can `require()` those from Lua.

## Additional Lua functions

See https://github.com/folays/luajit-go/tree/master/lua-helpers :

- Some `printf` / `snprintf` / `errorf` / `assertf`
- `log` / `logf`

## Behavioral differences

- Lua `debug.traceback()` and `pcall` errors will concat the Go backtrace after the Lua one.

## Propose your own `Go Module` exposing `.lua` files : TODO

I think I already implemented it.

## Propose your own `Go Module` exposing 3rd-party `C` (`Cgo`)

Want to propose your own additional `Go Module` exposing a `C library` as a Lua module ?

```
package yourmodule

import (
	"github.com/folays/luajit-go"
)

var (
	Luaopen_yourmodule_c  = luajit.LuaCFunction(C.luaopen_yourmodule)
	Luaopen_yourmodule_go = _luaopen_yourmodule
)

func _luaopen_yourmodule(L *luajit.State) (err error) {
	if err = L.Module_preload("yourmodule", Luaopen_yourmodule_c); err != nil {
		log.Fatalln(err)
	}

	return
}
```

Here above the C function `luaopen_yourmodule` is expected to behave like a regular C `luaopen_*()` Lua C module.

Any of your users would use your `github.com/yourname/yourmodule` as :

```
package main

import (
	"github.com/folays/luajit-go"
	"github.com/yourname/yourmodule"
)

func main() {
	var L *luajit.State = L = luajit.NewState()
	L.ReasonableDefaults()

	yourmodule.Luaopen_yourmodule_go(L)
	
	// ... now you *COULD* do in Lua : `local yourmodule = require("yourmodule")` and use it.
}
```

So it respects the semantics of how-to open Lua module in C, but the mechanism is exposed in Go.

## Backtraces

Work has been done to ensure that the calling Go routine backtrace appears concat'ed to Lua backtraces.

## Lua coroutines : TODO (or never if useless)

I didn't need them, so I did not think of how they would be useful w.r.t. a Go env already ultra parallelized.

## Sandboxing : TODO

- TODO/maybe : expose less "dangerous" functions to the global env
- TODO/maybe : otherwise maybe expose an "easy API" in Go to run sandboxed code in a non-global less-filled table 

## Level of Quality

I put some love in this module. It should work well, and I did some profiling.

I tried to do zero-copy of `string` as much as I could do.

# Real-world examples : TODO

- possibly expect some ImGui stuff
- possibly expect some SQLite stuff

# Licence

BSD 3-clause. I mean, at least for my original work.

Please consider the licences of included parts, including but not limited to :
- LuaJIT itself
- `*.lua` modules found in `lua-modules/`

# Build Instructions

You don't need to locally download `luajit-go`. It works as a "normal" Go module.

However, it uses `LuaJIT` and `Cgo`, so obviously, you will need LuaJIT to be installed somehow,
preferably find'able by `#cgo: pkg-config luajit`

### Build Instruction - Linux : TODO

I guess, try to have `pkg-config luajit` working. Do yourself a favor : use LuaJIT's latest commit (probably `v2.1` branch).

### Build Instruction - macOS X

I would like to advise to install LuaJIT HEAD via homebrew :
`brew install luajit --HEAD`

**HOWEVER**, on recent macOS (at least 11.6.2 for me, Intel CPU),
LuaJIT needs some special compile flags to work past simple cases, when used by Go ;

See https://github.com/LuaJIT/LuaJIT/issues/285#issuecomment-1451260836 and above...

The reason seems to be :
* LuaJIT would like to reserve some "relative jump area" pool of memory (2 GB on x86_64) for JITed code
* MacOS/Intel/Go seems to tight close together everything in virtual memory
* When the Go runtime let us a chance to load LuaJIT, it's already too late : LuaJIT is left with only 3% of the wanted 2 GB pool

That's without doing any work. As soon as you do some real work, those 3% goes down very quickly.

So under those circumstances, LuaJIT will bailout JIT'ing code.

I found it can be alleviated by building LuaJIT yourself :
```
git clone https://github.com/LuaJIT/LuaJIT.git
git checkout v2.1

export MACOSX_DEPLOYMENT_TARGET=11.0

make -C src/ libluajit.so LDFLAGS="-Wl,-image_base,0x488800000000"
make install
```

I choose 0x488800000000 hasardly, no specifics here. It's near the top of the 48 bits (~47,5 bits), with enough room before the 48 bits to not constrain the mcode region.

Then, please ensure that LuaJIT is find'able by `pkg-config luajit`.

All this to say that, whatever Go Module (this one or another) you would import to use LuaJIT in Go,
chances are that you would need to anyway specially-build LuaJIT as shown above.

P.S. : The above was for Intel Mac. For ARM Mac, see https://github.com/LuaJIT/LuaJIT/issues/285, people are also discussing patch(es).

### Build Instruction - Windows : TODO