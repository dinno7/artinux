
# Get help message
help:
  just --list


# Generate Swagger documentation
@gen-docs:
  swag init -d "pkg/server,pkg/response,internal/domain/entities,internal/adapters/http" -g server.go -o "./docs"
  echo "Formatiing..."
  swag fmt -d "pkg/server,pkg/response,internal/domain/entities,internal/adapters/http" -g server.go
