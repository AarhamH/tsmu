package essentials

/*
  PageId Structure
  |FID(32b)|PN(16b)|**(16b)|
  FID = File Identifier
  PN = Page Number
*/

const PAGE_SIZE = 4096 // 4K pages

type PageId struct {
  value uint64
}

func NewPageId(fileId int, page_num uint64) *PageId {
  computedValue := uint64(fileId) << uint32(32) | page_num << uint16(16)
  return &PageId{value: computedValue}
}

func (pid *PageId) IsPageIdValid() bool { return pid.value != ^uint64(0) }
func (pid *PageId) GetPageNumber() uint32 { return uint32(pid.value & 0xffffffff) >> 16 }
func (pid *PageId) GetFileId() uint32 { return uint32(pid.value >> uint32(32)) }


/*
  RecordId structure
  |FID(32b)|PN(16b)|SLOT(16b)|
  FID = File Identifier
  PN = Page Number
  SLOT = Slot Identifier
*/
type RecordId struct {
  PageId
}

func NewRecordId(pid *PageId, slot_id uint64) *RecordId {
  computedValue := pid.value | slot_id
  return &RecordId{PageId: PageId{value: computedValue}}
}

func (rid *RecordId) IsRecordIdValid() bool { return rid.value != ^uint64(0) }
func (rid *RecordId) GetSlotId() uint32 { return uint32(rid.value) & 0xffff }
