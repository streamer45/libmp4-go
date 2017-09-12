package main

import (
  "fmt"
  "os"
  "log"

  "encoding/json"

  mp4 "github.com/streamer45/mp4/lib"
)

func printUsage() {
  fmt.Println("Usage:");
  fmt.Println("  go run main.go input_file.mp4");
}

func main() {
  if (len(os.Args) != 2) {
    printUsage();
    return;
  }

  fname := os.Args[1];

  f,err := os.OpenFile(fname, os.O_RDONLY, 0600);

  if (err != nil) {
    log.Fatal(err)
  }

  res,err := mp4.Parse(f);

  if (err != nil) {
    log.Fatal(err);
  }

  js,e := json.Marshal(res.Boxes);

  if (e != nil) {
    fmt.Println(e);
  }

  jsFile,err := os.OpenFile("./out.json", os.O_RDWR | os.O_CREATE, 0600);

  if (err != nil) {
    log.Fatal(err)
  }

  _,err = jsFile.Write(js);

  if (err != nil) {
    log.Fatal(err)
  }

  fmt.Println("wrote full parsing result to ./out.json");

  err = jsFile.Close();

  if (err != nil) {
    log.Fatal(err)
  }

  err = f.Close();

  if (err != nil) {
    log.Fatal(err)
  }
}
