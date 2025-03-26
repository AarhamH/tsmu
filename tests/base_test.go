package tests

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	essentials "github.com/AarhamH/tsmu/essentials"
	basefile "github.com/AarhamH/tsmu/storage"
)

const filePath = "../tests/fixtures/testfile_1.txt"

func TestBaseFile(t *testing.T) {
  t.Run(fmt.Sprintf("test: instantiate basefile with non-existent fixtures file"), func(t *testing.T) {
    newlyCreatedFilePath := "../tests/fixtures/new_file.txt"
    
    // Initialize basefile
    basefile, file := basefile.NewBaseFile(newlyCreatedFilePath) 
    
    if basefile == nil {
      t.Logf("basefile is nil; not instantiated properly")
      t.Fail()
    }
  
    if _, err := os.Stat(newlyCreatedFilePath); errors.Is(err, os.ErrNotExist) {
      t.Logf("expecting file path %s to be created", filePath)
      t.Fail()
    }

    // Close
    basefile.CloseFile(*file)
    os.Remove(newlyCreatedFilePath)
  })

  t.Run(fmt.Sprintf("test: instantiate basefile with existing fixtures file"), func(t *testing.T) {
    basefile, file := basefile.NewBaseFile(filePath) 

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

    // Close
    basefile.CloseFile(*file)
  }) 

  t.Run(fmt.Sprintf("test: flush page to valid file"), func(t *testing.T) {
    // Initialize basefile
    basefile,file := basefile.NewBaseFile(filePath)
    
    fd := basefile.GetFd();
    
    pageId := essentials.NewPageId(fd, 0)
    buff := []byte("This is sample text")
    if err := basefile.Flush(*pageId, buff); err != nil {
      t.Errorf("error occured when flushing page %e", err)
      t.Fail()
    }
    
    // read contents of file and compare with buffer
    content, err := os.ReadFile(filePath)
    if err != nil {
      t.Errorf("error occured when reading file: %e", err)
      t.Fail()
    }

    if len(content) != essentials.PAGE_SIZE {
      t.Errorf("expecting the content size: %d, got: %d",essentials.PAGE_SIZE, len(content))
      t.Fail()
    }
    
    trimmedContent := strings.TrimSpace(string(content))
    if strings.Compare(string(buff), trimmedContent) == 0 {
      t.Errorf("expecting the content: %s, got: %s", string(buff), trimmedContent)
      t.Fail()
    }

    // Close
    basefile.CloseFile(*file)
  })

  t.Run(fmt.Sprintf("test: flush page given invalid buffer"), func(t *testing.T) {
    // Initialize basefile
    basefile,file := basefile.NewBaseFile(filePath)
    
    fd := basefile.GetFd();
    
    pageId := essentials.NewPageId(fd, 0)
    if err := basefile.Flush(*pageId, nil); err == nil {
      t.Errorf("expecting error for nil page buffer")
      t.Fail()
    }

    // Close
    basefile.CloseFile(*file)
  }) 
}

