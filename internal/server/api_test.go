package server

import (
"bytes"
"context"
"crypto/sha256"
"encoding/base64"
"encoding/hex"
"encoding/json"
"fmt"
"net/http"
"net/http/httptest"
"sync"
"testing"

"github.com/dependabot/cli/internal/model"
)

func Test_decodeWrapper(t *testing.T) {
t.Run("reject extra data", func(t *testing.T) {
_, err := decodeWrapper("update_dependency_list", []byte(`data: {"unknown": "value"}`))
if err == nil {
t.Error("expected decode would error on extra data")
}
})
}

func TestAPI_ServeHTTP(t *testing.T) {
t.Run("doesn't crash when unknown endpoint is used", func(t *testing.T) {
request := httptest.NewRequest("POST", "/unexpected-endpoint", nil)
response := httptest.NewRecorder()

api := NewAPI(nil, nil)
api.ServeHTTP(response, request)

if response.Code != http.StatusNotImplemented {
t.Errorf("expected status code %d, got %d", http.StatusNotImplemented, response.Code)
}
})
}

type Wrapper[T any] struct {
Data T `json:"data"`
}

func TestAPI_CreatePullRequest_ReplacesBinaryWithHash(t *testing.T) {
var stdout bytes.Buffer

api := NewAPI(nil, &stdout)
defer api.Stop()

content := base64.StdEncoding.EncodeToString([]byte("Hello, world!"))
hash := sha256.Sum256([]byte(content))
expectedHashedContent := hex.EncodeToString(hash[:])

createPullRequest := model.CreatePullRequest{
UpdatedDependencyFiles: []model.DependencyFile{
{
Content:         content,
ContentEncoding: "base64",
},
},
}
var body bytes.Buffer
if err := json.NewEncoder(&body).Encode(model.UpdateWrapper{Data: createPullRequest}); err != nil {
t.Fatalf("failed to encode request body: %v", err)
}

url := fmt.Sprintf("http://127.0.0.1:%d/create_pull_request", api.Port())
req, err := http.NewRequestWithContext(context.Background(), "POST", url, &body)
if err != nil {
t.Fatalf("failed to create request: %v", err)
}
req.Header.Set("Content-Type", "application/json")

client := &http.Client{}
resp, err := client.Do(req)
if err != nil {
t.Fatalf("failed to send request: %v", err)
}
defer resp.Body.Close()

if len(api.Errors) > 0 {
t.Fatalf("expected no errors, got %d errors: %v", len(api.Errors), api.Errors)
}

if len(api.Actual.Output) != 1 {
t.Fatalf("expected 1 output, got %d", len(api.Actual.Output))
}
if api.Actual.Output[0].Type != "create_pull_request" {
t.Fatalf("expected output type 'create_pull_request', got '%s'", api.Actual.Output[0].Type)
}
if api.Actual.Output[0].Expect.Data.(model.CreatePullRequest).UpdatedDependencyFiles[0].Content != expectedHashedContent {
t.Errorf("expected hashed content, got '%s'", api.Actual.Output[0].Expect.Data.(model.CreatePullRequest).UpdatedDependencyFiles[0].Content)
}

var wrapper Wrapper[model.CreatePullRequest]
if err := json.NewDecoder(&stdout).Decode(&wrapper); err != nil {
t.Fatalf("failed to decode stdout: %v", err)
}
if wrapper.Data.UpdatedDependencyFiles[0].Content != content {
t.Errorf("expected stdout to contain the original content, got '%s'", stdout.String())
}
}

func TestAPI_compareRecordEcosystemMeta(t *testing.T) {
t.Run("matching ecosystem meta", func(t *testing.T) {
expect := []model.RecordEcosystemMeta{
{
Ecosystem: model.Ecosystem{
Name: "bundler",
PackageManager: &model.VersionManager{
Name:       "bundler",
Version:    "2.7.2",
RawVersion: "2.7.2",
},
},
},
}
actual := []model.RecordEcosystemMeta{
{
Ecosystem: model.Ecosystem{
Name: "bundler",
PackageManager: &model.VersionManager{
Name:       "bundler",
Version:    "2.7.2",
RawVersion: "2.7.2",
},
},
},
}

if err := compareRecordEcosystemMeta(expect, actual); err != nil {
t.Errorf("expected no error, got %v", err)
}
})

t.Run("mismatched ecosystem meta", func(t *testing.T) {
expect := []model.RecordEcosystemMeta{
{
Ecosystem: model.Ecosystem{
Name: "bundler",
},
},
}
actual := []model.RecordEcosystemMeta{
{
Ecosystem: model.Ecosystem{
Name: "npm_and_yarn",
},
},
}

if err := compareRecordEcosystemMeta(expect, actual); err == nil {
t.Error("expected error for mismatched ecosystem meta")
}
})

t.Run("compare via compare function", func(t *testing.T) {
meta := []model.RecordEcosystemMeta{
{
Ecosystem: model.Ecosystem{
Name: "go_modules",
PackageManager: &model.VersionManager{
Name:       "gomod",
Version:    "1.21",
RawVersion: "1.21",
},
},
},
}
expectWrapper := &model.UpdateWrapper{Data: meta}
actualWrapper := &model.UpdateWrapper{Data: meta}

if err := compare(expectWrapper, actualWrapper); err != nil {
t.Errorf("expected no error from compare, got %v", err)
}
})
}

func TestAPI_compareDependencySubmissionRequest(t *testing.T) {
t.Run("ignores detector version", func(t *testing.T) {
expect := model.DependencySubmissionRequest{
Detector: model.DetectorMeta{
Version: "1.2.3",
},
}
actual := model.DependencySubmissionRequest{
Detector: model.DetectorMeta{
Version: "4.5.6",
},
}

if compareDependencySubmissionRequest(expect, actual) != nil {
t.Error("expected detector version to be ignored")
}
if expect.Detector.Version != "1.2.3" {
t.Error("expected expect detector version to be unchanged")
}
if actual.Detector.Version != "4.5.6" {
t.Error("expected actual detector version to be unchanged")
}
})
}

func TestAPI_ThreadSafety(t *testing.T) {
api := NewAPI(nil, nil)
defer api.Stop()

var wg sync.WaitGroup
workers := 20
for i := 0; i < workers; i++ {
wg.Add(1)
go func() {
defer wg.Done()
req := httptest.NewRequest("POST", "/increment_metric", bytes.NewReader([]byte(`data: {"metric": "test"}`)))
rec := httptest.NewRecorder()
api.ServeHTTP(rec, req)
}()
}
wg.Wait()
}

func TestAPI_BodyLimit(t *testing.T) {
api := NewAPI(nil, nil)
defer api.Stop()

big := make([]byte, 11<<20)
req := httptest.NewRequest("POST", "/create_pull_request", bytes.NewReader(big))
rec := httptest.NewRecorder()
api.ServeHTTP(rec, req)

if len(api.Errors) == 0 {
t.Error("expected error when body exceeds limit")
}
}

func TestAPI_Complete(t *testing.T) {
api := NewAPI([]model.Output{{Type: "close_pull_request"}}, nil)
api.Complete()
if len(api.Errors) != 1 {
t.Fatalf("expected 1 error for unmet expectation, got %d", len(api.Errors))
}
}
