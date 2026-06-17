-- Data-oriented ring buffer in LuaJIT.
-- No structs, no methods, no type system, no OOP.
-- Flat array, scalar head/count, freestanding functions.

local capacity = 1024

-- ringPush: store value at head, return new head and count.
-- head is 1-indexed (Lua convention). Capacity is 1024.
-- After writing, head wraps: (head % 1024) + 1
local function ringPush(buf, head, count, value)
    buf[head] = value
    head = (head % capacity) + 1
    if count < capacity then
        count = count + 1
    end
    return head, count
end

-- ringAvg: arithmetic mean of stored elements. 0 if empty.
-- Elements are stored contiguously starting at buf[1] when count < capacity.
local function ringAvg(buf, count)
    if count == 0 then
        return 0.0
    end
    local sum = 0.0
    for i = 1, count do
        sum = sum + buf[i]
    end
    return sum / count
end

-- ringReadRange: copy up to n most recent elements into result.
-- Returns the result array and number of elements written.
-- Go equivalent: idx = (head_go - n + i + capacity) % capacity
-- Lua head = head_go + 1, so: ((head - 1 - n + i + capacity) % capacity) + 1
local function ringReadRange(buf, head, count, n)
    if n <= 0 or count == 0 then
        return {}, 0
    end
    if n > count then
        n = count
    end
    local res = {}
    for i = 0, n - 1 do
        local idx = ((head - 1 - n + i + capacity) % capacity) + 1
        res[i + 1] = buf[idx]
    end
    return res, n
end

-- ringUtilization: ratio of stored count to capacity.
local function ringUtilization(count)
    return count / capacity
end

-- Main: read stdin, dispatch commands, write stdout
local storage = {}
for i = 1, capacity do
    storage[i] = 0.0
end
local head = 1   -- 1-indexed, points to next write position
local count = 0

for line in io.lines() do
    -- Trim whitespace
    line = line:match("^%s*(.-)%s*$")
    if line == "" then goto continue end

    local cmd, arg = line:match("^(%S+)%s*(.-)%s*$")
    if not cmd then goto continue end

    if cmd == "push" then
        local value = tonumber(arg)
        if value then
            head, count = ringPush(storage, head, count, value)
        end
    elseif cmd == "avg" then
        local avg = ringAvg(storage, count)
        io.write(string.format("%.1f\n", avg))
    elseif cmd == "readRange" then
        local n = tonumber(arg) or 1
        local res, written = ringReadRange(storage, head, count, n)
        for i = 1, written do
            io.write(string.format("%.1f\n", res[i]))
        end
    elseif cmd == "utilization" then
        local util = ringUtilization(count)
        io.write(string.format("%.4f\n", util))
    end

    ::continue::
end
