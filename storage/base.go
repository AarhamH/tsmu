package storage

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"syscall"

	essentials "github.com/AarhamH/tsmu/essentials"
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

func (b Base) Flush (pid essentials.PageId, page []byte) error {
  // fail on invalid page id and nullptr
  if(page == nil) {
    return fmt.Errorf("page is nil")
  }

  if(!pid.IsPageIdValid()) {
    return fmt.Errorf("page id is invalid")
  } 

  var fileInfo syscall.Stat_t
  if err := syscall.Fstat(b.fd, &fileInfo); err != nil {
    return fmt.Errorf("failed to read basefile: %w", err)
  }
 
  // write PAGE_SIZE bytes into page buffer
  writtenBytes, err := syscall.Pwrite(int(pid.GetFileId()), page, int64(pid.GetPageNumber() * essentials.PAGE_SIZE))
  if err != nil || writtenBytes != essentials.PAGE_SIZE {
    return fmt.Errorf("write failed or incomplete: %w", err)
  }
    
  // final check to verify all data has been written to storage
  if err := syscall.Fsync(b.fd); err != nil {
    return fmt.Errorf("data sync failed: %w", err)
  }

  return nil;
}
