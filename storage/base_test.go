package storage

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestBaseFile(t *testing.T) {
  t.Run(fmt.Sprintf("test: instantiate basefile with non-existent fixtures file"), func(t *testing.T) {
    filePath := "../storage/fixtures/NON_EXISTENT.txt"
    basefile := NewBaseFile(filePath) 
    
    if basefile == nil {
      t.Logf("basefile is nil; not instantiated properly")
      t.Fail()
    }
  
    if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
      t.Logf("expecting file path %s to be created", filePath)
      t.Fail()
    }

    os.Remove(filePath)
  }) 
  t.Run(fmt.Sprintf("test: instantiate basefile with existing fixtures file"), func(t *testing.T) {
    filePath := "../storage/fixtures/testfile_1.txt"
    basefile := NewBaseFile(filePath) 

    // assert that the filePath already exists
    if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
      t.Logf("expecting file path %s to exist", filePath)
      t.Fail()
    }

    fd := basefile.GetFd();
    page_count := basefile.GetPageCount()
    
    if(page_count != 0) {
      t.Logf("expecting page_count = 0, got %d", page_count)
      t.Fail()
    }
    
    if(fd == -1) {
      t.Logf("expecting fd != -1, got %d", fd)
      t.Fail()
    }
  }) 
}

