local mt_cdata = {}
mt_cdata.__tojson = function(t)
    local ptr = tostring(ffi.cast("intptr_t", t.cdata_ptr)):sub(1,-3)
    return sprintf('{"cdata_ptr":%s,"cdata_len":%d}', ptr, t.cdata_len)
end

function json_cdata(ptr, len)
    return setmetatable({cdata_ptr = ptr, cdata_len = len}, mt_cdata)
end
