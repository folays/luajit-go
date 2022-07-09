local Gprint = print

function print(...)
    Gprint(...)
    io.flush()
end

function sprintf(format, ...)
    return string.format(format, ...)
end

function printf(format, ...)
    local s = sprintf(format, ...)
    print(s)
    io.flush()
end

function errorf(format, ...)
    local s = sprintf(format, ...)
    error(s)
end

function assertf(v, format, ...)
    if v then return end

    if not format then
        return error("assertion failed!")
    else
        local s = sprintf(format, ...)
        return error("assertion failed! "..s)
    end
end

function log(...)
    Gprint(...)
    io.flush()
end

function logf(format, ...)
    Gprint(sprintf(format, ...))
    io.flush()
end

if os.getenv("SILENCE") then
    print = function() end
    printf = function() end
end

function toupper(s) return string.upper(s) end
function tolower(s) return string.lower(s) end

local function _pack_pcall(...)
    local status, err = ...
    if not status then return false, err end
    return true, nil, select('#', ...)-1, {select(2, ...)}
end
function pack_pcall(f, ...) -- returns status, err, nresults, {results...}
    return _pack_pcall(pcall(f, ...))
end

function pack_pcall_string(code)
    local fn, err = loadstring(code)
    if err then return false, err end
    local status, err, nresults, results = pack_pcall(fn)
    return status, err, nresults, results
end

