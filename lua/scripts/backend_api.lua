local app = require('lib.app')
local util = require('lib.util')

local backend_api = app:extend()

function backend_api:build(env_type)
    env_type = env_type or 'dev'
    if util:is_windows() then
        if env_type == 'prod' then
            print('production mode on windows not found')
            os.exit(0)
        end
		util:check_exec(
            'cmd /V /C "set GOOS=linux&set CGO_ENABLED=0&set GOARCH=amd64',
            'cd C:/Users/Daniel/Documents/lang/go/backend-api',
            'go build -a -installsuffix cgo -ldflags="-w -s" -o ./build/backend-api"'
        )
    else
        if env_type == 'dev' then
            print('development mode on linux not found')
            os.exit(0)
        end
        util:check_exec(
            'cd /mnt/c/Users/Daniel/Documents/lang/go/backend-api',
            'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o build/backend-api'
        )
    end
end

function backend_api:run_docker(env_type)
    env_type = env_type or 'dev'
	if util:is_windows() then
		os.execute('docker stop backend-api')
		os.execute('docker rm backend-api')
		util:check_exec('docker run -it --name backend-api -v C:/Users/Daniel/Documents/lang/go/backend-api/log:/go/bin/log -p 8078:8090 "trip/backend-api:latest"')
	else
		print('dev or prod mode on linux not found')
	end
end

function backend_api:build_docker(env_type)
    env_type = env_type or 'dev'
    if util:is_windows() then
        if env_type == 'prod' then
            print('production mode on windows not found')
            os.exit(0)
        end
        util:check_exec(
            'cd C:/Users/Daniel/Documents/lang/go/backend-api',
            'docker build -t "trip/backend-api:latest" -f "./build/Dockerfile" .'
        )
    else
        if env_type == 'dev' then
            print('development mode on linux not found')
            os.exit(0)
        end
        util:check_exec('mv /mnt/c/Users/Daniel/Documents/lang/go/backend-api/build/development.env /mnt/c/Users/Daniel/Documents/lang/go/backend-api/development.env')
        util:check_exec('mv /mnt/c/Users/Daniel/Documents/production_env/production.env /mnt/c/Users/Daniel/Documents/lang/go/backend-api/build/production.env')
        util:check_exec(
            'cd /mnt/c/Users/Daniel/Documents/lang/go/backend-api',
            'docker build -t smallwoods/backend-api -f build/Dockerfile .',
            'docker save -o build/backend-api smallwoods/backend-api',
            'rsync --inplace -e "ssh -i ~/.ssh/dans_key" -vzP build/backend-api app@api.smallwood.tools:/home/app/backend-api'
        )
    end
end

function backend_api:cleanup(env_type)
    env_type = env_type or 'dev'
    if util:is_windows() then
        print('no clean up on windows')
    else
        if env_type ~= 'prod' then
            print('no development mode on linux')
            os.exit(0)
        end
        util:check_exec(
            'cd /mnt/c/Users/Daniel/Documents/backend-api',
            'rm build/backend-api',
            'mv /mnt/c/Users/Daniel/Documents/lang/go/backend-api/development.env /mnt/c/Users/Daniel/Documents/lang/go/backend-api/build/development.env',
            'mv /mnt/c/Users/Daniel/Documents/lang/go/backend-api/build/production.env /mnt/c/Users/Daniel/Documents/production_env/production.env'
        )
    end
end

return backend_api

