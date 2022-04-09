FROM golang:latest AS build
WORKDIR /app
COPY go.mod go.sum main.go ./
#RUN export GOPROXY=https://goproxy.io,direct
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download && go mod tidy

RUN go build -v

FROM gcr.io/distroless/base-debian11 AS run
#RUN apt update && apt -y upgrade 
#RUN apt -y install ca-certificates
COPY --from=build /app/pomotodo-stats /usr/local/bin/pomotodo-stats
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT [ "/usr/local/bin/pomotodo-stats" ]
#CMD ["/usr/local/bin/pomotodo-stats"]
