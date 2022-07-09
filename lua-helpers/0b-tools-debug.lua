local Gtraceback = debug.traceback

function debug.traceback(...)
    local trace_lua = Gtraceback(...)
    local trace_go = luajit.stacktrace()

    trace_go = string.gsub(trace_go, "([^\n]*)\n", "\t%1\n")

    return trace_lua .. "\n" .. trace_go
end
luajit.set_error_handler(debug.traceback)
