# To install prerequisits:
#
# To install redoc-cli:
# $ npm install
#
# To install oapi-codegen in $GOPATH/bin, go outside this go module:
# $ go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen

API_REST_SPEC=./docs/openapi/openapi.yaml
API_REST_CODE_GEN_LOCATION=./docs/openapi/generated/oapigen/oapigen.go
API_REST_DOCO_GEN_LOCATION=./docs/openapi/generated/doc.html

# Open API Makefile targets
oapi-validate:
	./node_modules/.bin/oas-validate -v ${API_REST_SPEC}

oapi-doc: oapi-validate
	./node_modules/.bin/redoc-cli bundle ${API_REST_SPEC} -o ${API_REST_DOCO_GEN_LOCATION}


test:
	go test -v ./...

dev-worker:
	gow run cmd/worker/main.go

dev-server:
	gow run cmd/server/main.go

worker:
	go run cmd/worker/main.go

server:
	go run cmd/server/main.go