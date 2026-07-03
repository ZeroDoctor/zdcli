
local app = require('lib.app')
local util = require('lib.util')

local script = app:extend()

-- TODO: create seaweedfs that contains a zip file of updated zdcli
function script:zdcli(arg)
	if arg.os == "windows" then
		os.execute('mkdir %USERPROFILE%\\scripts\\bin')

		util:check_exec(
		'cd ..',
		'go build -o zd.exe .',
		-- 'copy .\\zd.exe %USERPROFILE%\\scripts',
		'xcopy .\\lua\\ %USERPROFILE%\\scripts /E/H'
		)

		print('copy .\\zd.exe %USERPROFILE%\\scripts')
		return
	end

	os.execute('mkdir ~/scripts/bin')

	util:check_exec(
	'cd ..',
	'go build -o zd .',
	'cp -r ./lua ~/scripts'
	)

	print('cp ./zd ~/scripts')
end

return script
