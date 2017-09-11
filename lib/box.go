package mp4

import (
  "fmt"
  "math"
  "errors"
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

type MediaHeaderBox struct {
  box FullBox
  creation_time uint64
  modification_time uint64
  timescale uint32
  duration uint64
  language string
}

type HandlerBox struct {
  box FullBox
  handler_type string
  name string
}

type AVCcBox struct {
  box Box
  version uint8
  profile uint8
  level uint8
  size_len uint8
  sps [][]byte
  pps [][]byte
}

type PixelAspectRatioBox struct {
  box Box
  hSpacing uint32
  vSpacing uint32
}

type VideoSampleDescription struct {
  width uint16
  height uint16
  depth uint16
}

type ESDescriptor struct {
  tag uint8
  len uint8
  id uint16
  config []byte
}

type ElementaryStreamDescBox struct {
  box Box
  esd ESDescriptor
}

type SoundSampleDescription struct {
  channels uint16
  sample_size uint16
  sample_rate float32
}

type SampleEntry struct {
  box Box
  data_ref_index uint16
  sample_desc interface{}
  extensions []interface{}
}

type SampleDescriptionBox struct {
  box FullBox
  entry_count uint32
  entries []SampleEntry
}

type SampleTableBox struct {
  box Box
  stsd SampleDescriptionBox
}

type MediaInfoBox struct {
  box Box
  stbl SampleTableBox
}

type MediaBox struct {
  box Box
  mdhd MediaHeaderBox
  hdlr HandlerBox
  minf MediaInfoBox
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
  copy(fb.flags[:], data[1:4]);
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

func parseMediaHeaderBox(data []byte, b *Box) (*MediaHeaderBox, error) {
  fb,_ := parseFullBox(data, b);
  mhb := MediaHeaderBox{box: *fb};

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

  code := [3]byte{};

  code[0] = (data[0] >> 2) + 0x60;
  code[1] = (((data[0] & 0x03) << 3) | data[1] >> 5) + 0x60;
  code[2] = (data[1] & 0x1f) + 0x60;

  mhb.language = string(code[0:3]);

  return &mhb, nil;
}

func parseHandlerBox(data []byte, b *Box) (*HandlerBox, error) {
  fb,_ := parseFullBox(data, b);
  hb := HandlerBox{box: *fb};
  data = data[4:];
  hb.handler_type = string(data[4:8]);
  hb.name = string(data[8:b.Size - b.headerSize - 4]);
  return &hb, nil;
}

func parsePixelAspectRatioBox(data []byte, b *Box) (*PixelAspectRatioBox, error) {
  pasp := PixelAspectRatioBox{box: *b};
  pasp.hSpacing = binary.BigEndian.Uint32(data[0:4]);
  pasp.vSpacing = binary.BigEndian.Uint32(data[4:8]);
  return &pasp, nil;
}

func parseAVCcBox(data []byte, b *Box) (*AVCcBox, error) {
  avcc := AVCcBox{box: *b};

  avcc.version = data[0];
  avcc.profile = data[1];
  avcc.level = data[3];
  avcc.size_len = (data[4] & 0x03) + 1;

  nsps := int(data[5] & 0x1f);
  data = data[6:];

  for i := 0; i < nsps; i++ {
    len := binary.BigEndian.Uint16(data[0:2]);
    sps := make([]byte, len);
    copy(sps, data[2: len + 2]);
    avcc.sps = append(avcc.sps, sps);
    data = data[len + 2:];
  }

  npps := int(data[0]);
  data = data[1:];

  for i := 0; i < npps; i++ {
    len := binary.BigEndian.Uint16(data[0:2]);
    pps := make([]byte, len);
    copy(pps, data[2: len + 2]);
    avcc.pps = append(avcc.pps, pps);
    data = data[len + 2:];
  }

  return &avcc, nil;
}

func parseVideoSampleDesc(data []byte, entry *SampleEntry) (*VideoSampleDescription, error) {
  vsd := VideoSampleDescription{};
  vsd.width = binary.BigEndian.Uint16(data[16:18]);
  vsd.height = binary.BigEndian.Uint16(data[18:20]);
  vsd.depth = binary.BigEndian.Uint16(data[66:68]);
  // TODO add missing info

  data = data[70:];

  tsize := entry.box.headerSize + 8 + 70;

  for {
    b,err := parseBox(data);

    if (err != nil) {
      return nil, err;
    }

    fmt.Println("              -", b.Type);

    switch (b.Type) {
    case "avcC":
      avcc,_ := parseAVCcBox(data[b.headerSize:], b);
      entry.extensions = append(entry.extensions, *avcc);
    case "pasp":
      pasp,_ := parsePixelAspectRatioBox(data[b.headerSize:], b);
      entry.extensions = append(entry.extensions, *pasp);
    }

    tsize += b.Size;

    if (tsize == entry.box.Size) {
      break;
    }

    data = data[b.Size:];
  }

  return &vsd, nil;
}

func parseESDescriptor(data []byte) (*ESDescriptor, error) {
  d := ESDescriptor{};

  d.tag = data[0];
  d.len = data[4];
  d.id = binary.BigEndian.Uint16(data[5:7]);

  data = data[8:];
  tag := data[0];

  if (tag != 0x04) {
    return nil, errors.New("invalid descriptor tag");
  }

  len := data[4];
  data = data[4 + 14:];
  tag = data[0];

  if (tag != 0x05) {
    return nil, errors.New("invalid descriptor tag");
  }

  len = data[4];
  data = data[5:];
  d.config = make([]byte, len);
  copy(d.config, data[:len]);

  return &d, nil;
}

func parseElementaryStreamDescBox(data []byte, b *Box) (*ElementaryStreamDescBox, error) {
  esds := ElementaryStreamDescBox{box: *b};
  data = data[4:];

  d,err := parseESDescriptor(data);

  if (err != nil) {
    fmt.Println(err)
  }

  esds.esd = *d;

  //fmt.Println(d);

  return &esds, nil;
}

func parseSoundSampleDesc(data []byte, entry *SampleEntry) (*SoundSampleDescription, error) {
  ssd := SoundSampleDescription{};
  ssd.channels = binary.BigEndian.Uint16(data[8:10]);
  ssd.sample_size = binary.BigEndian.Uint16(data[10:12]);
  fixed := binary.BigEndian.Uint32(data[16:20]);
  ssd.sample_rate = float32(fixed) / float32(math.Pow(2, 16));

  tsize := entry.box.headerSize + 8 + 20;

  data = data[20:];

  for {
    b,err := parseBox(data);

    if (err != nil) {
      return nil, err;
    }

    fmt.Println("              -", b.Type);

    switch (b.Type) {
    case "esds":
      esds,_ := parseElementaryStreamDescBox(data[b.headerSize:b.Size], b);
      entry.extensions = append(entry.extensions, *esds);
    }

    tsize += b.Size;

    if (tsize == entry.box.Size) {
      break;
    }

    data = data[b.Size:];
  }

  return &ssd, nil;
}

func parseSampleDescBox(data []byte, b *Box) (*SampleDescriptionBox, error) {
  fb,_ := parseFullBox(data, b);
  sdb := SampleDescriptionBox{box: *fb};

  data = data[4:];

  sdb.entry_count = binary.BigEndian.Uint32(data[0:4]);

  data = data[4:];

  tsize := b.headerSize + 4 + 4;

  for i := 0; i < int(sdb.entry_count); i++ {
    b,err := parseBox(data);

    if (err != nil) {
      return nil, err;
    }

    fmt.Println("            -", b.Type);

    entry := SampleEntry{box: *b};

    data = data[b.headerSize:];

    entry.data_ref_index = binary.BigEndian.Uint16(data[6:8]);

    data = data[8:];

    switch (b.Type) {
    case "avc1":
      vsd,_ := parseVideoSampleDesc(data, &entry);
      entry.sample_desc = *vsd;
    case "mp4a":
      ssd,_ := parseSoundSampleDesc(data, &entry);
      entry.sample_desc = *ssd;
    }

    sdb.entries = append(sdb.entries, entry);

    tsize += b.Size;

    if (tsize == sdb.box.box.Size) {
      break;
    }

    data = data[b.Size:];
  }

  return &sdb, nil;
}

func parseSampleTableBox(data []byte, b *Box) (*SampleTableBox, error) {
  stb := SampleTableBox{box: *b};

  tsize := b.headerSize;

  for {
    b,err := parseBox(data);

    if (err != nil) {
      return nil, err;
    }

    fmt.Println("          -", b.Type);

    switch b.Type {
    case "stsd":
      sdb,_ := parseSampleDescBox(data[b.headerSize:], b);
      stb.stsd = *sdb;
    }

    tsize += b.Size;

    if (tsize == stb.box.Size) {
      break;
    }

    data = data[b.Size:];
  }

  return &stb, nil;

}

func parseMediaInfoBox(data []byte, b *Box) (*MediaInfoBox, error) {
  mib := MediaInfoBox{box: *b};

  tsize := b.headerSize;

  for {
    b,err := parseBox(data);

    if (err != nil) {
      return nil, err;
    }

    fmt.Println("        -", b.Type);

    switch b.Type {
    case "stbl":
      stb,_ := parseSampleTableBox(data[b.headerSize:], b);
      mib.stbl = *stb;
    }

    tsize += b.Size;

    if (tsize == mib.box.Size) {
      break;
    }

    data = data[b.Size:];
  }

  return &mib, nil;
}

func parseMediaBox(data []byte, b *Box) (*MediaBox, error) {
  mb := MediaBox{box: *b};

  tsize := b.headerSize;

  for {
    b,err := parseBox(data);

    if (err != nil) {
      return nil, err;
    }

    fmt.Println("      -", b.Type);

    switch b.Type {
    case "mdhd":
      mhb,_ := parseMediaHeaderBox(data[b.headerSize:], b);
      mb.mdhd = *mhb;
    case "hdlr":
      hb,_ := parseHandlerBox(data[b.headerSize:], b);
      mb.hdlr = *hb;
    case "minf":
      mib,_ := parseMediaInfoBox(data[b.headerSize:], b);
      mb.minf = *mib;
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

    fmt.Println("    -", b.Type);

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

    fmt.Println("  -", b.Type);

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
