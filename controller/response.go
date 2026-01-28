package controller

import (
	"errors"
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	codeError   = -1
	codeSuccess = 0
	codeFail    = 1
)

func ResponseOk() hero.Response {
	return hero.Response{
		Code: iris.StatusOK,
		Object: Response{
			Code: codeSuccess,
			Msg:  "success",
			Data: "",
		},
	}
}

func ResponseData(data interface{}) hero.Response {
	return hero.Response{
		Code: iris.StatusOK,
		Object: Response{
			Code: codeSuccess,
			Msg:  "success",
			Data: data,
		},
	}
}

func ResponseFail(data interface{}) hero.Response {
	var msg string
    switch v := data.(type) {
    case string:
        msg = v
    case error:
        msg = v.Error()
    default:
        msg = fmt.Sprintf("%v", v)
    }

	return hero.Response{
		Code: iris.StatusOK,
		Object: Response{
			Code: codeFail,
			Msg:  msg,
			Data: data,
		},
	}
}

func ResponseError(err error) hero.Response {
	return hero.Response{
		Code: iris.StatusInternalServerError,
		Err:  err,
	}
}

func ResponseErrorCode(code int,err error) hero.Response {
	return hero.Response{
		Code: code,
		Err:  err,
	}
}

func ResponseErrorStr(err string) hero.Response {
	return hero.Response{
		Code: iris.StatusInternalServerError,
		Err:  errors.New(err),
	}
}

func ResponseErrorParam(err error) hero.Response {
	return hero.Response{
		Code: iris.StatusBadRequest,
		Err:  err,
	}
}

func ResponseErrorUnauthorized(err error) hero.Response {
	return hero.Response{
		Code: iris.StatusUnauthorized,
		Err:  err,
	}
}
