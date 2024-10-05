up:
    docker compose up -d
down:
    docker compose down -v    
run:
    cd cmd/app && air
build:
    go build -o ./tmp/main ./cmd/app
