package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"backend-assignment/model"
	"backend-assignment/store"
	"backend-assignment/util"
)

//이슈 상태값
var validStatuses = map[string]struct{}{
	"PENDING":     {},
	"IN_PROGRESS": {},
	"COMPLETED":   {},
	"CANCELLED":   {},
}

func getUserByID(uid uint) (*model.User, bool) {
	u, ok := store.Users[uid]
	return u, ok
}

func getIssueByID(iid uint) (*model.Issue, bool) {
	store.IssuesLock.Lock()
	defer store.IssuesLock.Unlock()
	iss, ok := store.Issues[iid]
	return iss, ok
}

func isValidStatus(status string) bool {
	_, ok := validStatuses[status]
	return ok
}

// 이슈 생성: title 필수, userId 있으면 IN_PROGRESS, 없으면 PENDING
func CreateIssue(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		UserID      *uint  `json:"userId"`
	}
	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteJSON(w, 400, util.ErrorResponse{"요청 데이터를 알 수 없습니다.", 400})
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		util.WriteJSON(w, 400, util.ErrorResponse{"title은 필수입니다.", 400})
		return
	}
	var user *model.User
	status := "PENDING"
	if req.UserID != nil {
		u, ok := getUserByID(*req.UserID)
		if !ok {
			util.WriteJSON(w, 400, util.ErrorResponse{"존재하지 않는 사용자입니다.", 400})
			return
		}
		user = u
		status = "IN_PROGRESS"
	}
	now := store.NowUTC()
	id := store.NextIssueID()
	issue := &model.Issue{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		User:        user,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	store.IssuesLock.Lock()
	store.Issues[id] = issue
	store.IssuesLock.Unlock()
	util.WriteJSON(w, 201, issue)
}

// 이슈 목록 조회: 전체 or status별 필터링
func ListIssues(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	store.IssuesLock.Lock()
	defer store.IssuesLock.Unlock()
	var out []model.Issue
	for _, iss := range store.Issues {
		if status == "" || iss.Status == status {
			tmp := *iss
			out = append(out, tmp)
		}
	}
	util.WriteJSON(w, 200, map[string]interface{}{"issues": out})
}

// 이슈 상세 조회 id로
func GetIssue(w http.ResponseWriter, r *http.Request) {
	idstr := strings.TrimPrefix(r.URL.Path, "/issue/")
	iid, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		util.WriteJSON(w, 400, util.ErrorResponse{"올바르지 않은 이슈 ID입니다.", 400})
		return
	}
	iss, ok := getIssueByID(uint(iid))
	if !ok {
		util.WriteJSON(w, 404, util.ErrorResponse{"이슈를 찾을 수 없습니다.", 404})
		return
	}
	util.WriteJSON(w, 200, iss)
}


// 이수 수정: title, description, status, userId 수정
func UpdateIssue(w http.ResponseWriter, r *http.Request) {
	idstr := strings.TrimPrefix(r.URL.Path, "/issue/")
	iid, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		util.WriteJSON(w, 400, util.ErrorResponse{"올바르지 않은 이슈 ID입니다.", 400})
		return
	}
	store.IssuesLock.Lock()
	iss, ok := store.Issues[uint(iid)]
	if !ok {
		store.IssuesLock.Unlock()
		util.WriteJSON(w, 404, util.ErrorResponse{"이슈를 찾을 수 없습니다.", 404})
		return
	}
	if iss.Status == "COMPLETED" || iss.Status == "CANCELLED" {
		store.IssuesLock.Unlock()
		util.WriteJSON(w, 400, util.ErrorResponse{"완료/취소된 이슈는 수정할 수 없습니다.", 400})
		return
	}
	type Req struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Status      *string `json:"status"`
		UserID      *uint   `json:"userId"`
	}
	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		store.IssuesLock.Unlock()
		util.WriteJSON(w, 400, util.ErrorResponse{"요청 데이터를 해석할 수 없습니다.", 400})
		return
	}

	if req.UserID != nil {
		if *req.UserID == 0 {
			iss.User = nil
			iss.Status = "PENDING"
		} else {
			u, ok := getUserByID(*req.UserID)
			if !ok {
				store.IssuesLock.Unlock()
				util.WriteJSON(w, 400, util.ErrorResponse{"존재하지 않는 사용자입니다.", 400})
				return
			}
			iss.User = u
			if iss.Status == "PENDING" && (req.Status == nil || *req.Status == "") {
				iss.Status = "IN_PROGRESS"
			}
		}
	}

	if req.Status != nil && *req.Status != "" {
		st := strings.ToUpper(*req.Status)
		if !isValidStatus(st) {
			store.IssuesLock.Unlock()
			util.WriteJSON(w, 400, util.ErrorResponse{"올바르지 않은 상태입니다.", 400})
			return
		}
		if iss.User == nil && st != "PENDING" && st != "CANCELLED" {
			store.IssuesLock.Unlock()
			util.WriteJSON(w, 400, util.ErrorResponse{"담당자 없는 이슈는 해당 상태로 변경할 수 없습니다.", 400})
			return
		}
		iss.Status = st
	}
	if req.Title != nil {
		iss.Title = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		iss.Description = *req.Description
	}
	iss.UpdatedAt = store.NowUTC()
	store.IssuesLock.Unlock()
	util.WriteJSON(w, 200, iss)
}
