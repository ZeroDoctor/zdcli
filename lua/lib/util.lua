local module = {}

function module:trim_all(s)
	return s:match('^%s*(.-)%s*$')
end

local success = true

function module:check_exec(...)
	local args = {...}
	local command = ''
	for _, v in ipairs(args) do
		print('\texec: '..tostring(v)..'...')
		command = command..'&&'..tostring(v)
	end
	print('')
	command = command:sub(3, command:len())
	local code = os.execute(command)
	if code ~= success then
		module.perror('command failed exit code: '..tostring(code))
		os.exit(code)
	end
end

function module:capture_exec(cmd)
	local output = ''

	local h = io.popen(cmd, 'r')
	if h == nil then
		module.perror('failed to capture exec')
		return output
	end

	output = h:read('*a')
	h:close()

	return output
end

function module:slice_str(str, first, last)
	local sliced = ''

	for i = first or 1, last or str:len(), 1 do
		sliced = sliced..str:sub(i,i)
	end

	return sliced
end

function module.perror(str)
	print('[ERROR] '..str)
end

return module

