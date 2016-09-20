# Golang configuration library

Simple configuration library for Go programs. For simple configuration files with just `[sections]` and `token=value`.

This library has been made for self-education and for personal use but feel free to improve it.

## Example

### Config file

If we have `/etc/example/example.conf` with :

    [section]
    valueA=aaaa

### Source code

```Go
package main

import(
  "fmt"
  config "github.com/lerenn/go-config"
)

func main(){
  // Create a new struct
  conf := config.New()

  // Read the conf file
  if err := conf.Read("/etc/example/example.conf"); err != nil{
    fmt.Println("Error when reading file.")
    return
  }

  // Get a string
  val, err := conf.GetString("section", "valueA")
  if err == nil {
    fmt.Println("Value A = ", val)
  }
}
```
