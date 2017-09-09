package main

import (
  "fmt"
  "os"
  "log"

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

  for _, x := range res.Boxes {
    switch x.(type) {
    case mp4.FileTypeBox:
      fmt.Printf("%+v\n", x.(mp4.FileTypeBox));
    case mp4.FreeSpaceBox:
      fmt.Printf("%+v\n", x.(mp4.FreeSpaceBox));
    case mp4.MediaDataBox:
      fmt.Printf("%+v\n", x.(mp4.MediaDataBox));
    case mp4.MovieBox:
      fmt.Printf("%+v\n", x.(mp4.MovieBox));
    }
  }

	err = f.Close();

  if (err != nil) {
		log.Fatal(err)
	}
}
