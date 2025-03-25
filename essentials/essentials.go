package essentials

type PageId struct {
  value uint64
}

func (pid PageId) IsPageIdValid() bool {
  return pid.value != ^uint64(0)
}

func (pid PageId) GetPageNumber() uint32 {
  return uint32(pid.value & 0xffffffff) >> 16
}

func (pid PageId) GetFileId() uint32 {
  return uint32(pid.value >> uint32(32))
}

func NewPageId(fileId int, value uint64) *PageId {
  computedValue := uint64(fileId) << uint32(32) | uint64(value) << uint16(16)
  p := PageId{value: computedValue}
  return &p
}
