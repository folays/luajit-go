ffi.cdef[[
void *malloc(size_t size);
void free(void *ptr);
]]

function FFI_cdef(def)
    def = string.gsub(def, "#[^\n]*\n", "")
    ffi.cdef(def)
end

function FFI_malloc(ct, nb)
    local ctype = ffi.typeof(ct)

    if nb ~= nil then
        local len = ffi.sizeof(ctype) * nb
        local buf = ffi.C.malloc(len)
        ffi.fill(buf, len)
        return ffi.cast(ffi.typeof("$(&)[$]", ctype, nb), buf)
    else
        local len = ffi.sizeof(ctype) * 1
        local buf = ffi.C.malloc(len)
        ffi.fill(buf, len)
        return ffi.cast(ffi.typeof("$(&)", ctype), buf)
    end
end

function FFI_copy(dst, src)
    if ffi.sizeof(dst) ~= ffi.sizeof(src) then
        error("FFI_copy: size differ")
    end
    ffi.copy(dst, src, ffi.sizeof(dst))
end

local ctype_unsigned_char_ptr = ffi.typeof("unsigned char *")
function FFI_voidAdd(cdata, len)
    local ctype = ffi.typeof(cdata)
    return ffi.cast(ctype, ffi.cast(ctype_unsigned_char_ptr, cdata) + len)
end

function FFI_offsetof(header, member) -- return offset of a struct member
    --local offset = reflect.typeof(header).element_type:member(member).offset
    local offset = ffi.offsetof(ffi.typeof(header), member)
    return offset
end

function FFI_sizeof(header, member) -- returns size of a struct member
    local size = reflect.typeof(header).element_type:member(member).type.size
    return size
end

function FFI_sizeelem(array) -- for a "TYPE(&)[N]" array, return sizeof(TYPE)
    local refct = reflect.typeof(array)

    if refct.what == "ref" then
        refct = refct.element_type
    end

    local size = refct.element_type.size
    return size
end

function FFI_nb(array) -- for a "TYPE(&)[N]" array, return N
    return ffi.sizeof(array) / FFI_sizeelem(array)
end
