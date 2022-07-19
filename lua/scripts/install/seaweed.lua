
local app = require('lib.app')
local util = require('lib.util')

local script = app:extend()

function script:hello_world()
	print('hello world!')
end

return script
