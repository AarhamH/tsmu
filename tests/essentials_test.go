package tests

import (
	"fmt"
	"testing"
  "math/rand/v2"
  essentials "github.com/AarhamH/tsmu/essentials"
)

func TestPageId(t *testing.T) {
  fileId := rand.IntN(100)
  pageNum := rand.IntN(100)

  t.Run(fmt.Sprintf("test: pageId generated with random fileId and pageNum"), func(t *testing.T) {
    pid := essentials.NewPageId(fileId, uint64(pageNum))
    isValid := pid.IsPageIdValid() 
    if !isValid {
      t.Logf("Expected IsPageIdValid = true, got %t", isValid)
      t.Fail()
    }
  })

  t.Run(fmt.Sprintf("test: pageId calculates accurate pageNum"), func(t *testing.T) {
    pid := essentials.NewPageId(fileId, uint64(pageNum))
    isValid := pid.IsPageIdValid() 
    if !isValid {
      t.Logf("Expected IsPageIdValid = true, got %t", isValid)
      t.Fail()
    }

    if pageNum != int(pid.GetPageNumber()) {
      t.Logf("Expected fileId = %d, got %d", pageNum, int(pid.GetPageNumber()))
      t.Fail()
    }
  })

  t.Run(fmt.Sprintf("test: pageId calculates accurate fileId"), func(t *testing.T) {
    pid := essentials.NewPageId(fileId, uint64(pageNum))
    isValid := pid.IsPageIdValid() 
    if !isValid {
      t.Logf("Expected IsPageIdValid = true, got %t", isValid)
      t.Fail()
    }

    if uint32(fileId) != pid.GetFileId() {
      t.Logf("Expected fileId = %d, got %d", fileId, int(pid.GetFileId()))
      t.Fail()
    }
  })
}

func TestRecordId(t *testing.T) {
  fileId := rand.IntN(100)
  pageNum := rand.IntN(100)
  slotId := rand.IntN(100)
  t.Run(fmt.Sprintf("test: recordId is valid"), func(t *testing.T) {
    pid := essentials.NewPageId(fileId, uint64(pageNum))
    isValid := pid.IsPageIdValid() 
    if !isValid {
      t.Logf("Expected IsPageIdValid = true, got %t", isValid)
      t.Fail()
    }
    rid := essentials.NewRecordId(pid, uint64(slotId))
    isValid = rid.IsRecordIdValid()
    if !isValid {
      t.Logf("Expected IsRecordIdValid =, got %t", isValid)
      t.Fail()
    }
  })

  t.Run(fmt.Sprintf("test: recordId calculates accurate slotId"), func(t *testing.T) {
    pid := essentials.NewPageId(fileId, uint64(pageNum))
    isValid := pid.IsPageIdValid() 
    if !isValid {
      t.Logf("Expected IsPageIdValid = true, got %t", isValid)
      t.Fail()
    }
    rid := essentials.NewRecordId(pid, uint64(slotId))
    isValid = rid.IsRecordIdValid()
    if !isValid {
      t.Logf("Expected IsRecordIdValid =, got %t", isValid)
      t.Fail()
    }

    if uint32(slotId) != rid.GetSlotId() {
      t.Logf("Expected slotId = %d, got %d", slotId, int(rid.GetSlotId()))
      t.Fail()
    }
  })
}


