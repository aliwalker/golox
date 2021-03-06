generate_ast: tools/generate_ast.go
	@go run tools/generate_ast.go ./lox
	@echo "generating ast..."

lox: generate_ast
	@echo "building golox..."
	@go build golox.go

clean:
	@rm golox

test:
	@go test -v github.com/aliwalker/golox/lox