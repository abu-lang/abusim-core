# Create the building image for compiling
FROM golang:1.16-alpine as build

RUN mkdir /coordinator
WORKDIR /coordinator

COPY ./abusim-core-dev/schema ./abusim-core-dev/schema
COPY ./goabu ./goabu

WORKDIR /coordinator/abusim-core-dev/coordinator

COPY ./abusim-core-dev/coordinator/go.mod .
COPY ./abusim-core-dev/coordinator/go.sum .
RUN go mod download -x

COPY ./abusim-core-dev/coordinator .

RUN CGO_ENABLED=0 go build

# Create the final image with the executable
FROM scratch as exec

COPY --from=build /coordinator/abusim-core-dev/coordinator/abusim-coordinator /coordinator

ENTRYPOINT [ "/coordinator" ]
