require('lib.class')

local App = class()

function App:build()
	print('no implementation for build method')
end

function App:run()
	print('no implementation for run method')
end

function App:cleanup()
	print('no implementation for cleanup method')
end

return App

