package goinside

import (
	"errors"
	"io"
	"net/url"
)

var (
	errInvalidIDorPW = errors.New("invalid ID or PW")
)

// GuestSession 구조체는 유동닉 세션을 표현합니다.
type GuestSession struct {
	id   string
	pw   string
	conn *Connection
}

// Guest 함수는 유동닉 세션을 반환합니다.
func Guest(id, pw string) (gs *GuestSession, err error) {
	if len(id) == 0 || len(pw) == 0 {
		err = errInvalidIDorPW
		return
	}
	gs = &GuestSession{id: id, pw: pw, conn: &Connection{}}
	return
}

func (gs *GuestSession) Connection() *Connection {
	if gs.conn == nil {
		gs.conn = &Connection{}
	}
	return gs.conn
}

// Write 메소드는 쓰기 가능한 객체를 전달받아 작성 요청을 보냅니다.
func (gs *GuestSession) Write(wa writable) error {
	return wa.write(gs)
}

func (gs *GuestSession) articleWriteForm(ad *ArticleDraft) (io.Reader, string) {
	return multipartForm(map[string]string{
		"app_id":   AppID,
		"mode":     "write",
		"name":     gs.id,
		"password": gs.pw,
		"id":       ad.GallID,
		"subject":  ad.Subject,
		"content":  ad.Content,
	}, ad.Images...)
}

func (gs *GuestSession) commentWriteForm(cd *CommentDraft) (io.Reader, string) {
	return makeForm(map[string]string{
		"app_id":       AppID,
		"comment_nick": gs.id,
		"comment_pw":   gs.pw,
		"id":           cd.Target.Gall.ID,
		"no":           cd.Target.Number,
		"comment_memo": cd.Content,
		"mode":         "comment_nonmember",
	}), defaultContentType
}

// Delete 메소드는 삭제 가능한 객체를 전달받아 삭제 요청을 보냅니다.
func (gs *GuestSession) Delete(da deletable) error {
	return da.delete(gs)
}

func (gs *GuestSession) articleDeleteForm(a *Article) (io.Reader, string) {
	return makeForm(map[string]string{
		"app_id":   AppID,
		"mode":     "board_del",
		"write_pw": gs.pw,
		"id":       a.Gall.ID,
		"no":       a.Number,
	}), defaultContentType
}

func (gs *GuestSession) commentDeleteForm(c *Comment) (io.Reader, string) {
	return makeForm(map[string]string{
		"app_id":     AppID,
		"comment_pw": gs.pw,
		"id":         c.Parents.Gall.ID,
		"no":         c.Parents.Number,
		"mode":       "comment_del",
		"comment_no": c.Number,
	}), defaultContentType
}

// ThumbsUp 메소드는 해당 글에 추천 요청을 보냅니다.
func (gs *GuestSession) ThumbsUp(a *Article) error {
	return a.thumbsUp(gs)
}

// ThumbsDown 메소드는 해당 글에 비추천 요청을 보냅니다.
func (gs *GuestSession) ThumbsDown(a *Article) error {
	return a.thumbsDown(gs)
}

func (gs *GuestSession) actionForm(a *Article) (io.Reader, string) {
	return makeForm(map[string]string{
		"app_id": AppID,
		"id":     a.Gall.ID,
		"no":     a.Number,
	}), nonCharsetContentType
}

// Report 메소드는 해당 글에 메모와 함께 신고 요청을 보냅니다.
func (gs *GuestSession) Report(a *Article, memo string) error {
	return a.report(gs, memo)
}

func (gs *GuestSession) reportForm(URL, memo string) (io.Reader, string) {
	_Must := func(s string, err error) string {
		if err != nil {
			return ""
		}
		return s
	}
	return makeForm(map[string]string{
		"name":     _Must(url.QueryUnescape(gs.id)),
		"password": _Must(url.QueryUnescape(gs.pw)),
		"choice":   "4",
		"memo":     _Must(url.QueryUnescape(memo)),
		"no":       articleNumber(URL),
		"id":       gallID(URL),
		"app_id":   AppID,
	}), nonCharsetContentType
}
