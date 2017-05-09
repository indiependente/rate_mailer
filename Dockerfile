FROM golang:alpine 
RUN apk add --no-cache git curl
RUN mkdir /app
ENV GOPATH /app
ARG hostip
ENV HOSTIP=$hostip
ADD . /app/
WORKDIR /app 
RUN go get github.com/jordan-wright/email
RUN go build -o rate_mailer . 
COPY myscript /etc/periodic/hourly/myscript
RUN chmod +x /etc/periodic/hourly/myscript 
CMD crond -l 2 -f
