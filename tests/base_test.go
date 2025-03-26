package tests

import (
	"errors"
	"fmt"
	"os"
	"testing"
  basefile "github.com/AarhamH/tsmu/storage"
  essentials "github.com/AarhamH/tsmu/essentials"
)

const newlyCreatedFilePath = "../tests/fixtures/new_file.txt"
const filePath = "../tests/fixtures/testfile_1.txt"

func TestBaseFile(t *testing.T) {
  t.Run(fmt.Sprintf("test: instantiate basefile with non-existent fixtures file"), func(t *testing.T) {
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
    basefile.CloseFile(*file)
  }) 

  t.Run(fmt.Sprintf("test: flush page to valid file"), func(t *testing.T) {
    basefile,file := basefile.NewBaseFile(newlyCreatedFilePath)

    // assert that the filePath has been created
    if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
      t.Logf("expecting file path %s to be created", newlyCreatedFilePath)
      t.Fail()
    }
    
    fd := basefile.GetFd();
    
    pageId := essentials.NewPageId(fd, 0)
    buff := []byte("This is sample text")
    if err := basefile.Flush(*pageId, buff); err != nil {
      t.Errorf("error occured when flushing page %e", err)
      t.Fail()
    }
    content, err := os.ReadFile(newlyCreatedFilePath)
    if err != nil {
      t.Errorf("error occured when reading file: %e", err)
      t.Fail()
    }

    if string(buff) != string(content) {
      t.Errorf("expecting the content: %s, got: %s", string(buff), string(content))
      t.Fail()
    }
    basefile.CloseFile(*file)
    os.Remove(newlyCreatedFilePath)
  }) 
}

