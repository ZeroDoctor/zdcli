local util = require('lib.util')
local file = require('lib.file')
local ptr_tbl = require('lib.table_dump')
local argparse = require('lib.arg')

local function find_script(script)
	local str = string.gsub(script, '%.', '/')

	return file.exists('./scripts/'..str..'.lua')
end

local function set_flags(parser)
	-- parser:flag("-v --verbose"):count("*")
	parser:option("-a --arch", "set architecture")
	parser:option("-v --version", "set version")
end

local function main()
	local parser = argparse("build-app", "initial script to call other script")
	parser:argument("script", "name of script to run")
	parser:argument("funcs", "name of script to run"):args("*")
	set_flags(parser)
	local args = parser:parse()

	local app_name = args.script
	if not find_script(app_name) then
		util.perror('failed to find script named ['..app_name..']')
		return
	end
	local app = require('scripts.'..app_name)

	if #args.funcs < 1 then
		util.perror('must call a function in [script='..app_name..']')
		return
	end

	print('args: '..ptr_tbl(args, 2, false))
	for i=1, #args.funcs do
		local command = util:trim_all(args.funcs[i])

		if app[command] == nil then
			util.perror('failed to find [function='..command..']')
			return
		end

		print('step: '..command..'...')
		app[command](app, args)
	end

end

main()

