package main

import (
  "fmt"
  "os"
  "log"
  "bufio"

  mp4 "github.com/streamer45/mp4/lib"
)

func printUsage() {

}

func main() {
  if (len(os.Args) != 1) {
    printUsage();
  }

  fname := os.Args[1];

  f,err := os.OpenFile(fname, os.O_RDONLY, 0600);

  if err != nil {
		log.Fatal(err)
	}

  r := bufio.NewReader(f);

  err = mp4.Parse(r);

  if (err != nil) {
    log.Fatal(err);
  }

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

  fmt.Println("done");
}
