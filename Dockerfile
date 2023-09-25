FROM golang:alpine

#RUN apt-get update && apt-get install -y iputils-ping
#CMD bash


#FROM larjim/kademlialab:latest

# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#
# $ docker build . -t kadlab

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download


COPY main/*.go ./main/
COPY server/*.go ./server/
COPY client/*.go ./client/

# Build
RUN go build -o /kademlialab ./main

EXPOSE 8080

# Run
CMD [ "/kademlialab" ]
