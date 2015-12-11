FROM golang:1.5.1

# Setting up working directory
WORKDIR /go/src/github.com/alexstoick/budgie-backend/
Add . /go/src/github.com/alexstoick/budgie-backend/

# Get godeps from main repo
RUN go get github.com/tools/godep

# Restore godep dependencies
RUN godep restore

# Install
RUN go install github.com/alexstoick/budgie-backend/

# Setting up environment variables
ENV ENV dev

# My web app is running on port 8080 so exposed that port for the world
EXPOSE 3000
ENTRYPOINT ["/go/bin/budgie-backend"]
