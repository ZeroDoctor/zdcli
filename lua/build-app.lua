local util = require('lib.util')

local function main() 
	if #arg < 1 then
		print('error: app name not given')
		return
	end

	local app_name = util:trim_all(arg[1])	
	local env_flag = ''
	local env_type = 'dev'

	if #arg >= 2 then
		env_flag = util:trim_all(arg[2])
	end
	
	local command_start = 2
	if env_flag == '-t' then
		env_type = util:trim_all(arg[3])
		command_start = 4
	end

	local app = require('scripts.'..app_name)
	if app == nil then
		print('error: failed to find app name: '..app_name)
		return
	end

	if command_start > #arg then
		print('error: command not found for app: '..app_name)
		return
	end

	print('currently in '..env_type)
	for i=command_start, #arg do
		local command = util:trim_all(arg[i])
                print(command)
		local fn = app[command]
		if fn == nil then
			print('error: failed to find command: '..command)
			return
		else
			print('step: '..command..'...')
			fn(env_type)
		end
	end
end

main()

