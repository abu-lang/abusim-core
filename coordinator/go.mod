module abusim-coordinator

go 1.16

require (
	github.com/gorilla/mux v1.8.0
	github.com/rs/cors v1.8.0
	github.com/abu-lang/abusim-core/schema v0.0.0
)

replace github.com/abu-lang/goabu => ../../goabu
replace github.com/abu-lang/abusim-core/schema => ../schema
