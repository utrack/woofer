package ihttp

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/utrack/woofer/domain"
	"github.com/utrack/woofer/lib/session"
	"github.com/utrack/woofer/service"
)

// Handler is an HTTP interface for the Woofer service.
type Handler struct {
	svc  *service.Woofer
	sess session.Storage
}

// NewHandler creates a new Handler using services provided.
func NewHandler(svc *service.Woofer, sess session.Storage) *Handler {
	return &Handler{svc: svc, sess: sess}
}

type tweetRequest struct {
	Text string `json:"text"`
}

type tweetResponse struct {
	TweetID uint64 `json:"tweet_id"`
}

type userCreateResponse struct {
	UserID domain.UserID `json:"user_id"`
}

// Tweet is a POST request containing tweetRequest.
// Returns tweetResponse.
func (h Handler) Tweet(w http.ResponseWriter, r *http.Request) {
	var req tweetRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logrus.Info(err)
		renderError(w, errors.Wrap(err, "error when parsing JSON body"), 400)
		return
	}
	id, err := h.svc.Tweet(r.Context(), req.Text)
	if err != nil {
		renderError(w, err, 500)
		return
	}
	json.NewEncoder(w).Encode(tweetResponse{TweetID: id})
}

// GetTweetPage is a GET request that has ?from URI param.
func (h Handler) GetTweetPage(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	fromUint, _ := strconv.ParseUint(from, 10, 0)
	tweets, err := h.svc.GetTweetPage(r.Context(), fromUint)
	if err != nil {
		renderError(w, err, 500)
	}
	json.NewEncoder(w).Encode(tweets)
}

// GetProfileTweets is a GET request that has ?from URI param and nickname path URI param.
func (h Handler) GetProfileTweets(w http.ResponseWriter, r *http.Request) {
	targetStr := chi.URLParam(r, "nickname")
	from := r.URL.Query().Get("from")
	fromUint, _ := strconv.ParseUint(from, 10, 0)
	tweets, err := h.svc.GetTweetsForProfile(r.Context(), targetStr, fromUint)
	if err != nil {
		renderError(w, err, 500)
	}
	json.NewEncoder(w).Encode(tweets)
}

// Subscribe is a GET request that has path URI param 'nickname'.
func (h Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	targetStr := chi.URLParam(r, "nickname")
	err := h.svc.Subscribe(r.Context(), targetStr)
	renderError(w, err, 500)
}

// Unsubscribe is a GET request that has path URI param 'nickname'.
func (h Handler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	targetStr := chi.URLParam(r, "nickname")
	err := h.svc.Unsubscribe(r.Context(), targetStr)
	renderError(w, err, 500)
}

// GetUser is a GET request that has path URI param 'nickname'.
func (h Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	targetStr := chi.URLParam(r, "nickname")
	user, err := h.svc.UserByNickname(r.Context(), targetStr)
	if err != nil {
		renderError(w, err, 500)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// Subscriptions is a GET request without any parameters.
func (h Handler) Subscriptions(w http.ResponseWriter, r *http.Request) {
	subs, err := h.svc.Subscriptions(r.Context())
	if err != nil {
		renderError(w, err, 500)
	}
	json.NewEncoder(w).Encode(subs)
}

// Subscribers is a GET request without any parameters.
func (h Handler) Subscribers(w http.ResponseWriter, r *http.Request) {
	subs, err := h.svc.Subscribers(r.Context())
	if err != nil {
		renderError(w, err, 500)
	}
	json.NewEncoder(w).Encode(subs)
}

// UserCreate is a POST request that should contain domain.UserWithPassword
// JSON in its body.
func (h Handler) UserCreate(w http.ResponseWriter, r *http.Request) {
	var req domain.UserWithPassword
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		renderError(w, errors.Wrap(err, "error when parsing JSON body"), 400)
		return
	}
	id, err := h.svc.UserCreate(r.Context(), req)
	if err != nil {
		renderError(w, err, 500)
		return
	}
	json.NewEncoder(w).Encode(userCreateResponse{UserID: id})
}

// Login is a POST form with username and password as form params.
func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		renderError(w, err, 400)
		return
	}

	user := r.FormValue("username")
	pass := r.FormValue("password")
	err = h.svc.CheckPassword(r.Context(), user, pass)
	if err != nil {
		renderError(w, err, 500)
		return
	}
	userObj, err := h.svc.UserByNickname(r.Context(), user)
	if err != nil {
		renderError(w, err, 500)
		return
	}

	expiration := time.Now().Add(14 * 24 * time.Hour)
	sessID, err := h.sess.SaveID(userObj.ID)
	if err != nil {
		renderError(w, err, 500)
		return
	}
	cookie := http.Cookie{Name: cookieSessID, Value: sessID, Expires: expiration}
	http.SetCookie(w, &cookie)

}
