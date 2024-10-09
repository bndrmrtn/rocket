build:
	@go build -o bin/rck

generate: build
	@./bin/rck generate --file="./code/user.rocket" --language="go" --database="mysql"

gen: build
	@./bin/rck generate --file="./code/$(name).rocket" --out="*_gen.{ext}" --language="go" --database="mysql"

gend: build
	@./bin/rck generate --file="./code/$(name)" --out="*_gen.{ext}" --language="go" --database="mysql"

execute: build
	@./bin/rck execute --query="$(name)" --args="$(args)"
