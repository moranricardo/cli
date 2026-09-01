package server

import (
"bytes"
"context"
"crypto/sha256"
"encoding/hex"
"encoding/json"
"errors"
"fmt"
"io"
"log"
"net"
"net/http"
"os"
"reflect"
"runtime"
"strings"
"sync"
"time"

"github.com/dependabot/cli/internal/model"
"gopkg.in/yaml.v3"
)

// API intercepts calls to the Dependabot API
type API struct {
mu              sync.Mutex
Expectations    []model.Output
Errors          []error
Actual          model.SmokeTest

server          *http.Server
cursor          int
hasExpectations bool
port            int
writer          io.Writer
}

// NewAPI creates a new API instance and starts the server
func NewAPI(expected []model.Output, writer io.Writer) *API {
fakeAPIHost := "127.0.0.1"
if runtime.GOOS == "linux" {
fakeAPIHost = "0.0.0.0"
if version, err := os.ReadFile("/proc/version"); err == nil && strings.Contains(string(version), "Microsoft") {
fakeAPIHost = "127.0.0.1"
}
}
if os.Getenv("FAKE_API_HOST") != "" {
fakeAPIHost = os.Getenv("FAKE_API_HOST")
}

port := "0"
if os.Getenv("FAKE_API_PORT") != "" {
port = os.Getenv("FAKE_API_PORT")
}

l, err := (&net.ListenConfig{}).Listen(context.Background(), "tcp", fakeAPIHost+":"+port)
if err != nil {
panic(err)
}

server := &http.Server{
ReadTimeout:       5 * time.Second,
ReadHeaderTimeout: 5 * time.Second,
WriteTimeout:      10 * time.Second,
IdleTimeout:       60 * time.Second,
}

api := &API{
server:          server,
Expectations:    expected,
writer:          writer,
cursor:          0,
hasExpectations: len(expected) > 0,
port:            l.Addr().(*net.TCPAddr).Port,
}
server.Handler = api

go func() {
if err := server.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
api.mu.Lock()
api.Errors = append(api.Errors, err)
api.mu.Unlock()
log.Println("Server error:", err)
}
}()

return api
}

// Port returns the port the API is listening on
func (a *API) Port() int {
return a.port
}

// Stop stops the server
func (a *API) Stop() {
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()
_ = a.server.Shutdown(ctx)
}

// Complete adds any remaining expectations to the error queue
func (a *API) Complete() {
a.mu.Lock()
defer a.mu.Unlock()
for i := a.cursor; i < len(a.Expectations); i++ {
exp := &a.Expectations[i]
a.Errors = append(a.Errors, fmt.Errorf("expectation not met: %v\n%v", exp.Type, exp.Expect))
}
}

// ServeHTTP handles requests to the server
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10MB limit to prevent DoS
data, err := io.ReadAll(r.Body)
if err != nil {
a.pushError(fmt.Errorf("failed to read body: %w", err))
return
}
if err = r.Body.Close(); err != nil {
a.pushError(fmt.Errorf("failed to close body: %w", err))
return
}

parts := strings.Split(r.URL.String(), "/")
kind := parts[len(parts)-1]

actual, err := decodeWrapper(model.OutputType(kind), data)
if err != nil {
a.pushError(err)
}

a.outputRequestData(model.OutputType(kind), actual)

if kind == "create_pull_request" && actual != nil {
createPR := actual.Data.(model.CreatePullRequest)
createPR.UpdatedDependencyFiles = replaceBinaryWithHash(createPR.UpdatedDependencyFiles)
actual.Data = createPR
} else if kind == "update_pull_request" && actual != nil {
updatePR := actual.Data.(model.UpdatePullRequest)
updatePR.UpdatedDependencyFiles = replaceBinaryWithHash(updatePR.UpdatedDependencyFiles)
actual.Data = updatePR
}

if actual == nil {
w.WriteHeader(http.StatusNotImplemented)
return
}

if kind == "increment_metric" || kind == "record_ecosystem_meta" {
return
}

if err := a.pushResult(model.OutputType(kind), actual); err != nil {
a.pushError(err)
return
}

a.mu.Lock()
hasExp := a.hasExpectations
a.mu.Unlock()

if !hasExp {
return
}

a.assertExpectation(model.OutputType(kind), actual)
}

func (a *API) assertExpectation(kind model.OutputType, actual *model.UpdateWrapper) {
a.mu.Lock()
if len(a.Expectations) <= a.cursor {
a.mu.Unlock()
a.pushError(fmt.Errorf("missing expectation"))
return
}
expect := a.Expectations[a.cursor]
a.cursor++
a.mu.Unlock()

if kind != expect.Type {
a.pushError(fmt.Errorf("type was unexpected: expected %v got %v", expect.Type, kind))
return
}

data, err := json.Marshal(expect.Expect)
if err != nil {
panic(err)
}
expected, err := decodeWrapper(expect.Type, data)
if err != nil {
panic(err)
}

if err = compare(expected, actual); err != nil {
a.pushError(err)
}
}

func (a *API) outputRequestData(kind model.OutputType, data *model.UpdateWrapper) {
if a.writer != nil {
if err := json.NewEncoder(a.writer).Encode(map[string]any{
"type": kind,
"data": data.Data,
}); err != nil {
log.Panicln("Failed to write to stdout: ", err)
}
}
}

