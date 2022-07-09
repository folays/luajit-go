local function luajit_preload_loader(modname)
    local preload = package.preload[modname]
    if not preload then return end

    return function(modname)
        local mod = preload(modname)

        if type(mod) ~= "table" then
            mod = package.loaded[modname]
        end

        return mod
    end
end
table.remove(package.loaders,1)
table.insert(package.loaders,1, luajit_preload_loader)

local function luajit_embed_loader(modname)
    return luajit_embed:loadfile("lua-modules/"..modname..".lua")
end
table.insert(package.loaders, 2, luajit_embed_loader)
--while #package.loaders > 2 do table.remove(package.loaders) end

local mt_defer = {}
function mt_defer.__index(t,k)
    local modname = t[mt_defer]
    _G[modname] = nil
    setmetatable(t, nil)
    local mod = require(modname)
    --if type(mod) == "table" then _G[modname] = mod end
    _G[modname] = mod
    return mod[k]
end
local function defer(modname)
    local t = {[mt_defer] = modname} -- use mt_defer as a private key
    setmetatable(t, mt_defer)
    _G[modname] = t
end

for modname in pairs(package.preload) do
    if not _G[modname] then defer(modname) end
end
for _,entry in ipairs(luajit_embed:readdir("lua-modules")) do
    local modname = string.match(entry, "^(.*)%.lua$")
    if modname then defer(modname) end
end
