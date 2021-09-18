# Create the building image for compiling
FROM golang:1.16-alpine as build

RUN mkdir /coordinator
WORKDIR /coordinator

COPY ./schema ./schema

WORKDIR /coordinator/abusim-coordinator

COPY ./abusim-coordinator/go.mod .
COPY ./abusim-coordinator/go.sum .
RUN go mod download -x

COPY ./abusim-coordinator .

RUN CGO_ENABLED=0 go build

# Create the final image with the executable
FROM scratch as exec

COPY --from=build /coordinator/abusim-coordinator/abusim-coordinator /coordinator

ENTRYPOINT [ "/coordinator" ]
