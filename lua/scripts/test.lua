local app = require('lib.app')
local util = require('lib.util')

local test = app:extend()

function test:sup(env_type)
  print('enter name: ')
  local name = io.read()
  print('hello '..name)
end

return test

