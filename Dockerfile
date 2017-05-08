FROM golang:alpine 
RUN apk add --no-cache git
RUN mkdir /app
ENV GOPATH /app
ADD . /app/
WORKDIR /app 
RUN go get github.com/jordan-wright/email
RUN go build -o rate_mailer . 
COPY myscript /etc/periodic/15min/myscript
RUN chmod +x /etc/periodic/15min/myscript 
CMD crond -l 2 -f
