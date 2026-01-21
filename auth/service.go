package auth

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	terrors "github.com/ueebee/tachibanashi/errors"
	"github.com/ueebee/tachibanashi/model"
)

const (
	authPath          = "auth/"
	defaultJSONFormat = "5"

	clmAuthLoginRequest  = "CLMAuthLoginRequest"
	clmAuthLogoutRequest = "CLMAuthLogoutRequest"
)

type Client interface {
	DoJSON(ctx context.Context, method, path string, req, resp any) error
	TokenStore() TokenStore
	SetVirtualURLs(urls VirtualURLs)
	VirtualURLs() VirtualURLs
}

type Service struct {
	client Client
}

func NewService(client Client) *Service {
	return &Service{client: client}
}

type Credentials struct {
	LoginID  string
	Password string
}

type VirtualURLs struct {
	Request string `json:"sUrlRequest"`
	Master  string `json:"sUrlMaster"`
	Price   string `json:"sUrlPrice"`
	Event   string `json:"sUrlEvent"`
}

type LoginResponse struct {
	model.CommonResponse
	VirtualURLs
}

type loginRequest struct {
	model.CommonParams
	CLMID    string `json:"sCLMID"`
	UserID   string `json:"sUserId"`
	Password string `json:"sPassword"`
}

func (r *loginRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type logoutRequest struct {
	model.CommonParams
	CLMID string `json:"sCLMID"`
}

func (r *logoutRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

func (s *Service) Login(ctx context.Context, creds Credentials) (*LoginResponse, error) {
	if creds.LoginID == "" {
		return nil, &terrors.ValidationError{Field: "login_id", Reason: "required"}
	}
	if creds.Password == "" {
		return nil, &terrors.ValidationError{Field: "password", Reason: "required"}
	}

	if store := s.client.TokenStore(); store != nil {
		store.Reset()
	}

	req := loginRequest{
		CommonParams: model.CommonParams{JsonOfmt: defaultJSONFormat},
		CLMID:        clmAuthLoginRequest,
		UserID:       creds.LoginID,
		Password:     creds.Password,
	}

	var resp LoginResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, authPath, &req, &resp); err != nil {
		return nil, err
	}

	if !resp.VirtualURLs.isZero() {
		s.client.SetVirtualURLs(resp.VirtualURLs)
	}
	if resp.PNo != "" {
		if v, err := parseInt64(resp.PNo); err == nil {
			if store := s.client.TokenStore(); store != nil {
				store.Set(v)
			}
		}
	}

	return &resp, nil
}

func (s *Service) Logout(ctx context.Context) error {
	req := logoutRequest{
		CommonParams: model.CommonParams{JsonOfmt: defaultJSONFormat},
		CLMID:        clmAuthLogoutRequest,
	}

	if err := s.client.DoJSON(ctx, http.MethodGet, authPath, &req, nil); err != nil {
		return err
	}

	if store := s.client.TokenStore(); store != nil {
		store.Reset()
	}
	s.client.SetVirtualURLs(VirtualURLs{})

	return nil
}

func (s *Service) VirtualURL(ctx context.Context) (*VirtualURLs, error) {
	_ = ctx
	urls := s.client.VirtualURLs()
	if urls.isZero() {
		return nil, errors.New("tachibanashi: virtual URL not set")
	}
	return &urls, nil
}

func (v VirtualURLs) isZero() bool {
	return v.Request == "" && v.Master == "" && v.Price == "" && v.Event == ""
}

func parseInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
