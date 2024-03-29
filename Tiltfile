# Run server dependencies
docker_compose(
    './server/dev/docker-compose.yaml',
)

# Add labels to Docker services
dc_resource('postgres', labels=["Services"])
dc_resource('redis', labels=["Services"])
dc_resource('asynqmon', labels=["Services"])

# Run local API
local_resource(
    'scrimpton-api',
    dir='./server',
    serve_dir='./server',
    cmd='make build',
    serve_cmd='make run',
    ignore=['./server/bin', './server/proto', './server/var'],
    deps='./server',
    labels=['API'],
    resource_deps=['postgres', 'redis', 'asynqmon'],
    readiness_probe=probe(
        period_secs=15,
        http_get=http_get_action(port=8888, path="/api/metadata")
    ),
)

local_resource(
    'generate',
    auto_init=False,
    dir='./server',
    cmd='make generate',
    labels=['API'],
    deps='./server/proto',
)

# Run local UI
local_resource(
    'scrimpton-ui',
    dir='./gui',
    serve_dir='./gui',
    serve_cmd='npm run start',
    deps=['scrimpton-api'],
    labels=['UI'],
)

