local app = require('lib.app')
local util = require('lib.util')
local file = require('lib.file')

local script = app:extend()

function script:seaweedfs()
	-- TODO: add version option
	local version = "3.16"

	if util:is_windows() then
		local path = '../bin/swfs.zip'

		util:check_exec(
		'curl -o "../bin/swfs.zip" -L https://github.com/chrislusf/seaweedfs/releases/download/'..version..'/windows_amd64.zip',
		'powershell.exe -nologo -noprofile -command '..
		'"Expand-Archive -Force \''..path..'\' \''..file.get_parent_dir(path)..'\'"',
		'del ../bin/swfs.zip'
		)

		return
	end

	-- TODO: add arch option
	local path = '../bin/swfs.tar.gz'

	util:check_exec(
	'curl -o "../bin/swfs.tar.gz" -L https://github.com/chrislusf/seaweedfs/releases/download/'..version..'/linux_amd64.tar.gz',
	'tar -xvf '..path..' -C '..file.get_parent_dir(path),
	'rm ../bin/swfs.tar.gz'
	)
end

return script

