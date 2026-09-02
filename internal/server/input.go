package server

import (
"context"
"encoding/json"
"errors"
"log"
"net"
"net/http"
"reflect"
"time"

"github.com/moranricardo/cli/internal/model"
)

type credServer struct {
server *http.Server
data   model.Input
done   chan struct{}
}

func (s *credServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
return
}
r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB max para creds
defer r.Body.Close()

var input model.Input
if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
return
}

s.data = input
w.WriteHeader(http.StatusOK)

go func() {
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
_ = s.server.Shutdown(ctx)
close(s.done)
}()
}

func Input(listener net.Listener) (*model.Input, error) {
handler := &credServer{done: make(chan struct{})}
srv := &http.Server{
Handler:           handler,
ReadHeaderTimeout: 5 * time.Second,
ReadTimeout:       10 * time.Second,
WriteTimeout:      10 * time.Second,
}
handler.server = srv

log.Println("waiting for input on", listener.Addr())
if err := srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
return nil, err
}

<-handler.done // espera shutdown limpio
if reflect.DeepEqual(handler.data, model.Input{}) {
return nil, errors.New("no input received")
}
return &handler.data, nil
}
