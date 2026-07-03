local app = require('lib.app')

local test = app:extend()

function test:func_win()
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

function test:func_unix()
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

