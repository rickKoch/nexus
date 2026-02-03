.PHONY: lint
lint:
	@./scripts/lint.sh segments
 	
.PHONY: fmt
fmt:
	goimports -l -w internal/

test:
	@./scripts/test.sh segments 