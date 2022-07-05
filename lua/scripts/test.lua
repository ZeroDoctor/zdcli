local app = require('lib.app')

local test = app:extend()

function test:get_name(env_type)
  io.write('enter name: ')
  local name = io.read()
  print('whoa there '..name)

  io.write('is it the best of both worlds? (y/n) ')
  local ans = io.read()

  if ans == 'y'then
    print('whoop! whoop!')
  elseif ans == 'n' then
    print('oh on!')
  else
    print('why...')
  end
end

function test:get_name_bash(env_type)
  print('enter name: ')
  local name = io.read('l')
  print('whoa there '..name)

  print('is it the best of both worlds? (y/n)')
  local ans = io.read('l')

  if ans == 'y'then
    print('whoop! whoop!')
  elseif ans == 'n' then
    print('oh on!')
  else
    print('why...')
  end
end

return test

