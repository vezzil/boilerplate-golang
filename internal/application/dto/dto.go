package dto

import (
    "reflect"
)

type ResponseDto struct {
    Code  int         `json:"code"` 
    Msg   string      `json:"msg"` 
    Data  interface{} `json:"data,omitempty"` 
    Count int64       `json:"count"` 
}

type IncentiveResponseDto struct {
    Code       int         `json:"code"` 
    Msg        string      `json:"msg"` 
    Data       interface{} `json:"data,omitempty"` 
    Count      int64       `json:"count"` 
    TotalPrice float64     `json:"total_price"` 
}

func Success(data interface{}) *ResponseDto {
    if data == nil {
        data = []interface{}{}
    } else if reflect.TypeOf(data).Kind() == reflect.Slice && reflect.ValueOf(data).IsNil() {
        data = []interface{}{}
    }
    return &ResponseDto{
        Code:  0,
        Msg:   "SUCCESS",
        Data:  data,
        Count: 0,
    }
}

func SuccessMessage(message string, data interface{}) *ResponseDto {
    return &ResponseDto{
        Code:  0,
        Msg:   message,
        Data:  data,
        Count: 0,
    }
}

func SuccessCount(data interface{}, count int64) *ResponseDto {
    return &ResponseDto{
        Code:  0,
        Msg:   "SUCCESS",
        Data:  data,
        Count: count,
    }
}

func SuccessIncentiveCount(data interface{}, count int64, totalPrice float64) *IncentiveResponseDto {
    return &IncentiveResponseDto{
        Code:       0,
        Msg:        "SUCCESS",
        Data:       data,
        Count:      count,
        TotalPrice: totalPrice,
    }
}

func FailIncentive(msg string) *IncentiveResponseDto {
    return &IncentiveResponseDto{
        Code: 1,
        Msg:  msg,
    }
}

func Fail(msg string) *ResponseDto {
    return &ResponseDto{
        Code: 1,
        Msg:  msg,
    }
}

func FailCode(code int) *ResponseDto {
    return &ResponseDto{
        Code: code,
        Msg:  "No Permission",
    }
}

type NullDto struct{}