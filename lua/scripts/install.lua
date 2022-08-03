local app = require('lib.app')
local util = require('lib.util')
local file = require('lib.file')

local script = app:extend()

local function ca_crt(name)
	local path = '../bin/crt'

	os.execute('mkdir '..path)

	print('creating '..name..' ca...')
	util:check_exec(
	'openssl genrsa -aes256 -out '..name..'-ca.key 4096',
	'openssl req -x509 -new -nodes -key '..name..'-ca.key '..
	'-sha256 -days 1826 -out '..name..'-ca.crt -subj '..
	'"/CN=zd'..name..' ca/C=US/ST=Texas/L=Longview/O=ZeroDoc Solutions"'
	)
end

function script:seaweedfs(arg)
	-- TODO: add version option
	local version = "3.16"

	if arg.version ~= nil then
		version = arg.version
	end

	if arg.os == "windows" then
		local path = '../bin/swfs.zip'
		local link = 'https://github.com/chrislusf/seaweedfs/releases/download/'..version..'/windows_'..arg.arch..'.zip'

		util:check_exec(
		'curl -o "../bin/swfs.zip" -L '..link,
		'powershell.exe -nologo -noprofile -command '..
		'"Expand-Archive -Force \''..path..'\' \''..file.get_parent_dir(path)..'\'"',
		'del ../bin/swfs.zip'
		)

		return
	end

	local path = '../bin/swfs.tar.gz'
	local link = 'https://github.com/chrislusf/seaweedfs/releases/download/'..version..'/linux_'..arg.arch..'.tar.gz'

	util:check_exec(
	'curl -o "../bin/swfs.tar.gz" -L '..link,
	'tar -xvf '..path..' -C '..file.get_parent_dir(path),
	'rm ../bin/swfs.tar.gz'
	)

	ca_crt('swfs')
end

function script:consul(arg)
	if arg.os == "windows" then
		print('use WSL')
	end
	-- TODO: add version option
	local version = "1.12.3"

	if arg.version ~= nil then
		version = arg.version
	end


	local path = '../bin/consul.zip'
	local link = "https://releases.hashicorp.com/consul/"..version.."/consul_"..version.."_linux_"..arg.arch..".zip"

	util:check_exec(
	'curl -o "../bin/consul.zip" -L '..link,
	'unzip '..path..' -d '..file.get_parent_dir(path),
	'rm ../bin/consul.zip'
	)
end

function script:nomad(arg)
	if arg.os == "windows" then
		print('use WSL')
	end
	-- TODO: add version option
	local version = "1.3.2"

	if arg.version ~= nil then
		version = arg.version
	end


	local path = '../bin/nomad.zip'
	local link = "https://releases.hashicorp.com/nomad/"..version.."/nomad"..version.."_linux_"..arg.arch..".zip"

	util:check_exec(
	'curl -o "../bin/nomad.zip" -L '..link,
	'unzip '..path..' -d '..file.get_parent_dir(path),
	'rm ../bin/nomad.zip'
	)
end

return script

