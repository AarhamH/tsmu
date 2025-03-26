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
  os.Mkdir("../tests/fixtures", os.ModePerm)
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

  t.Run(fmt.Sprintf("test: load page given valid buffer"), func(t *testing.T) {
    // Initialize basefile
    basefile, file:= basefile.NewBaseFile(filePath)
    inputString := strings.Repeat("A", 4096)
    os.WriteFile(filePath, []byte(inputString), 0664)
    
    fd := basefile.GetFd();
    
    pageId := essentials.NewPageId(fd, 0)
    outpufBuffer := make([]byte, essentials.PAGE_SIZE)
    if err := basefile.Load(*pageId, outpufBuffer); err != nil {
      t.Errorf("error occured when loading page: %s", err)
      t.Fail()
    }
    
    outputBuffString := strings.TrimSpace(string(outpufBuffer))
    if strings.Compare(outputBuffString, inputString) != 0 {
      t.Errorf("expecting output buffer: %s, got: %s", inputString, outputBuffString)
      t.Fail()
    }

    // Close
    os.WriteFile(filePath, []byte(""), 0664)
    basefile.CloseFile(*file)
    
  }) 

  t.Run(fmt.Sprintf("test: load page given invalid buffer"), func(t *testing.T) {
    // Initialize basefile
    basefile, file:= basefile.NewBaseFile(filePath)
    inputString := strings.Repeat("A", 4096)
    os.WriteFile(filePath, []byte(inputString), 0664)
    
    fd := basefile.GetFd();
    
    pageId := essentials.NewPageId(fd, 0)
    if err := basefile.Load(*pageId, nil); err == nil {
      t.Errorf("expecting error when loading to nil output buffer")
      t.Fail()
    }

    // Close
    os.WriteFile(filePath, []byte(""), 0664)
    basefile.CloseFile(*file)
  })

  t.Run(fmt.Sprintf("test: flush page and immediately load same page"), func(t *testing.T) {
    // Initialize basefile
    basefile, file:= basefile.NewBaseFile(filePath)
    fd := basefile.GetFd();
    
    repeatingText := strings.Repeat("A", 4096)
    inputBuffer := []byte(repeatingText)
    pageId := essentials.NewPageId(fd, 0)
    outpufBuffer := make([]byte, essentials.PAGE_SIZE)

    if err := basefile.Flush(*pageId, inputBuffer); err != nil {
      t.Errorf("error occured when flushing page: %s", err)
      t.Fail()
    }

    if err := basefile.Load(*pageId, outpufBuffer); err != nil {
      t.Errorf("error occured when loading page: %s", err)
      t.Fail()
    }
    
    inputBuffString := string(inputBuffer)
    outputBuffString := string(outpufBuffer)
    if strings.Compare(outputBuffString, inputBuffString) != 0 {
      t.Errorf("expecting output buffer: %s, got: %s", inputBuffString, outputBuffString)
      t.Fail()
    }

    // Close
    os.WriteFile(filePath, []byte(""), 0664)
    basefile.CloseFile(*file)
  }) 

  t.Run(fmt.Sprintf("test: create multiple pages in serial"), func(t *testing.T) {
    // Initialize basefile
    basefile, file:= basefile.NewBaseFile(filePath)
      
    kPages := 1000
    for i := 0; i< kPages; i++ {
      newPageId, err := basefile.Construct()

      if err != nil {
        t.Errorf("%s", err)
        t.Fail()
      }

      if newPageId.GetFileId() != uint32(basefile.GetFd()) {
        t.Errorf("expecting file id: %d, got %d", basefile.GetFd(), newPageId.GetFileId())
        t.Fail()
      }

      if newPageId.GetPageNumber() != uint32(i) {
        t.Errorf("expecting file id: %d, got %d", i, newPageId.GetPageNumber())
        t.Fail()
      }
    }
    basefile.CloseFile(*file)
  })

  t.Run(fmt.Sprintf("test: construct, load, flush, and load"), func(t *testing.T) {
    // Initialize basefile
    basefile, file:= basefile.NewBaseFile(filePath)
      
    kPages := 1000
    for i := 0; i< kPages; i++ {
      newPageId, err := basefile.Construct()

      if err != nil {
        t.Errorf("error occured when constructing %s", err)
        t.Fail()
      }
        
      outputBuffer := make([]byte, essentials.PAGE_SIZE)
      err = basefile.Load(*newPageId, outputBuffer)
      
      if err != nil {
        t.Errorf("error occured when loading page %s", err)
        t.Fail()
      }

      for j:=0; j<essentials.PAGE_SIZE/4; j++ {
        outputBuffer[j] = byte(j) 
      }

      err = basefile.Flush(*newPageId, outputBuffer)
      if err != nil {
        t.Errorf("error occured when flushing page %s", err)
        t.Fail()
      }
      
      secondOutputBuffer := make([]byte, essentials.PAGE_SIZE)
      err = basefile.Load(*newPageId, secondOutputBuffer)
      if err != nil {
        t.Errorf("error occured when loading page for the second time %s", err)
        t.Fail()
      }

      for j:=0; j<essentials.PAGE_SIZE/4; j++ {
        secondOutputBuffer[j] = outputBuffer[j] 
      }
    }
    basefile.CloseFile(*file)
  })

}

