
local function slice_tbl(tbl, first, last)
	local sliced = {}

	for i = first or 1, last or #tbl, 1 do
		sliced[#sliced+1] = tbl[i]
	end

	return sliced
end

local function slice_str(str, first, last)
	local sliced = ''

	for i = first or 1, last or #str, 1 do
		sliced = sliced..str:sub(i,i)
	end

	return sliced
end

local function exists(file)
	local f = io.open(file, "rb")
	if f then f:close() end
	return f ~= nil
end

local function lines_from(file)
	if not exists(file) then return {} end

	local lines = {}
	local count = 1
	for line in io.lines(file) do
		lines[count] = line
		count = count + 1
	end

	return lines
end

local function find_replace_word(file, find, replace)
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

local function find_word(file, find, replace)
	local lines = find_replace_word(file, find, replace)

	local str = ""
	for _,line in ipairs(lines) do
		str = str..line..'\n'
	end

	if replace ~= nil then
		local f = io.open(file, 'w+')
		if f ~= nil then
			f:write(str)
			f:close()
		end
	end
end

local function get_parent_dir(path)
	if path:sub(#path, #path) == '/' or path:sub(#path, #path) == '\\' then
		path = slice_str(path, 1, #path-1)
	end

	local function last_slash()
		for i = #path or 1, 1, -1 do
			if path:sub(i, i) == '/' or path:sub(i, i) == '\\' then
				return i
			end
		end
		return -1
	end

	return slice_str(path, 1, last_slash()-1)
end

return {
	slice_str = slice_str,
	find_word = find_word,
	exists = exists,
	get_parent_dir = get_parent_dir,
}

