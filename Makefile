# Usage: make [target]

.PHONY: help
help:  ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  make %-10s %s\n", $$1, $$2}'

up:  ## Start all services in detached mode
	docker-compose up -d --build

down:  ## Stop all services
	docker-compose down

logs:  ## Follow container logs
	docker-compose logs -f

clean: down  ## Stop services and remove volumes
	docker-compose down -v

mongo:  ## Access MongoDB shell
	docker exec -it mongodb mongosh -u admin -p admin

ui:  ## Open Mongo-Express in browser (macOS/Linux)
	@xdg-open http://localhost:8081 2>/dev/null || open http://localhost:8081