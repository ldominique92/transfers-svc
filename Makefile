generate-server:
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
	oapi-codegen --package=api --generate server,types,spec -o api.gen.go api-spec.yml
	go get

install:
	go get -u github.com/gorilla/mux