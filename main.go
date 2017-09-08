package main

import (
  "fmt"
  "os"
  "log"
  "bufio"

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

  r := bufio.NewReader(f);

  res,err := mp4.Parse(r);

  if (err != nil) {
    log.Fatal(err);
  }

  for _, x := range res.Boxes {
    switch x.(type) {
    case mp4.FileTypeBox:
      fmt.Println(x.(mp4.FileTypeBox));
    case mp4.FreeSpaceBox:
      fmt.Println(x.(mp4.FreeSpaceBox));
    }
  }

	err = f.Close();

  if (err != nil) {
		log.Fatal(err)
	}

  fmt.Println("done");
}
