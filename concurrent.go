package main

import(
  "fmt"
  "io"
  "io/ioutil"
  "net/http"
  "os"
  "time"
)

func main() {
  start := time.Now()
  ch := make(chan string) // channel is a communication that allows one goroutine to pass values of a specified type to another goroutine
  for _, url := range os.Args[1:] {
    go fetch(url, ch) // goroutine concurrent function execution
  }
  for range os.Args[1:] {
    fmt.Println(<-ch) // receive from channel ch
  }
  fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
  start := time.Now()
  resp, err := http.Get(url)
  if err != nil {
    ch <- fmt.Sprint(err) // send to channel ch
    return
  }

  nbytes, err := io.Copy(ioutil.Discard, resp.Body)
  resp.Body.Close() // don't leak resources
  if err != nil {
    ch <- fmt.Sprintf("while reading %s: %v", url, err)
    return
  }
  ch <- fmt.Sprintf("%.2fs, %7d, %s", time.Since(start).Seconds(), nbytes, url)
}