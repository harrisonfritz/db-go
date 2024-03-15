##   specify the base image we need for our
## go application
FROM golang:1.16-alpine3.13
## We create an /app directory within our
## image that will hold our application source
## files
RUN mkdir /app
## Environment variable
ENV pgConnP="test"
RUN echo $pgConnP
## We copy everything in the root directory
## into our /app directory
ADD . /app
RUN ls
## We specify that we now wish to execute 
## any further commands inside our /app
## directory
WORKDIR /app
## we run go build to compile the binary
## executable of our Go program
# RUN go build pg-connect.go 
## Our start command which kicks off
## our newly created binary executable
# CMD ["sleep 3600"]
