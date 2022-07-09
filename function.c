#include <lua.h>

#include "_cgo_export.h"

#include "function.h"

int function_call_c(lua_State *L) {
    int ret = function_call_go(L);

    if (ret == -1) {
        return lua_error(L);
    }

    return ret;
}
