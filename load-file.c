#include "_cgo_export.h"

#include "load-file.h"

int file_dump_writer_c(lua_State *L, const void* p, size_t sz, void* ud) {
    return file_dump_writer_go(L, (void *)p, sz, ud);
}
