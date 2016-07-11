package main

import (
  "fmt"
  "log"
  "net/http"
  "time"
)

type Broker struct {
  clients map[chan string]bool
  newClients chan chan string
  defunctClients chan chan string
  messages chan string
}

func (b *Broker) Start() {
  go func() {
    for {
      select {
      case s := <-b.newClients:
        b.clients[s] = true
        log.Println("Added new client")

      case s := <-b.defunctClients:
        delete(b.clients, s)
        close(s)

        log.Println("Removed client")

      case msg := <-b.messages:
        for s, _ := range b.clients { s <- msg }
        log.Printf("Broadcast message to %d clients", len(b.clients))
      }
    }
  }()
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  f, ok := w.(http.Flusher)
  if !ok {
    http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
    return
  }

  messageChan := make(chan string)
  b.newClients <- messageChan
  notify := w.(http.CloseNotifier).CloseNotify()

  go func() {
    <-notify
    b.defunctClients <- messageChan
    log.Println("HTTP connection just closed.")
  }()

  w.Header().Set("Content-Type", "text/event-stream")
  w.Header().Set("Cache-Control", "no-cache")
  w.Header().Set("Connection", "keep-alive")

  for {
    msg, open := <-messageChan
    if !open { break }
    fmt.Fprintf(w, "data: Message: %s\n\n", msg)
    f.Flush()
  }

  log.Println("Finished HTTP request at ", r.URL.Path)
}

func main() {
  b := &Broker{
    make(map[chan string]bool),
    make(chan (chan string)),
    make(chan (chan string)),
    make(chan string),
  }

  b.Start()

  http.Handle("/", b)

  go func() {
    for i := 0; ; i++ {
      b.messages <- fmt.Sprintf("%d - the time is %v", i, time.Now())
      log.Printf("Sent message %d ", i)
      time.Sleep(5 * 1e9)
    }
  }()

  http.ListenAndServe(":3000", nil)
}
