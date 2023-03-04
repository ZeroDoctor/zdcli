local str = require('lib.str')

local function slice_tbl(tbl, first, last)
	local sliced = {}

	for i = first or 1, last or #tbl, 1 do
		sliced[#sliced+1] = tbl[i]
	end

	return sliced
end


local function is_dir(path)
	if type(path) ~= "string" then return false end
	return os.execute("cd "..path)
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


local function find_replace_output(file, find, replace)
	local lines = lines_from(file)
	lines = str.find_replace_word(lines, find, replace)

	local words = ""
	for _,line in ipairs(lines) do
		words = words..line..'\n'
	end

	if replace ~= nil then
		local f = io.open(file, 'w+')
		if f ~= nil then
			f:write(words)
			f:close()
		end
	end
end

local function get_parent_dir(path)
	if path:sub(#path, #path) == '/' or path:sub(#path, #path) == '\\' then
		path = str.slice(path, 1, #path-1)
	end

	local function last_slash()
		for i = #path or 1, 1, -1 do
			if path:sub(i, i) == '/' or path:sub(i, i) == '\\' then
				return i
			end
		end
		return -1
	end

	return str.slice(path, 1, last_slash()-1)
end

return {
	exists = exists,
	is_dir = is_dir,
	slice_tbl = slice_tbl,
	lines_from = lines_from,
	find_replace_output = find_replace_output,
	get_parent_dir = get_parent_dir
}

