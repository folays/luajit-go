#include "_cgo_export.h"

#include "metatable.h"

int metatable_gc_c(lua_State *L) {
    return metatable_gc_go(L);
}
