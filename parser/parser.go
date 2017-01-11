package parser

import (
	"encoding/json"
	"fmt"
	"reflect"
	"serverskeleton/utils"
	"unicode"
	"unicode/utf8"
)

const (
	ErrorCode_OK = 0 + iota
	ErrorCode_JsonError
	ErrorCode_MethodNotFound
	ErrorCode_ParameterNotMatch
)

type Request struct {
	FuncName string        `json:"func_name"`
	Params   []interface{} `json:"params"`
}

type Response struct {
	FuncName  string        `json:"func_name"`
	Data      []interface{} `json:"data"`
	ErrorCode int           `json:"errorcode"`
}

type MethodInfo struct {
	Method reflect.Method
	Host   reflect.Value
	Idx    int
}

func RegisterMethod(MethodMap map[string]*MethodInfo, v interface{}) {

	reflectType := reflect.TypeOf(v)
	host := reflect.ValueOf(v)

	for i := 0; i < reflectType.NumMethod(); i++ {
		m := reflectType.Method(i)

		char, _ := utf8.DecodeRuneInString(m.Name)
		//非导出函数不注册
		if !unicode.IsUpper(char) {
			continue
		}

		MethodMap[m.Name] = &MethodInfo{Method: m, Host: host, Idx: i}
	}

}

func GenRequest(data []byte) (*Request, bool) {
	req := &Request{}
	err := json.Unmarshal(data, req)

	if nil != err {
		fmt.Println("json.Unmarshal err:", err)
		return nil, false
	}

	return req, true
}

func GenErrRespones(errorcode int) *Response {
	resp := &Response{}
	resp.FuncName = "Error"
	resp.ErrorCode = errorcode

	return resp
}

func Invoke(methodMap map[string]*MethodInfo, req *Request, defaultParams ...interface{}) *Response {

	methodInfo, found := methodMap[req.FuncName]

	if !found {
		return GenErrRespones(ErrorCode_MethodNotFound)
	}

	if len(req.Params) != methodInfo.Method.Type.NumIn()-1-len(defaultParams) {
		return GenErrRespones(ErrorCode_ParameterNotMatch)
	}

	req.Params = append(defaultParams, req.Params...)
	// fmt.Println(Params)
	paramsValue := make([]reflect.Value, 0, len(req.Params))

	//跳过 receiver
	for i := 1; i < methodInfo.Method.Type.NumIn(); i++ {
		inParaType := methodInfo.Method.Type.In(i)
		value, ok := utils.ConvertParamType(req.Params[i-1], inParaType)
		if !ok {
			return GenErrRespones(ErrorCode_ParameterNotMatch)
		}
		paramsValue = append(paramsValue, value)
	}

	data := &Response{}
	data.FuncName = methodInfo.Method.Name
	result := methodInfo.Host.Method(methodInfo.Idx).Call(paramsValue)
	for _, x := range result {
		data.Data = append(data.Data, x.Interface())
	}

	return data
}
