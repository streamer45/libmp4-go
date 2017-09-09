package mp4

import (
  "fmt"
  "math"
  "encoding/binary"
)

const BOX_HDR_SZ = 8;
const BOX_HDR_SZ_EXT = 16;

type Box struct {
  Size uint64
  Type string
  headerSize uint64
}

type FullBox struct {
  box Box
  version uint8
  flags [3]byte
}

type FileTypeBox struct {
  box Box
  major_brand string
  minor_version uint32
  compatible_brands []string
}

type FreeSpaceBox struct {
  box Box
  data []byte
}

type MediaDataBox struct {
  box Box
  data []byte
}

type MovieHeaderBox struct {
  box FullBox
  creation_time uint64
  modification_time uint64
  timescale uint32
  duration uint64
  rate float32
  volume float32
  matrix [3][3]float32
  pre_defined [6][4]byte
  next_track_id uint32
}

type TrackHeaderBox struct {
  box FullBox
  creation_time uint64
  modification_time uint64
  track_id uint32
  duration uint64
  width float32
  height float32
}

type MediaBox struct {
  box Box
}

type TrackBox struct {
  box Box
  tkhd TrackHeaderBox
  mdia MediaBox
}

type MovieBox struct {
  box Box
  mvhd MovieHeaderBox
  tracks []TrackBox
}

func parseBox(data []byte) (*Box, error) {
  var b Box;

  hdr := data[0:BOX_HDR_SZ];

  b.headerSize = BOX_HDR_SZ;

  b.Size = uint64(binary.BigEndian.Uint32(hdr[0:4]));
  b.Type = string(hdr[4:8]);

  if (b.Size == 1) {
    b.headerSize = BOX_HDR_SZ_EXT;
    hdr = data[BOX_HDR_SZ:BOX_HDR_SZ_EXT]
    b.Size = binary.BigEndian.Uint64(hdr);
  }

  return &b,nil;
}

func parseFullBox(data []byte, b *Box) (*FullBox, error) {
  fb := FullBox{box: *b};
  fb.version = data[0];
  copy(fb.flags[:], data[1:]);
  return &fb, nil;
}

func parseFileTypeBox(data []byte, b *Box) (*FileTypeBox, error) {
  ftb := FileTypeBox{box: *b};

  ftb.major_brand = string(data[0:4]);
  ftb.minor_version = binary.BigEndian.Uint32(data[4:8]);

  data = data[8:];

  n := int((ftb.box.Size - ftb.box.headerSize - 8) / 4);

  ftb.compatible_brands = make([]string, n);

  for i := 0; i < n; i++ {
    ftb.compatible_brands[i] = string(data[i * 4: (i + 1) * 4]);
  }

  return &ftb, nil;
}

func parseFreeSpaceBox(data []byte, b *Box) (*FreeSpaceBox, error) {
  fsb := FreeSpaceBox{box: *b};
  fsb.data = data;
  return &fsb, nil;
}

func parseMediaDataBox(data []byte, b *Box) (*MediaDataBox, error) {
  mdb := MediaDataBox{box: *b};
  mdb.data = data;
  return &mdb, nil;
}

func parseMovieHeaderBox(data []byte, b *Box) (*MovieHeaderBox, error) {
  fb,_ := parseFullBox(data, b);
  mhb := MovieHeaderBox{box: *fb};

  data = data[4:];

  if (fb.version == 1) {
    mhb.creation_time = binary.BigEndian.Uint64(data[0:8]);
    mhb.modification_time = binary.BigEndian.Uint64(data[8:16]);
    mhb.timescale = binary.BigEndian.Uint32(data[16:20]);
    mhb.duration = binary.BigEndian.Uint64(data[20:28]);
    data = data[28:];
  } else {
    mhb.creation_time = uint64(binary.BigEndian.Uint32(data[0:4]));
    mhb.modification_time = uint64(binary.BigEndian.Uint32(data[4:8]));
    mhb.timescale = binary.BigEndian.Uint32(data[8:12]);
    mhb.duration = uint64(binary.BigEndian.Uint32(data[12:16]));
    data = data[16:];
  }

  fixed := binary.BigEndian.Uint32(data[0:4]);
  mhb.rate = float32(fixed) / float32(math.Pow(2, 16));

  fixed2 := binary.BigEndian.Uint16(data[4:6]);
  mhb.volume = float32(fixed2) / float32(math.Pow(2, 8));

  data = data[16:];

  for i:=0; i < 3; i++ {
    mhb.matrix[i][0] = float32(binary.BigEndian.Uint32(data[0:4])) / float32(math.Pow(2, 16));
    mhb.matrix[i][1] = float32(binary.BigEndian.Uint32(data[4:8])) / float32(math.Pow(2, 16));
    mhb.matrix[i][2] = float32(binary.BigEndian.Uint32(data[8:12])) / float32(math.Pow(2, 30));
    data = data[12:];
  }

  data = data[9 * 4 + 6 * 4:];

  mhb.next_track_id = binary.BigEndian.Uint32(data[0:4]);

  return &mhb, nil;
}

