# Schemes.go

A tiny library for working with [Schemes](https://github.com/judah-caruso/schemes).


## Usage

```sh
go get github.com/judah-caruso/schemes.go
```

```go
package main

import (
   "github.com/judah-caruso/schemes.go"
   "io/ioutil"
   "fmt"
)

func main() {
   source, err := ioutil.ReadFile("my-scheme.svg")
   if err != nil {
      // ...
   }

   scheme, err := schemes.ReadScheme(source)
   if err != nil {
      // ...
   }

   fmt.Println("Opened scheme:", scheme)

   // Apply/work with scheme colors
   // Scheme.C0-C7     : SchemeColor
   // Scheme.Palette() : []SchemeColor

   fmt.Println("The hex value of c3 is:", scheme.C3.Hex())

   svg := schemes.ExportScheme(scheme)
   fmt.Println("Exported scheme:\n", svg)
}
```