func (a *API) pushError(err error) {
a.mu.Lock()
defer a.mu.Unlock()
escapedError := strings.NewReplacer("\n", "", "\r", "").Replace(err.Error())
log.Println(escapedError)
a.Errors = append(a.Errors, err)
}

func (a *API) pushResult(kind model.OutputType, actual *model.UpdateWrapper) error {
a.mu.Lock()
defer a.mu.Unlock()

output := model.Output{
Type:   kind,
Expect: *actual,
}
a.Actual.Output = append(a.Actual.Output, output)

if msg, ok := actual.Data.(model.MarkAsProcessed); ok {
a.Actual.Input.Job.Source.Commit = msg.BaseCommitSha
}

return nil
}

func decodeWrapper(kind model.OutputType, data []byte) (actual *model.UpdateWrapper, err error) {
actual = &model.UpdateWrapper{}
switch kind {
case "update_dependency_list":
actual.Data, err = decode[model.UpdateDependencyList](data)
case "create_pull_request":
actual.Data, err = decode[model.CreatePullRequest](data)
case "create_dependency_submission":
actual.Data, err = decode[model.DependencySubmissionRequest](data)
case "update_pull_request":
actual.Data, err = decode[model.UpdatePullRequest](data)
case "close_pull_request":
actual.Data, err = decode[model.ClosePullRequest](data)
case "mark_as_processed":
actual.Data, err = decode[model.MarkAsProcessed](data)
case "record_ecosystem_versions":
actual.Data, err = decode[model.RecordEcosystemVersions](data)
case "record_ecosystem_meta":
actual.Data, err = decode[[]model.RecordEcosystemMeta](data)
case "record_update_job_error":
actual.Data, err = decode[model.RecordUpdateJobError](data)
case "record_update_job_unknown_error":
actual.Data, err = decode[model.RecordUpdateJobUnknownError](data)
case "increment_metric":
actual.Data, err = decode[model.IncrementMetric](data)
default:
return nil, fmt.Errorf("unexpected output type: %s", kind)
}
return actual, err
}

func replaceBinaryWithHash(files []model.DependencyFile) []model.DependencyFile {
for i := range files {
file := &files[i]
if file.ContentEncoding == "base64" {
file.ContentEncoding = "sha256"
decoded, decodeErr := hex.DecodeString(file.Content)
var hash [32]byte
if decodeErr == nil {
hash = sha256.Sum256(decoded)
} else {
hash = sha256.Sum256([]byte(file.Content))
}
file.Content = hex.EncodeToString(hash[:])
}
}
return files
}

func decode[T any](data []byte) (T, error) {
var wrapper struct {
Data T `json:"data" yaml:"data"`
}
decoder := yaml.NewDecoder(bytes.NewBuffer(data))
decoder.KnownFields(true)
err := decoder.Decode(&wrapper)
if err != nil {
return *new(T), err
}
return wrapper.Data, nil
}

func compare(expect, actual *model.UpdateWrapper) error {
switch v := expect.Data.(type) {
case model.DependencySubmissionRequest:
act, ok := actual.Data.(model.DependencySubmissionRequest)
if !ok {
return unexpectedBody("dependency_submission_request")
}
return compareDependencySubmissionRequest(v, act)
case model.UpdateDependencyList:
return compareGeneric(v, actual.Data, "update_dependency_list")
case model.CreatePullRequest:
return compareGeneric(v, actual.Data, "create_pull_request")
case model.UpdatePullRequest:
return compareGeneric(v, actual.Data, "update_pull_request")
case model.ClosePullRequest:
return compareGeneric(v, actual.Data, "close_pull_request")
case model.RecordEcosystemVersions:
return compareGeneric(v, actual.Data, "record_ecosystem_versions")
case model.MarkAsProcessed:
return compareGeneric(v, actual.Data, "mark_as_processed")
case model.RecordUpdateJobError:
return compareGeneric(v, actual.Data, "record_update_job_error")
case model.RecordUpdateJobUnknownError:
return compareGeneric(v, actual.Data, "record_update_job_unknown_error")
case []model.RecordEcosystemMeta:
act, ok := actual.Data.([]model.RecordEcosystemMeta)
if !ok {
return unexpectedBody("record_ecosystem_meta")
}
return compareRecordEcosystemMeta(v, act)
default:
return fmt.Errorf("unexpected type: %s", reflect.TypeOf(v))
}
}

func compareGeneric[T any](expect T, actualAny any, kind string) error {
actual, ok := actualAny.(T)
if !ok {
return unexpectedBody(kind)
}
if reflect.DeepEqual(expect, actual) {
return nil
}
return unexpectedBody(kind)
}

func unexpectedBody(kind string) error {
return fmt.Errorf("unexpected body for %s", kind)
}

func compareDependencySubmissionRequest(expect, actual model.DependencySubmissionRequest) error {
expect.Detector.Version = ""
actual.Detector.Version = ""
if reflect.DeepEqual(expect, actual) {
return nil
}
return unexpectedBody("dependency_submission_request")
}

func compareRecordEcosystemMeta(expect, actual []model.RecordEcosystemMeta) error {
return compareGeneric(expect, actual, "record_ecosystem_meta")
}
