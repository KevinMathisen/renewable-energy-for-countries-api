FROM golang:1.19

LABEL maintainer=kevinnm@stud.ntnu.no
LABEL maintainer=raphaesl@stud.ntnu.no

# Copy source files into image
COPY ./cmd /go/src/app/cmd
COPY ./handlers /go/src/app/handlers
COPY ./utils /go/src/app/utils
COPY ./go.mod /go/src/app/go.mod

# Starting in following working div
WORKDIR /go/src/app/cmd

# Install external dependencies (firestone)
RUN go get assignment2/utils/db
RUN go get assignment2/handlers
RUN go get github.com/stretchr/testify/assert

# Compile executable
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags -static' -o server

# Application will run on port 8080
EXPOSE 8080

# Run executable
CMD ["./server"]