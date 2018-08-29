generate_ast: tools/generate_ast.go
	@go run tools/generate_ast.go ./lox
	@echo "generating ast..."