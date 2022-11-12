local util = require('lib.util')
local file = require('lib.file')
local ptr_tbl = require('lib.table_dump')
local argparse = require('lib.arg')
local table_dump = require('lib.table_dump')

local function find_script(script)
	local str = string.gsub(script, '%.', '/')

	return file.exists('./scripts/'..str..'.lua')
end

local function set_flags(parser)
	parser:flag("-v --verbose"):count("*")
	parser:option("--pwd", "get the current working directory (read-only)")
	parser:option("--os_i", "set operating system internally")
	parser:option("--arch_i", "set architecture internally")
	parser:option("--os", "set operating system (if blank os_i will be used)")
	parser:option("--arch", "set architecture (if blank arch_i wiill be used)")
	parser:option("--version", "set version")
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

	if args.arch == nil or args.arch == "" then
		args.arch = args.arch_i
	end

	if args.os == nil or args.os == "" then
		args.os = args.os_i
	end

	if args.verbose == 1 then
		print('args: '..ptr_tbl(args, 2, false))
	end

	for i=1, #args.funcs do
		local command = util:trim_all(args.funcs[i])

		if args.funcs[i] == 'ls' then
			print('\n-- available methods below:')
			print(table_dump(app["_"], 2), '\n--\n')
			return
		end

		if app[command] == nil then
			util.perror('failed to find [function='..command..']')
			return
		end

		print('step: '..command..'...')
		app[command](app, args)
	end

end

main()

