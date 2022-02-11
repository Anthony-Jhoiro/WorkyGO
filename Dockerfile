FROM golang:1.17 as builder


WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /wf

FROM docker:dind

COPY --from=builder /wf .


ENTRYPOINT [ "./wf" ]
