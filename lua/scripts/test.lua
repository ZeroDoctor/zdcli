local app = require('lib.app')
local util = require('lib.util')

local test = app:extend()

function test:get_name(env_type)
  io.write('enter name: ')
  local name = io.read()
  print('hello '..name)

  io.write('is it the best of both worlds? (y/n) ')
  local ans = io.read()

  if ans == 'y'then
    print('whoop whoop!')
  elseif ans == 'n' then
    print('oh on!')
  else
    print('why...')
  end
      
  
end

return test
