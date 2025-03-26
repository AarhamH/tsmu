package storage

import (
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"syscall"

	essentials "github.com/AarhamH/tsmu/essentials"
)

type Base struct {
  fd int // fd represents the file descriptor after opening a file  
  pageCount uint32 // pageCount is treated as an atomic integer
  baseLatch sync.Mutex 
}

func NewBaseFile(fileName string) (*Base, *os.File) {
  file, err := os.OpenFile(fileName, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0664)
  if err != nil {
    log.Fatal(err)
  }

  b := Base{ fd: int(file.Fd()), pageCount: 0 }

  return &b, file;
}

func (b *Base) CloseFile(file os.File) {
  if err := file.Close(); err != nil {
    log.Fatal(err)
  }
}

// inlines for encapsulation
func (b *Base) GetFd() int { return b.fd }
func (b *Base) GetPageCount() uint32 { return b.pageCount }
func (b *Base) IncrementPageCount() { atomic.AddUint32(&b.pageCount, 1) }

// Given a page id and a page buffer, write content of PAGE_SIZE to the file pointed by the page_id
func (b *Base) Flush (pid essentials.PageId, page []byte) error {
  // check for valid parameters before flushing
  if err := isBufferAndIdValid(pid, page, pid.GetFileId()); err != nil {
    return err
  }

  // Ensure the page buffer has exactly PAGE_SIZE (4096) bytes
  if len(page) > essentials.PAGE_SIZE {
      return fmt.Errorf("page exceeds the maximum size of %d bytes", essentials.PAGE_SIZE)
  }
  
  paddedPage := make([]byte, essentials.PAGE_SIZE)
  copy(paddedPage, page) // Copy the contents of the original page

  // write PAGE_SIZE bytes into page buffer
  writtenBytes, err := syscall.Pwrite(int(pid.GetFileId()), paddedPage, int64(pid.GetPageNumber() * essentials.PAGE_SIZE))
  if err != nil {
    return err
  }
  if writtenBytes != essentials.PAGE_SIZE {
    return fmt.Errorf("write incompete: only wrote %d bytes", writtenBytes)
  }

  // final check to verify all data has been written to storage
  if err := syscall.Fsync(int(pid.GetFileId())); err != nil {
    return err
  }

  return nil;
}

// given page id and output buffer, write contents pointed to by page_id to buffer
func (b *Base) Load(pid essentials.PageId, buffer[]byte) error {
  // check for valid parameters before loading
  if err := isBufferAndIdValid(pid, buffer, pid.GetFileId()); err != nil {
    return err
  }
  
  readBytes, err := syscall.Pread(b.fd, buffer, (int64(b.GetPageCount()) * essentials.PAGE_SIZE))
  if err != nil {
    return err
  }

  if readBytes < essentials.PAGE_SIZE && readBytes > 0 {
    return fmt.Errorf("read incomplete")
  }

  return nil
}

func (b *Base) Construct() (*essentials.PageId, error) {
  // treat entire code block as critical section
  b.baseLatch.Lock()
  defer b.baseLatch.Unlock()

  newFileSize := essentials.PAGE_SIZE * (b.pageCount + 1)
  
  // tries to check if file was not adjusted, if so return error
  if err := syscall.Ftruncate(b.fd, int64(newFileSize)); err != nil {
    return nil, err
  }

  b.IncrementPageCount()

  // build new page id
  newPID := essentials.NewPageId(b.fd, uint64(b.pageCount - 1))

  buff := make([]byte, essentials.PAGE_SIZE)
  if err := b.Flush(*newPID, buff); err != nil {
    fmt.Println("Error2 here")
    return nil, err
  }

  return newPID, nil
}

// helpers
func isBufferAndIdValid(pid essentials.PageId, page []byte, fd uint32) error {
  if(page == nil) {
    return fmt.Errorf("page is nil")
  }

  if(!pid.IsPageIdValid()) {
    return fmt.Errorf("page id is invalid")
  }

  var fileInfo syscall.Stat_t
  if err := syscall.Fstat(int(fd), &fileInfo); err != nil {
    return fmt.Errorf("failed to read basefile: %w", err)
  }

  return nil;
}
