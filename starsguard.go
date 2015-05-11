package main

import (
  "fmt"
  "github.com/google/go-github/github"
  "golang.org/x/oauth2"
  "strings"
  "io/ioutil"
  "time"
)

// tokenSource is an oauth2.TokenSource
// which returns a static access token
type tokenSource struct {
  token *oauth2.Token
}

// Token implements the oauth2.TokenSource interface
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

  // list all repositories for the authenticated user
  repos, _, err := client.Activity.ListStarred("p0123n", nil)

  if err != nil {
    fmt.Println( err )
  }

  // fmt.Println(reflect.TypeOf(client))
  for k:=range repos {
    result, url := checkCommits(*repos[k].CommitsURL, client)

    if result == true {
      fmt.Println("Alive: ", url)
    } else {
      fmt.Println("Dead: ", url)
    }
  }
}
