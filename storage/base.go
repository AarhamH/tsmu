package storage

import (
	"log"
	"os"
	"sync/atomic"
)

type Base struct {
  id int 
  page_count uint32
}

// can represent basefile id or file descriptor
func (b Base) GetId() int { return b.id }
func (b Base) GetPageCount() uint32 { return b.page_count }
func (b Base) IncrementPageCount() { atomic.AddUint32(&b.page_count, 1) }

func NewBaseFile(fileName string) *Base {
  file, err := os.OpenFile(fileName, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0664)
  if err != nil {
    log.Fatal(err)
  }

  b := Base{ id: int(file.Fd()), page_count: 0 }
  
  err = file.Close()
  if err != nil {
    log.Fatal(err)
  }

  return &b;
}



