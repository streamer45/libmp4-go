package mp4

const BOX_HDR_SZ = 8;
const BOX_HDR_SZ_EXT = 16;

type Box struct {
  Size uint64 `json:"size"`
  Type string `json:"type"`
  headerSize uint64
}

type FullBox struct {
  box Box
  version uint8
  flags [3]byte
}

type FileTypeBox struct {
  Box Box
  major_brand string
  minor_version uint32
  compatible_brands []string
}

type FreeSpaceBox struct {
  Box Box
  data []byte
}

type MediaDataBox struct {
  Box Box
  data []byte
}

type MovieHeaderBox struct {
  Box FullBox
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

type TimeToSampleBox struct {
  box FullBox
  entry_count uint32
  sample_count []uint32
  sample_delta []uint32
}

type SyncSampleBox struct {
  box FullBox
  entry_count uint32
  sample_number []uint32
}

type CompTimeToSampleBox struct {
  box FullBox
  entry_count uint32
  sample_count []int32
  sample_offset []int32
}

type SampleToChunkBox struct {
  box FullBox
  entry_count uint32
  first_chunk []int32
  samples_per_chunk []int32
  sample_desc_index []int32
}

type SampleSizeBox struct {
  box FullBox
  sample_size uint32
  sample_count uint32
  entry_size []uint32
}

type ChunkOffsetBox struct {
  box FullBox
  entry_count uint32
  chunk_offset []uint32
}

type SampleTableBox struct {
  box Box
  stsd SampleDescriptionBox
  stts TimeToSampleBox
  stss SyncSampleBox
  ctss CompTimeToSampleBox
  stsc SampleToChunkBox
  stsz SampleSizeBox
  stco ChunkOffsetBox
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
  Box Box
  mvhd MovieHeaderBox
  tracks []TrackBox
}
