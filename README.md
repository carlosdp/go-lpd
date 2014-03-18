# LPD Library for Go
go-lpd is a Line Printer Daemon Client/Server library for [Go](golang.org). It implements the [LPD Protocol](https://tools.ietf.org/html/rfc1179) allowing developers to easily create virtual network printers to create transparent applications that deal with printable files.

```go
package main

import (
  "fmt"
  "github.com/carlosdp/go-lpd"
  "io"
  "io/ioutil"
  "os"
)

func main() {
  server, err := lpd.NewServer(515)

  if err != nil {
    fmt.Println(err)
    return
  }

  client, err := lpd.NewClient("localhost:515")

  if err != nil {
    fmt.Println(err)
    return
  }

  file, err := ioutil.TempFile(os.TempDir(), "print")

  if err != nil {
    fmt.Println(err)
    return
  }

  err := client.PrintFile(file, "default") // Print file to "default" queue

  if err != nil {
    fmt.Println(err)
    return
  }

  receivedFile, err := server.ReceiveFile("default") // Receive file from "default" queue

  if err != nil {
    fmt.Println(err)
    return
  }

  defer receivedFile.Close()

  saveFile, err := os.Create("/some/path")

  if err != nil {
    fmt.Println(err)
    return
  }

  defer saveFile.Close()

  _, err := io.Copy(saveFile, receivedFile) // Copy the temp file to a permanent save location

  if err != nil {
    fmt.Println(err)
    return
  }

  fmt.Println("Printed file to: ", saveFile.Name())
}
```
