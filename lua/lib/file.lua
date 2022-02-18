
local function slice_tbl(tbl, first, last)
  local sliced = {}
  
  for i = first or 1, last or #tbl, 1 do
    sliced[#sliced+1] = tbl[i]
  end

  return sliced
end

local function slice_str(str, first, last)
  local sliced = '' 

  for i = first or 1, last or #tbl, 1 do
    sliced = sliced..str:sub(i,i)
  end

  return sliced
end

local function file_exists(file)
  local f = io.open(file, "rb")
  if f then f:close() end
  return f ~= nil
end

local function lines_from(file)
  if not file_exists(file) then return {} end
  
  local lines = {}
  local count = 1
  for line in io.lines(file) do
    lines[count] = line
    count = count + 1
  end

  return lines
end

local function find_replace_file(file, find, replace) 
  local lines = lines_from(file)
  local result = {}

  for k, v in pairs(lines) do
    local start_index, end_index = string.find(v, find)
    if start_index then
      print('found ['..k..'] '..v)
      local a = slice_str(v, 1, start_index-1)
      local b = slice_str(v, end_index+1, #v)
      v = a..replace..b
      print('\treplace '..v)
    end

    result[k] = v
  end

  return result
end

function find_file(file, find, replace)
  local lines = find_replace_file(file, find, replace)

  local str = ""
  for i,line in ipairs(lines) do 
    str = str..line..'\n'
  end

  if replace ~= nil then
    local f = io.open(file, 'w+')
    f:write(str)
    f:close()
  end
end

return find_file

