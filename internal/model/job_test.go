package model
import (
	"encoding/json"
	"strings"
	"testing"
	"gopkg.in/yaml.v3"
)
func TestInput(t *testing.T) {
	yml := `---
job:
  package-manager: npm_and_yarn
  source:
    provider: github
    repo: dependabot/test
    directory: "/"
  existing-pull-requests:
    - pr-number: 123
      dependencies:
        - dependency-name: dep-a
          dependency-version: 1.0.0
    - dependency-name: dep-b
      dependency-version: 2.0.0
`
	var in Input
	yaml.Unmarshal([]byte(yml), &in)
}
func TestExistingPullRequests(t *testing.T) {
	yml := `---
job:
  package-manager: npm_and_yarn
  source:
    provider: github
    repo: a/b
    directory: "/"
  existing-pull-requests:
    - pr-number: 42
      dependencies:
        - dependency-name: antd
          dependency-version: 6.3.2
        - dependency-name: node-fetch
          dependency-removed: true
`
	var in Input
	yaml.Unmarshal([]byte(yml), &in)
	deps := *in.Job.ExistingPullRequests[0].Dependencies
	if!deps[1].DependencyRemoved {
		t.Error("removed")
	}
	b,_:=json.Marshal(in)
	if!strings.Contains(string(b),"dependency-removed"){
		t.Error("lost")
	}
}
