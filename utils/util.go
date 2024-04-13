package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

type ResponsePayload map[string]interface{}

func CreateFullPath(baseUrl string, p string) (*url.URL, error) {
	if baseUrl == "" || p == "" {
		return nil, errors.New("base path or postfix path can't be empty")
	}

	if baseUrl[len(baseUrl)-1] == '/' && p[0] == '/' {
		baseUrl = baseUrl[0 : len(baseUrl)-1]
	}

	result, err := url.Parse(baseUrl + p)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't parse url, error: %s, input: %s", err.Error(), baseUrl+p))
	}

	return result, nil
}

func Message(status bool, data ResponsePayload, message string) map[string]interface{} {
	payload := make(map[string]interface{})

	payload["success"] = status

	if message != "" {
		payload["message"] = message
	}

	for k, v := range data {
		payload[k] = v
	}

	return payload
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return
	}
}

// WriteResponse compiles and writes out response body and status code
func WriteResponse(w http.ResponseWriter, statusCode int, message string, data map[string]interface{}) {
	// success value in response
	var sBool bool

	switch statusCode {
	case http.StatusOK, http.StatusAccepted, http.StatusCreated, http.StatusFound:
		sBool = true
	default:
		sBool = false
	}

	w.WriteHeader(statusCode)
	Respond(w, Message(sBool, data, message))
}

func IntToTime(val interface{}) time.Time {
	kind := reflect.TypeOf(val).Kind()
	var tm time.Time
	switch kind {
	case reflect.String:
		i, err := strconv.ParseInt(fmt.Sprintf("%s", val), 10, 64)
		if err != nil {
			panic(err)
		}
		tm = time.Unix(i, 0)
	case reflect.Int64:
		tm = time.Unix(val.(int64), 0)
	case reflect.Int32:
		tm = time.Unix(int64(val.(int32)), 0)
	default:
		tm = time.Unix(int64(val.(int16)), 0)
	}

	return tm
}
