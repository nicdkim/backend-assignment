package store

import (
	"sync"
	"time"
	"backend-assignment/model"
)

var (
	issueSeq   uint = 1
	Issues          = make(map[uint]*model.Issue)
	IssuesLock sync.Mutex

	Users = map[uint]*model.User{
		1: {ID: 1, Name: "김개발"},
		2: {ID: 2, Name: "이디자인"},
		3: {ID: 3, Name: "박기획"},
	}
)

func NextIssueID() uint {
	IssuesLock.Lock()
	defer IssuesLock.Unlock()
	id := issueSeq
	issueSeq++
	return id
}

func NowUTC() time.Time {
	return time.Now().UTC()
}
