package storage

import (
	"log"
	"os"
	"sync/atomic"
)

type Base struct {
  fd int // fd represents the file descriptor after opening a file  
  page_count uint32 // page_count is treated as an atomic integer
}

func NewBaseFile(fileName string) *Base {
  file, err := os.OpenFile(fileName, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0664)
  if err != nil {
    log.Fatal(err)
  }

  b := Base{ fd: int(file.Fd()), page_count: 0 }
  
  err = file.Close()
  if err != nil {
    log.Fatal(err)
  }

  return &b;
}

func (b Base) GetFd() int { return b.fd }
func (b Base) GetPageCount() uint32 { return b.page_count }
func (b Base) IncrementPageCount() { atomic.AddUint32(&b.page_count, 1) }
