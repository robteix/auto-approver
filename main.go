package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v24/github"
	"golang.org/x/oauth2"
)

func usage() {
	fmt.Println("Usage: maniaction [-user username] (merge|approve)")
	os.Exit(0)
}

func main() {
	var username = flag.String("user", "", "optional GitHub username to restrict this action to")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 2 {
		usage()
	}

	triggerName := os.Getenv("GITHUB_EVENT_NAME")
	if triggerName != "pull_request" {
		log.Printf("Ignoring trigger %s", triggerName)
		return
	}

	tok := os.Getenv("GITHUB_TOKEN")
	if tok == "" {
		log.Fatal("You must enable GITHUB_TOKEN access for this action")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: tok})

	oc := oauth2.NewClient(context.Background(), ts)
	c := github.NewClient(oc)

	b, err := ioutil.ReadFile(os.Getenv("GITHUB_EVENT_PATH"))
	if err != nil {
		log.Fatalf("could not read event info: %v", err)
	}

	var event github.PullRequestEvent
	if err := json.Unmarshal(b, &event); err != nil {
		log.Fatalf("could not unmarshal event info: %v", err)
	}

	pr := event.GetPullRequest()
	if *username != "" && *username != pr.GetUser().GetLogin() {
		fmt.Printf("Ignoring user %s; we only care about %s", *username, pr.GetUser().GetLogin())
		return
	}
	fullRepo := os.Getenv("GITHUB_REPOSITORY")
	repo := strings.SplitN(fullRepo, "/", 2)
	if len(repo) != 2 {
		log.Fatalf("Repo name %q invalid; should be in the form owner/repo", fullRepo)
	}

	switch flag.Arg(1) {
	case "approve":
		en := "APPROVE"
		req := &github.PullRequestReviewRequest{Event: &en}
		_, _, err = c.PullRequests.CreateReview(context.Background(), repo[0], repo[1], *event.Number, req)
		if err != nil {
			log.Fatalf("could not approve request: %v", err)
		}
		log.Printf("Approved PR %d", *event.Number)
	case "merge":
		pr := event.GetPullRequest()
		message := pr.GetTitle()
		opt := github.PullRequestOptions{MergeMethod: "rebase"}
		_, _, err = c.PullRequests.Merge(context.Background(), repo[0], repo[1], *event.Number, message, &opt)
		if err != nil {
			log.Fatalf("could not merge request: %v", err)
		}
		log.Printf("Merged PR %d", *event.Number)
	default:
		usage()
	}
}
