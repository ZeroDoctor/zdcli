local util = require('lib.util')
local file = require('lib.file')

local function find_script(script)
	local str = string.gsub(script, '%.', '/')

	return file.exists('./scripts/'..str..'.lua')
end

local function main()
	if #arg < 1 then
		util.perror('app name not given')
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

	if not find_script(app_name) then
		util.perror('failed to find script named: '..app_name)
		return
	end

	local app = require('scripts.'..app_name)
	if app == nil then
		util.perror('failed to find app name: '..app_name)
		return
	end

	if command_start > #arg then
		util.perror('command not found for app: '..app_name)
		return
	end

	print('currently in '..env_type)
	for i=command_start, #arg do
		local command = util:trim_all(arg[i])

		if app[command] == nil then
			util.perror('failed to find command: '..command)
			return
		end

		print('step: '..command..'('..env_type..')'..'...')
		app[command](app, env_type)
	end

end

main()

