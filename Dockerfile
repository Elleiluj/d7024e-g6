FROM larjim/kademlialab:latest
RUN apt-get clean
RUN apt-get update && apt-get install -y iputils-ping
CMD bash

#FROM larjim/kademlialab:latest

#FROM alpine:latest

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
#WORKDIR /app

# Download Go modules
#COPY go.mod go.sum ./
#RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
#COPY main/*.go ./

# Build
#RUN CGO_ENABLED=0 GOOS=linux go build -o /repository

# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can (optionally) document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose

# Run
#CMD [ "/repository" ]
