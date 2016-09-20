# log
Multiple access log library for Go programs.
Its first goal was to learn how to correctly test a package

## Example

```Go
package main

import (
  "fmt"
  "sync"
  "github.com/Lerenn/log"
  "time"
)

func main() {
  // Prepare log ---------------------------------------------------------------
  // Create a new log instance
  l := log.New()

  // Launch log parallelization
  l.Start("/tmp/test.log")

  // Clear log file
  l.Clear()

  // Write synchronously in log ------------------------------------------------
  l.Print("Hello")
  // Should be printed after the previous line, never before because of synchronisation :
  fmt.Println("World")

  // Log parallelization example -----------------------------------------------
  fmt.Println("## Log parallelization example")

  var wg sync.WaitGroup
  wg.Add(2)
  go func() {
    defer wg.Done()
    for i := 0; i < 10; i++ {
      l.Print("1..")
      fmt.Println("..2")
      time.Sleep(time.Millisecond)
    }
  }()
  go func() {
    defer wg.Done()
    for i := 0; i < 10; i++ {
      l.Print("3..")
      fmt.Println("..4")
      time.Sleep(time.Millisecond)
    }
  }()

  wg.Wait()
  fmt.Println("Done.")
  // ---------------------------------------------------------------------------
}
```
