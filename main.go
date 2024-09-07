package main

import (
	"fmt"
	"netExperiment/api"
	"time"
)

// 状态机
const (
	SigningUp       = iota // 注册中
	TokenExpired           // token 过期
	TokenNotExpired        // token 未过期
	TokenUpdating          // token 更新中
	GettingSub             // 获取内容
	Submitting             // 提交内容
)

// server state
const (
	server_crashed = iota
	server_running
)

func main() {
	acount, err := api.SignUp(`byr`)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(acount)

	token, err := api.Login(acount)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(token)
	fmt.Println(token.ValidTime)
	fmt.Println(time.Now())

	// time.Sleep(2 * time.Second)
	sub, err := api.GetSubmission(token)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(sub)

	ok, err := api.SubmitCode(token, sub)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(ok))
}
