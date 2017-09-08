package mp4

import (
  "fmt"
  "bufio"
)

func Parse(r *bufio.Reader) error {
  buf := make([]byte, 4096);

  n,err := r.Read(buf);

  if (err != nil) {
    return err;
  }

  if (n == 0) {
    return nil;
  }

  fmt.Println("read", n);

  return nil;
}
