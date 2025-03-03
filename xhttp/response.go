package xhttp

import "github.com/go-hao/zero/xerrors"

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func newResponse(v interface{}) Response {
	var resp Response
	switch data := v.(type) {
	case *xerrors.Error:
		resp.Code = data.Code
		resp.Msg = data.Msg
	case xerrors.Error:
		resp.Code = data.Code
		resp.Msg = data.Msg
	case error:
		resp.Code = CodeErr
		resp.Msg = data.Error()
	default:
		resp.Code = CodeOk
		resp.Msg = MsgOk
		resp.Data = v
	}

	return resp
}
