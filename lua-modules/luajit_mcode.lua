local mcode = {}

local ffi = require("ffi")

ffi.cdef[[
void *malloc(size_t size);
void free(void *ptr);

int getpid(void);

typedef struct lua_State lua_State;
lua_State *(luaL_newstate) (void);
void luaL_openlibs(lua_State *L);
int luaL_loadfile(lua_State *L, const char *filename);
void lua_call(lua_State *L, int nargs, int nresults);
int lua_pcall(lua_State *L, int nargs, int nresults, int errfunc);
const char *lua_tolstring(lua_State *L, int idx, size_t *len);
void lua_settop(lua_State *L, int idx);

void *mmap(void *addr, size_t len, int prot, int flags, int fd, uint64_t offset);
int munmap(void *addr, size_t len);
]]

function mcode.eat_memory(kB, nb) -- raw eat (1024) * kB * nb
    for i=1,nb do
        local buf = ffi.C.malloc(1024*kB)
        buf = ffi.cast(ffi.typeof("uint8_t(&)[$]", 1024*kB), buf)
        for i_kB=0,kB-1 do -- populate pages
            buf[1024*i_kB] = 0x42
        end
    end
end
function mcode.eat_100_MiB_by_1_MiB()
    eat_memory(1024, 100)
end

local pid = ffi.C.getpid()
local si_suffix_table = {K=1024, M=1024*1024, G=1024*1024*1024}
function mcode.vm_size()
    local total_MB

    local file = io.popen("vmmap "..pid)
    for line in file:lines() do
        local total, suffix = line:match("^TOTAL +(%d+.?%d*)([KMG])")
        if total then
            total_MB = total * si_suffix_table[suffix] / 1024 / 1024
            --print("total_MB", total_MB)
        end

        local range1, range2 = line:match("^[%w_() ]+ +(%x+)-(%x+) +%[")
        if range1 then
        end
    end
    file:close()

    return total_MB
end

local LJ_TARGET_JUMPRANGE = 31
local sizemcode = 32 * 1024
local maxmcode = 512 * 1024
local t_ranges
function mcode._ranges_populate(target) -- Figuring out the whole mcode from code excerpt from mcode_alloc()
    t_ranges = {}
    target = bit.band(target, bit.bnot(0xffffULL))
    --printf("TARGET 0x%x", target)
    --local range = (1u << (LJ_TARGET_JUMPRANGE-1)) - (1u << 21);
    local range = bit.lshift(1ULL, LJ_TARGET_JUMPRANGE-1) - bit.lshift(1ULL, 21)
    --printf("RANGE 0x%x", range)
    local max = tonumber(bit.lshift(1ULL, 52)-1)

    for i=1,1000000 do
        local hint = bit.band(math.random(0,max), (bit.lshift(1ULL, LJ_TARGET_JUMPRANGE) - 0x10000))
        --printf("xxxx: %30x", tonumber(hint))
        if hint + sizemcode < range+range then
            --printf("yyyy: %30x %x", hint + sizemcode, range+range)
            --printf(".... target %x + hint %x - range %x", target,hint,range)
            hint = target + hint - range
            --printf("HINT: %30x", hint)
            hint = tonumber(bit.rshift(hint, 16))
            t_ranges[hint] = true
        end
    end

    local ranges_nb = 0
    for k, v in pairs(t_ranges) do
        ranges_nb = ranges_nb + 1
    end
    --printf("TOTAL RANGES : %d", ranges_nb)
end

local PROT_NONE = 0
local MAP_ANON = 0x1000
local MAP_FIXED = 0x0010
local MAP_PRIVATE = 0x0002
function mcode.ranges_check()
    if not t_ranges then mcode._ranges_populate(lj_vm_exit_handler()) end

    local t_ranges_available = {[true]=0,[false]=0}

    for range in pairs(t_ranges) do
        local range = ffi.cast("void *", bit.lshift(ffi.cast("uint64_t", range), 16))
        --printf("MMAP 0x%x", range)
        local mmap = ffi.C.mmap(range, sizemcode, PROT_NONE, bit.bor(MAP_ANON,MAP_PRIVATE), -1, 0)
        --print("MMAP", mmap, mmap == range)
        ffi.C.munmap(mmap, sizemcode)

        t_ranges_available[mmap == range] = t_ranges_available[mmap == range] + 1
    end

    local t = {
        ok = t_ranges_available[true],
        ok_pct = t_ranges_available[true] / (t_ranges_available[true] + t_ranges_available[false]) * 100,
        fail = t_ranges_available[false],
    }
    return t
end

function mcode.diagnose(label)
    local vm_size = mcode.vm_size()
    local ranges = mcode.ranges_check()
    if not ranges_ok_previous then ranges_ok_previous = ranges.ok end
    logf("[%5s] MCODE RANGES: [available %5d (%5.2f%%) %+6d] [not available %5d] ; virtual memory size %9.3f MB",
            label,
            ranges.ok, ranges.ok_pct, ranges.ok - ranges_ok_previous,
            ranges.fail, vm_size
    )
    ranges_ok_previous = ranges.ok
end

return mcode
