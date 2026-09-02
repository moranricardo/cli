package server

import (
"bytes"
"context"
"fmt"
"net"
"net/http"
"os"
"testing"
"time"

"github.com/moranricardo/cli/internal/model"
)

func newListener(t *testing.T) net.Listener {
t.Helper()
ip := ""
if os.Getenv("GOOS") == "darwin" {
ip = "127.0.0.1"
}
l, err := (&net.ListenConfig{}).Listen(context.Background(), "tcp", ip+":0")
if err != nil {
t.Fatalf("listener: %v", err)
}
return l
}

func TestInput(t *testing.T) {
tests := []struct {
name       string
method     string
body       string
wantStatus int
wantPM     string
}{
{"happy_path", "POST", `{"job":{"package-manager":"test"},"credentials":[{"credential":"value"}]}`, 200, "test"},
{"invalid_json", "POST", `{"job":`, 400, ""},
{"wrong_method", "GET", `{}`, 405, ""},
}

for _, tc := range tests {
t.Run(tc.name, func(t *testing.T) {
l := newListener(t)
inputCh := make(chan *model.Input, 1)
errCh := make(chan error, 1)

go func() {
input, err := Input(l)
if err != nil {
errCh <- err
return
}
inputCh <- input
}()

time.Sleep(50 * time.Millisecond)
url := fmt.Sprintf("http://%s", l.Addr())

var resp *http.Response
var err error
if tc.method == "POST" {
resp, err = http.Post(url, "application/json", bytes.NewReader([]byte(tc.body)))
} else {
req, _ := http.NewRequest(tc.method, url, nil)
resp, err = http.DefaultClient.Do(req)
}
if err != nil {
t.Fatalf("request failed: %v", err)
}
defer resp.Body.Close()

if resp.StatusCode != tc.wantStatus {
t.Fatalf("want %d got %d", tc.wantStatus, resp.StatusCode)
}

if tc.wantStatus == 200 {
select {
case input := <-inputCh:
if input.Job.PackageManager != tc.wantPM {
t.Errorf("want PM %s got %s", tc.wantPM, input.Job.PackageManager)
}
if input.Credentials[0]["credential"] != "value" {
t.Errorf("bad credential")
}
case <-time.After(2 * time.Second):
t.Fatal("server didn't shutdown")
}
}
})
}
}

func TestInput_BodyLimit(t *testing.T) {
l := newListener(t)
go func() { _, _ = Input(l) }()
time.Sleep(50 * time.Millisecond)
url := fmt.Sprintf("http://%s", l.Addr())

big := bytes.Repeat([]byte("a"), (1<<20)+1)
resp, err := http.Post(url, "application/json", bytes.NewReader(big))
if err != nil {
t.Fatal(err)
}
defer resp.Body.Close()
if resp.StatusCode != 400 {
t.Errorf("expected 400 for big body, got %d", resp.StatusCode)
}
}
