module github.com/abu-lang/abusim-core/abusim-coordinator

go 1.16

require (
	github.com/abu-lang/abusim-core/schema v1.0.0
	github.com/gorilla/mux v1.8.0
	github.com/rs/cors v1.8.0
)

replace github.com/abu-lang/abusim-core/schema => ../schema
