local function exec(...)
--   print not working?
  print(redis.call(unpack(arg)))
end
 
exec("FLUSHALL")
exec("HMSET", 2, "id", "hillaryclinton", "state", "new", "max", 100)
exec("HMSET", 3, "id", "potus44", "state", "new", "max", 3, "error", "e", "data", "abc")
exec("LPUSH", "genie:followfinder", 2)
exec("LPUSH", "genie:followfinder", 3)
exec("HGETALL", 2)
exec("HGETALL", 3)