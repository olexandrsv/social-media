package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"social-media/internal/common"
	"social-media/internal/common/app"
	"social-media/internal/common/app/config"
	"social-media/internal/users/endpoint"
	"testing"
)

type mockEndpoints struct {
	createUser        func(context.Context, interface{}) (interface{}, error)
	login             func(context.Context, interface{}) (interface{}, error)
	getUser           func(context.Context, interface{}) (interface{}, error)
	updateUser        func(context.Context, interface{}) (interface{}, error)
	getLoginsByInfo   func(context.Context, interface{}) (interface{}, error)
	followUser        func(context.Context, interface{}) (interface{}, error)
	getFollowedLogins func(context.Context, interface{}) (interface{}, error)
}

func (e mockEndpoints) CreateUser(ctx context.Context, request interface{}) (interface{}, error) {
	return e.createUser(ctx, request)
}

func (e mockEndpoints) Login(ctx context.Context, request interface{}) (interface{}, error) {
	if e.login == nil {
		return nil, nil
	}
	return e.login(ctx, request)
}

func (e mockEndpoints) GetUser(ctx context.Context, request interface{}) (interface{}, error) {
	if e.getUser == nil {
		return nil, nil
	}
	return e.getUser(ctx, request)
}

func (e mockEndpoints) UpdateUser(ctx context.Context, request interface{}) (interface{}, error) {
	if e.updateUser == nil {
		return nil, nil
	}
	return e.updateUser(ctx, request)
}

func (e mockEndpoints) GetLoginsByInfo(ctx context.Context, request interface{}) (interface{}, error) {
	if e.getLoginsByInfo == nil {
		return nil, nil
	}
	return e.getLoginsByInfo(ctx, request)
}

func (e mockEndpoints) FollowUser(ctx context.Context, request interface{}) (interface{}, error) {
	if e.followUser == nil {
		return nil, nil
	}
	return e.followUser(ctx, request)
}

func (e mockEndpoints) GetFollowedLogins(ctx context.Context, request interface{}) (interface{}, error) {
	if e.getFollowedLogins == nil {
		return nil, nil
	}
	return e.getFollowedLogins(ctx, request)
}

var mockErr = common.NewError(500, "mock error")

func TestCreateUser(t *testing.T) {
	app.InitMock(config.AppConfig{}, t)
	createBobReq := endpoint.CreateUserReq{Login: "bob", Name: "Bob", Surname: "Arnum", Password: "bob"}
	createBobResp := endpoint.AuthResp{Token: "bob"}

	data := []struct {
		e           mockEndpoints
		req         endpoint.CreateUserReq
		resp        endpoint.AuthResp
		err         error
		excludeForm bool
	}{
		{
			e: mockEndpoints{
				createUser: func(ctx context.Context, request interface{}) (interface{}, error) {
					req, ok := request.(endpoint.CreateUserReq)
					if !ok {
						t.Errorf("can't cast to endpoint.CreateUserReq: %+v", req)
					}
					if createBobReq.Login != req.Login || createBobReq.Name != req.Name ||
						createBobReq.Surname != req.Surname || createBobReq.Password != req.Password {
						t.Errorf("wrong request %+v, expected: %+v", req, createBobReq)
					}

					return createBobResp, nil
				},
			},
			req:  createBobReq,
			resp: createBobResp,
		},
		{
			e: mockEndpoints{
				createUser: func(ctx context.Context, request interface{}) (interface{}, error) {
					return nil, mockErr
				},
			},
			req: createBobReq,
			err: mockErr,
		},
		{
			excludeForm: true,
			err:         common.ErrInvalidData,
		},
	}

	for _, d := range data {
		srv := newMockServer(d.e)
		defer srv.Close()

		r := NewRequest("POST", srv.URL+"/register")
		if !d.excludeForm {
			r.WithForm(map[string]string{
				"login":       d.req.Login,
				"first_name":  d.req.Name,
				"second_name": d.req.Surname,
				"passw1":      d.req.Password,
			})
		}
		resp, _, err := r.Execute()
		if err != nil {
			t.Error(err)
		}

		correctResp, err := getCorrectResp(d.resp, d.err)
		if err != nil {
			t.Error(err)
		}

		if correctResp != resp {
			t.Errorf("wrong response %s, expected: %s", resp, correctResp)
		}

	}
}

