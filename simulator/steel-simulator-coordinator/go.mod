module steel-simulator-coordinator

go 1.16

replace (
	steel-lang => ../../src
	steel-simulator-config => ../steel-simulator-config
)

require (
	github.com/gorilla/mux v1.8.0
	steel-lang v0.0.0
	steel-simulator-config v0.0.0
)
