FROM golang:1.12

LABEL "com.github.actions.name"="Auto Approve PR"
LABEL "com.github.actions.description"="Auto approve PRs"
LABEL "com.github.actions.icon"="thumbs-up"
LABEL "com.github.actions.color"="green"

WORKDIR /go/src/app
COPY . .
RUN GO111MODULE=on go build -o autoapprover .
ENTRYPOINT ["/go/src/app/autoapprover"]