func TestLogin(t *testing.T) {
	app.InitMock(config.AppConfig{}, t)
	loginBobReq := endpoint.LoginReq{Login: "bob", Password: "123"}
	loginBobResp := endpoint.AuthResp{Token: "bob"}

	data := []struct {
		e           mockEndpoints
		req         endpoint.LoginReq
		resp        endpoint.AuthResp
		err         error
		excludeForm bool
	}{
		{
			e: mockEndpoints{
				login: func(ctx context.Context, request interface{}) (interface{}, error) {
					req, ok := request.(endpoint.LoginReq)
					if !ok {
						t.Errorf("can't cast to endpoint.LoginReq: %+v", req)
					}
					if loginBobReq.Login != req.Login || loginBobReq.Password != req.Password {
						t.Errorf("wrong request %+v, expected: %+v", req, loginBobReq)
					}

					return loginBobResp, nil
				},
			},
			req:  loginBobReq,
			resp: loginBobResp,
		},
		{
			e: mockEndpoints{
				login: func(ctx context.Context, request interface{}) (interface{}, error) {
					return nil, mockErr
				},
			},
			req: loginBobReq,
			err: mockErr,
		},
		{
			excludeForm: true,
			err:         common.ErrInvalidData,
		},
	}

	for _, d := range data {
		srv := newMockServer(d.e)
		defer srv.Close()

		r := NewRequest("POST", srv.URL+"/login")
		if !d.excludeForm {
			r.WithForm(map[string]string{
				"login":    d.req.Login,
				"password": d.req.Password,
			})
		}
		resp, _, err := r.Execute()
		if err != nil {
			t.Error(err)
		}

		correctResp, err := getCorrectResp(d.resp, d.err)
		if err != nil {
			t.Error(err)
		}

		if correctResp != resp {
			t.Errorf("wrong response %s, expected: %s", resp, correctResp)
		}

	}
}

func TestGetUser(t *testing.T) {
	app.InitMock(config.AppConfig{}, t)
	getBobReq := endpoint.GetUserReq{Login: "bob"}
	getBobResp := endpoint.GetUserResp{Login: "bob", Name: "Bob", Surname: "Smith", Bio: "student", Interests: "hockey"}
	getBobErr := common.ErrInvalidData

	data := []struct {
		e    mockEndpoints
		req  endpoint.GetUserReq
		resp endpoint.GetUserResp
		err  error
	}{
		{
			e: mockEndpoints{
				getUser: func(ctx context.Context, request interface{}) (interface{}, error) {
					req, ok := request.(endpoint.GetUserReq)
					if !ok {
						t.Errorf("can't cast to endpoint.GetUserReq: %+v", req)
					}
					if req.Login != getBobReq.Login {
						t.Errorf("wrong request %+v, expected: %+v", req, getBobReq)
					}
					return getBobResp, nil
				},
			},
			req:  getBobReq,
			resp: getBobResp,
		},
		{
			e: mockEndpoints{
				getUser: func(ctx context.Context, request interface{}) (interface{}, error) {
					return nil, getBobErr
				},
			},
			req: getBobReq,
			err: getBobErr,
		},
	}

	for _, d := range data {
		srv := newMockServer(d.e)
		defer srv.Close()

		resp, _, err := NewRequest("GET", srv.URL+"/info/"+d.req.Login).Execute()
		if err != nil {
			t.Error(err)
		}

		correctResp, err := getCorrectResp(d.resp, d.err)
		if err != nil {
			t.Error(err)
		}

		if correctResp != resp {
			t.Errorf("wrong response %s, expected: %s", resp, correctResp)
		}
	}
}

