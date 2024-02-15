FROM golang:1.22-alpine

WORKDIR /home/goregex/web/app

CMD [ "sh", "-c", \
    "go mod tidy && go run ./" \
]
