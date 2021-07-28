module steel-simulator-coordinator

go 1.16

replace (
	steel-lang => ../../src
	steel-simulator-common => ../steel-simulator-common
)

require (
	github.com/gorilla/mux v1.8.0
	github.com/rs/cors v1.8.0
	steel-lang v0.0.0
	steel-simulator-common v0.0.0
)