func TestUpdateUser(t *testing.T) {
	app.InitMock(config.AppConfig{}, t)
	updateBobReq := endpoint.UpdateUserReq{Token: "bob", Name: "Bob", Surname: "Smith", Bio: "developer", Interests: "skiing"}
	data := []struct {
		e             mockEndpoints
		req           endpoint.UpdateUserReq
		resp          any
		err           error
		excludeForm   bool
		excludeCookie bool
	}{
		{
			e: mockEndpoints{
				updateUser: func(ctx context.Context, request interface{}) (interface{}, error) {
					req, ok := request.(endpoint.UpdateUserReq)
					if !ok {
						t.Errorf("can't cast to endpoint.UpdateUserReq: %+v", req)
					}

					if req.Token != updateBobReq.Token || req.Name != updateBobReq.Name ||
						req.Surname != updateBobReq.Surname || req.Bio != updateBobReq.Bio ||
						req.Interests != updateBobReq.Interests {
						t.Errorf("wrong request: %+v, expected: %+v", req, updateBobReq)
					}
					return nil, nil
				},
			},
			req:  updateBobReq,
			resp: nil,
		},
		{
			excludeForm: true,
			err:         common.ErrInvalidData,
		},
		{
			excludeCookie: true,
			err:           common.ErrNoToken,
		},
	}

	for _, d := range data {
		srv := newMockServer(d.e)
		defer srv.Close()

		r := NewRequest("POST", srv.URL+"/info")
		if !d.excludeForm {
			r.WithForm(map[string]string{
				"first_name":  d.req.Name,
				"second_name": d.req.Surname,
				"bio":         d.req.Bio,
				"interests":   d.req.Interests,
			})
		}
		if !d.excludeCookie {
			r.WithCookie(map[string]string{
				"token": d.req.Token,
			})
		}
		resp, _, err := r.Execute()
		if err != nil {
			t.Error(err)
		}

		correctResp, err := getCorrectResp(d.resp, d.err)
		if err != nil {
			t.Error(err)
		}

		if correctResp != resp {
			t.Errorf("wrong response %s, expected: %s", resp, correctResp)
		}
	}
}

func TestGetLoginsByInfo(t *testing.T) {
	app.InitMock(config.AppConfig{}, t)
	getBobsReq := endpoint.GetLoginsByInfoReq{Info: "Bob"}
	getBobsResp := endpoint.LoginsResp{Logins: []string{"bob01", "bob02"}}
	getBobsErr := common.ErrInternal

	data := []struct {
		e           mockEndpoints
		req         endpoint.GetLoginsByInfoReq
		resp        endpoint.LoginsResp
		err         error
		excludeForm bool
	}{
		{
			e: mockEndpoints{
				getLoginsByInfo: func(ctx context.Context, request interface{}) (interface{}, error) {
					req, ok := request.(endpoint.GetLoginsByInfoReq)
					if !ok {
						t.Errorf("can't cast %+v to endpoint.GetLoginsByInfoReq", req)
					}

					if req.Info != getBobsReq.Info {
						t.Errorf("wrong request: %+v, expected: %+v", req, getBobsReq)
					}

					return getBobsResp, nil
				},
			},
			req:  getBobsReq,
			resp: getBobsResp,
		},
		{
			e: mockEndpoints{
				getLoginsByInfo: func(ctx context.Context, request interface{}) (interface{}, error) {
					return nil, getBobsErr
				},
			},
			err: getBobsErr,
		},
		{
			excludeForm: true,
			err:         common.ErrInvalidData,
		},
	}

	for _, d := range data {
		s := NewHTTPServer(d.e)
		srv := httptest.NewServer(s.router)

		r := NewRequest("POST", srv.URL+"/filter")
		if !d.excludeForm {
			r.WithForm(map[string]string{
				"info": d.req.Info,
			})
		}
		resp, _, err := r.Execute()
		if err != nil {
			t.Error(err)
		}

		correctResp, err := getCorrectResp(d.resp, d.err)
		if err != nil {
			t.Error(err)
		}

		if correctResp != resp {
			t.Errorf("wrong response %s, expected: %s", resp, correctResp)
		}
	}
}

