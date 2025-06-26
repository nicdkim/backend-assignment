package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"backend-assignment/model"
	"backend-assignment/store"
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

