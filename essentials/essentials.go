package essentials

/*
  PageId Structure
  |-- 32-bit --|--16-bit--|--16-bit--|
  |   FileId   | Page Num |  Extra   |
*/
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

func NewPageId(fileId int, page_num uint64) *PageId {
  computedValue := uint64(fileId) << uint32(32) | page_num << uint16(16)
  return &PageId{value: computedValue}
}

/*
  PageId Structure
  |-- 32-bit --|--16-bit--|--16-bit--|
  |   FileId   | Page Num |  Slot Id |
*/
type RecordId struct {
  PageId
}

func (rid RecordId) IsRecordIdValid() bool {
  return rid.value != ^uint64(0)
}

func (rid RecordId) GetSlotId() uint32 {
  return uint32(rid.value) & 0xffff
}

func NewRecordId(pid *PageId, slot_id uint64) *RecordId {
  computedValue := pid.value | slot_id
  return &RecordId{PageId: PageId{value: computedValue}}
}

