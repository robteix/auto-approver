FROM golang:1.12

LABEL "com.github.actions.name"="Auto Approver of PR"
LABEL "com.github.actions.description"="Auto approver of PRs"
LABEL "com.github.actions.icon"="thumbs-up"
LABEL "com.github.actions.color"="green"

WORKDIR /go/src/app
COPY . .
RUN GO111MODULE=on go build -o maniaction .
ENTRYPOINT ["/go/src/app/maniaction"]