func TestFollowUser(t *testing.T) {
	app.InitMock(config.AppConfig{}, t)
	followBobReq := endpoint.FollowUserReq{Token: "ben", Login: "bob"}
	
	data := []struct {
		e             mockEndpoints
		req           endpoint.FollowUserReq
		resp          any
		err           error
		excludeCookie bool
	}{
		{
			e: mockEndpoints{
				followUser: func(ctx context.Context, request interface{}) (interface{}, error) {
					req, ok := request.(endpoint.FollowUserReq)
					if !ok {
						t.Errorf("can't cast %+v to endpoint.FollowUserReq", req)
					}

					if req.Token != followBobReq.Token || req.Login != followBobReq.Login {
						t.Errorf("wrong request: %+v, expected: %+v", req, followBobReq)
					}

					return nil, nil
				},
			},
			req:  followBobReq,
			resp: nil,
		},
		{
			e: mockEndpoints{
				followUser: func(ctx context.Context, request interface{}) (interface{}, error) {
					return nil, mockErr
				},
			},
			req: followBobReq,
			err: mockErr,
		},
		{
			req:           followBobReq,
			excludeCookie: true,
			err:           common.ErrNoToken,
		},
	}

	for _, d := range data {
		srv := newMockServer(d.e)
		defer srv.Close()

		r := NewRequest("POST", srv.URL+"/follow/"+d.req.Login)
		if !d.excludeCookie {
			r.WithCookie(map[string]string{
				"token": d.req.Token,
			})
		}
		resp, _, err := r.Execute()
		if err != nil {
			t.Error(err)
		}

		correctResp, err := getCorrectResp(d.resp, d.err)
		if err != nil {
			t.Error(err)
		}

		if correctResp != resp {
			t.Errorf("wrong response: %s, expected: %s", resp, correctResp)
		}
	}
}

func TestGetFollowedLogins(t *testing.T) {
	app.InitMock(config.AppConfig{}, t)
	getBobFollowersReq := endpoint.TokenReq{Token: "bob"}
	getBobFollowersResp := endpoint.LoginsResp{Logins: []string{"ben", "bill"}}
	data := []struct {
		e             mockEndpoints
		req           endpoint.TokenReq
		resp          endpoint.LoginsResp
		err           error
		excludeCookie bool
	}{
		{
			e: mockEndpoints{
				getFollowedLogins: func(ctx context.Context, request interface{}) (interface{}, error) {
					req, ok := request.(endpoint.TokenReq)
					if !ok {
						t.Errorf("can't cast request %+v to endpoint.TokenReq", req)
					}

					if req.Token != getBobFollowersReq.Token {
						t.Errorf("wrong request: %+v, expected: %+v", req, getBobFollowersReq)
					}
					return getBobFollowersResp, nil
				},
			},
			req:  getBobFollowersReq,
			resp: getBobFollowersResp,
		},
		{
			e: mockEndpoints{
				getFollowedLogins: func(ctx context.Context, request interface{}) (interface{}, error) {
					return nil, mockErr
				},
			},
			err: mockErr,
		},
		{
			excludeCookie: true,
			err: common.ErrNoToken,
		},
	}

	for _, d := range data {
		srv := newMockServer(d.e)
		defer srv.Close()

		r := NewRequest("GET", srv.URL+"/follow")
		if !d.excludeCookie {
			r.WithCookie(map[string]string{
				"token": d.req.Token,
			})
		}
		resp, _, err := r.Execute()
		if err != nil {
			t.Error(err)
		}

		correctResp, err := getCorrectResp(d.resp, d.err)
		if err != nil {
			t.Error(err)
		}

		if correctResp != resp {
			t.Errorf("wrong response: %s, expected: %s", resp, correctResp)
		}
	}
}

func getCorrectResp(v any, err error) (string, error) {
	if err != nil {
		return err.Error(), nil
	}
	if v == nil {
		return "", nil
	}
	resp, err := structToJsonResp(v)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func structToJsonResp(v any) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", nil
	}
	return string(bytes) + "\n", nil
}

func newMockServer(e endpoint.Endpoints) *httptest.Server {
	s := NewHTTPServer(e)
	return httptest.NewServer(s.router)
}

func createForm(data map[string]string) (io.Reader, string, error) {
	body := new(bytes.Buffer)
	mWriter := multipart.NewWriter(body)
	defer mWriter.Close()
	for k, v := range data {
		if err := mWriter.WriteField(k, v); err != nil {
			return nil, "", err
		}
	}
	return body, mWriter.FormDataContentType(), nil
}
