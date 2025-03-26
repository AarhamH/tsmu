package storage

import "sync/atomic"

type Base struct {
  id int
  page_count *atomic.Uint32
}

// can represent basefile id or file descriptor
func (b Base) GetId() int { return b.id }
func (b Base) GetPageCount() *atomic.Uint32 { return b.page_count }
func (b Base) IncrementPageCount() { b.page_count.Add(1) }



