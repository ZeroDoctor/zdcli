local module = {}

local env_type = 'dev'
function module:set_env(env) env_type = env end
function module:get_env() return env_type end

local is_windows = false
local success = true

os_check, err = io.popen('uname -o 2>/dev/null', 'r')
if os_check:read() == nil then
	is_windows = true
end

if is_windows then
    print('os: windows')
else
    print('os: unix')
end

function module:is_windows()
	return is_windows
end

function module:trim_all(s)
	return s:match('^%s*(.-)%s*$')
end

function module:check_exec(...) 
	local args = {...}
	local command = ''
	for i, v in ipairs(args) do
		print('\texec: '..tostring(v)..'...')
		command = command..'&&'..tostring(v)
	end
	print('')
	command = command:sub(3, command:len())
	local code = os.execute(command)
	if code ~= success then
		print('command failed exit code: '..tostring(code))
		os.exit(code)
	end
end

function module:capture_exec(cmd)
  local output = ''
  
  local h = io.popen(cmd, 'r')
  output = h:read('*a')
  h:close()
  
  return output
end

function module:file_exists(name)
	local f = io.open(name,'r')
	if f ~= nil then
		io.close(f)
		return true
	end

	return false
end

function module:slice_str(str, first, last)
  local sliced = '' 

  for i = first or 1, last or #tbl, 1 do
    sliced = sliced..str:sub(i,i)
  end

  return sliced
end

return module

