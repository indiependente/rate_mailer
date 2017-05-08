FROM gliderlabs/alpine:latest
RUN apk add --no-cache go
RUN go get github.com/jordan-wright/email
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -o main .
COPY main /etc/periodic/15min/rate_mailer
CMD crond -l 2 -f