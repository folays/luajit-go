--jit.off() jit.off(true, true)
--jit.flush() jit.flush(true,true)

--jit.opt.start(3)

--jit.opt.start("maxtrace=99999") -- NOPE
--jit.opt.start("maxrecord=0") -- (default 4000) =0 fix it
--jit.opt.start("maxirconst=0") -- (default 500) =0 or =99999 fix it
--jit.opt.start("maxside=99999") -- NOPE
--jit.opt.start("maxsnap=99999") -- NOPE

--jit.opt.start("hotloop=0","hotexit=0") -- NOPE
--jit.opt.start("tryside=9999") -- NOPE

--jit.opt.start("instunroll=99999") -- NOPE
--jit.opt.start("loopunroll=99999") -- NOPE
--jit.opt.start("callunroll=99999") -- NOPE
--jit.opt.start("recunroll=99999") -- NOPE

--jit.opt.start("sizemcode=4") -- NOPE (only tried 128)
--jit.opt.start("maxmcode=8192") -- NOPE (only tried 8192)

--print("real",jit.status())

if true then -- helps mcode_alloc() PRNG ; https://github.com/LuaJIT/LuaJIT/issues/285
    local seed = tonumber(tostring(_G):match("(0x%x+)"))
    math.randomseed(seed)

    for i=1,100 do end -- Force allocation of the first segment.
end

function _profile_message(color, format, ...)
    local file = io.open("/tmp/jit-msg", "a+")
    local str = string.format("\027[%sm",color)..string.format(format, ...).."\027[39;49m\n"
    file:write(str)
    file:close()
end

function profile_start()
    _profile_message("43", "======== PROFILE START ========")

    if false then -- trace OR dump
        require("jit.v").on("/tmp/jit-trace")
    else
        require("jit.dump").on("Atbsx", "/tmp/jit-dump")
    end

    do
        require("jit.p").start("vlpr", "/tmp/jit-profile") -- r is "raw sample counts instead of %"
    end
end

function profile_stop()
    require("jit.p").stop()
    require("jit.dump").off()
    require("jit.v").off()

    _profile_message("42", "======== PROFILE STOP ========")
end
