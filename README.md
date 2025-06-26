#실행방법

## go 모듈 초기화
go mod init backend-assignment

## 서버 실행
go run main.go

## API 테스트 방법 윈도우 기준
### 이슈 생성
curl -X POST http://localhost:8080/issue -H "Content-Type: application/json" -d "{\"title\":\"버그 수정 필요\", \"description\":\"로그인 페이지에서 오류 발생\", \"userId\":1}"

### 이슈 목록 전체 조회
curl http://localhost:8080/issues

### 상태별 이슈 목록 조회
curl http://localhost:8080/issues?status=IN_PROGRESS

### 이슈 상세 조회
curl http://localhost:8080/issue/1

### 이슈 수정
curl -X PATCH http://localhost:8080/issue/1 -H "Content-Type: application/json" -d "{\"status\":\"COMPLETED\"}"

### 사용자가 존재하지 않을때
curl -X POST http://localhost:8080/issue -H "Content-Type: application/json" -d "{\"title\":\"에러테스트\", \"userId\":999}"
