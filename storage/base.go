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

func NewBaseFile(fileName string) (*Base, *os.File) {
  file, err := os.OpenFile(fileName, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0664)
  if err != nil {
    log.Fatal(err)
  }

  b := Base{ fd: int(file.Fd()), page_count: 0 }

  return &b, file;
}

func (b Base) CloseFile(file os.File) {
  if err := file.Close(); err != nil {
    log.Fatal(err)
  }
}

// inlines for encapsulation
func (b Base) GetFd() int { return b.fd }
func (b Base) GetPageCount() uint32 { return b.page_count }
func (b Base) IncrementPageCount() { atomic.AddUint32(&b.page_count, 1) }

// Given a page id and a page buffer, write content of PAGE_SIZE to the file pointed by the page_id
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
 
  // Ensure the page buffer has exactly PAGE_SIZE (4096) bytes
  if len(page) > essentials.PAGE_SIZE {
      return fmt.Errorf("page exceeds the maximum size of %d bytes", essentials.PAGE_SIZE)
  }
  
  paddedPage := make([]byte, essentials.PAGE_SIZE)
  copy(paddedPage, page) // Copy the contents of the original page

  // write PAGE_SIZE bytes into page buffer
  writtenBytes, err := syscall.Pwrite(int(pid.GetFileId()), paddedPage, int64(pid.GetPageNumber() * essentials.PAGE_SIZE))
  if err != nil || writtenBytes != essentials.PAGE_SIZE{
    return fmt.Errorf("write failed or incomplete: %w", err)
  }
    
  // final check to verify all data has been written to storage
  if err := syscall.Fsync(int(pid.GetFileId())); err != nil {
    return fmt.Errorf("data sync failed: %w", err)
  }

  return nil;
}
