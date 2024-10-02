build:
	@go build -o bin/rck

generate: build
	@./bin/rck generate --file="./code/user.rocket" --language="go" --database="mysql"

gen: build
	@./bin/rck generate --file="./code/$(name).rocket" --language="go" --database="mysql"

execute: build
	@./bin/rck execute --query="$(name)" --args="$(args)"
