package main

import (
  "fmt"
  "github.com/google/go-github/github"
  "golang.org/x/oauth2"
  "strings"
  "io/ioutil"
  "time"
)

type tokenSource struct {
  token *oauth2.Token
}

func (t *tokenSource) Token() (*oauth2.Token, error){
  return t.token, nil
}

func checkCommits(CommitsURL string, client *github.Client) (bool, string) {
  newURL := strings.Replace(CommitsURL, "{/sha}", "", 1)

  opt := &github.CommitsListOptions {
    Since:  time.Now().Add(- (3 * 30 * 24 * time.Hour) ),
  }

  splURL := strings.Split(newURL, "/")

  commits, _, err := client.Repositories.ListCommits(splURL[4], splURL[5], opt)
  if err != nil {
    fmt.Println( err )
  }

  if len(commits) > 0 {
    return true, newURL
  }
  return false, newURL
}

func getStarred(user string, client *github.Client) ([]github.Repository, error) {
  var starred []github.Repository
  moar := true

  xPage := 0
  for moar {
    opts := &github.ActivityListStarredOptions {
        ListOptions: github.ListOptions{
          Page:    xPage,
          PerPage: 99,
        },
    }
    repos, _, err := client.Activity.ListStarred(user, opts)
    if err != nil {
      return nil, err
    }
    for _, rp := range repos {
      starred = append(starred, rp)
    }
    if len(repos) < 99 {
      moar = false
    }
    xPage = xPage + 1
  }
  return starred, nil // second 'nil' coz "no errors". facepalm.
}

func main() {
  token, err := ioutil.ReadFile("starsguard.conf")
  if err != nil {
    fmt.Println( err )
  }

  ts := &tokenSource {
    &oauth2.Token{AccessToken: string(token)},
  }

  tc := oauth2.NewClient(oauth2.NoContext, ts)

  client := github.NewClient(tc)

  starred, err := getStarred("p0123n", client)

  if err != nil {
    fmt.Println( err )
  }

  for k:=range starred {
    result, url := checkCommits(*starred[k].CommitsURL, client)

    if result == false {
      url = strings.Replace(url, "api.", "", 1)
      url = strings.Replace(url, "/repos", "", 1)
      url = strings.Replace(url, "/commits", "", 1)
      fmt.Println("Dead: ", url)
    }
  }
}
