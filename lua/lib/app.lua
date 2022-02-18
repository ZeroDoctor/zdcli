require('lib.class')

local App = class()

function App:build(env_type)
	print('this app does not have build method')
end

function App:build_docker(env_type)
	print('this app does not have build_docker method')
end

function App:run(env_type)
	print('this app does not have run method')
end

function App:run_docker(env_type)
	print('this app does not have run_docker method')	
end

function App:cleanup(env_type)
	print('this app does not have cleanup method')	
end

return App