func parseTrackHeaderBox(data []byte, b *Box) (*TrackHeaderBox, error) {
  fb,_ := parseFullBox(data, b);
  thb := TrackHeaderBox{box: *fb};

  data = data[4:];

  if (fb.version == 1) {
    thb.creation_time = binary.BigEndian.Uint64(data[0:8]);
    thb.modification_time = binary.BigEndian.Uint64(data[8:16]);
    thb.track_id = binary.BigEndian.Uint32(data[16:20]);
    thb.duration = binary.BigEndian.Uint64(data[24:32]);
    data = data[32:];
  } else {
    thb.creation_time = uint64(binary.BigEndian.Uint32(data[0:4]));
    thb.modification_time = uint64(binary.BigEndian.Uint32(data[4:8]));
    thb.track_id = binary.BigEndian.Uint32(data[8:12]);
    thb.duration = uint64(binary.BigEndian.Uint32(data[16:20]));
    data = data[20:];
  }

  data = data[52:];

  fixed := binary.BigEndian.Uint32(data[0:4]);
  thb.width = float32(fixed) / float32(math.Pow(2, 16));

  fixed2 := binary.BigEndian.Uint32(data[4:8]);
  thb.height = float32(fixed2) / float32(math.Pow(2, 16));

  return &thb, nil;
}

func parseMediaBox(data []byte, b *Box) (*MediaBox, error) {
  mb := MediaBox{box: *b};

  tsize := b.headerSize;

  for {
    b,err := parseBox(data);

    if (err != nil) {
      return nil, err;
    }

    fmt.Println("   -", b.Type);

    switch b.Type {
    }

    tsize += b.Size;

    if (tsize == mb.box.Size) {
      break;
    }

    data = data[b.Size:];
  }

  return &mb, nil;

}

func parseTrackBox(data []byte, b *Box) (*TrackBox, error) {
  tb := TrackBox{box: *b};

  tsize := b.headerSize;

  for {
    b,err := parseBox(data);

    if (err != nil) {
      return nil, err;
    }

    fmt.Println("  -", b.Type);

    switch b.Type {
    case "tkhd":
      thb,_ := parseTrackHeaderBox(data[b.headerSize:], b);
      tb.tkhd = *thb;
    case "mdia":
      mb,_ := parseMediaBox(data[b.headerSize:], b);
      tb.mdia = *mb;
    }

    tsize += b.Size;

    if (tsize == tb.box.Size) {
      break;
    }

    data = data[b.Size:];
  }

  return &tb, nil;
}

func parseMovieBox(data []byte, b *Box) (*MovieBox, error) {
  mb := MovieBox{box: *b};

  tsize := b.headerSize;

  for {
    b,err := parseBox(data);

    if (err != nil) {
      return nil, err;
    }

    fmt.Println(" -", b.Type);

    switch b.Type {
    case "mvhd":
      mhb,_ := parseMovieHeaderBox(data[b.headerSize:], b);
      mb.mvhd = *mhb;
    case "trak":
      tb,_ := parseTrackBox(data[b.headerSize:], b);
      mb.tracks = append(mb.tracks, *tb);
    }

    tsize += b.Size;

    if (tsize == mb.box.Size) {
      break;
    }

    data = data[b.Size:];
  }

  return &mb, nil;
}
