# Create the building image for compiling
FROM golang:1.16-alpine as build

RUN mkdir /coordinator
WORKDIR /coordinator

COPY ./simulator/steel-simulator-common ./simulator/steel-simulator-common
COPY ./src ./src

WORKDIR /coordinator/simulator/steel-simulator-coordinator

COPY ./simulator/steel-simulator-coordinator/go.mod .
COPY ./simulator/steel-simulator-coordinator/go.sum .
RUN go mod download -x

COPY ./simulator/steel-simulator-coordinator .

RUN CGO_ENABLED=0 go build

# Create the final image with the executable
FROM scratch as exec

COPY --from=build /coordinator/simulator/steel-simulator-coordinator/steel-simulator-coordinator /coordinator

ENTRYPOINT [ "/coordinator" ]
