package main

import (
	"fmt"
	"netExperiment/api"
	"time"
)

// client state
const (
	SigningUp       = iota // 注册中
	TokenExpired           // token 过期
	TokenNotExpired        // token 未过期
)

// server state
const (
	server_crashed = iota
	server_running
)

// submission state
const (
	sub_getting = iota
	sub_got
)

const crashedTime = 10 * time.Second
const tokenExpTime = 5 * time.Second
const subExpTime = 30 * time.Second

var clientState = SigningUp
var serverState = server_running
var submissionState = sub_getting

func main() {
	// sign up
	acount, err := api.SignUp(`byr`)
	if err != nil {
		for err != nil {
			fmt.Println("Sign up error:", err)
			time.Sleep(10 * time.Second)
			acount, err = api.SignUp(`byr`)
		}
	}

	clientState = TokenExpired
	var tokenData api.TokenData
	var submission api.Submission
	var hbTimer <-chan time.Time
	exp := make(chan bool, 1)
	exp <- true

	// fmt.Println(acount)

	// login
	go func() {
		for {
			switch serverState {
			case server_running:
				select {
				case <-exp: // token 过期
					tokenData, err = api.Login(acount)

					if err != nil {
						fmt.Println("Login error:", err)
						if err, ok := err.(*api.NmsError); ok {
							if err.ServerCrashed() { // 服务器崩溃 计划再次登录
								serverState = server_crashed
							}
						}
						exp <- true
					} else { // 服务器正常 token 获取成功 计划下次更新
						clientState = TokenNotExpired
						hbTimer = time.After(tokenExpTime)
					}

				case <-hbTimer: // token 更新
					tokenData, err = api.HeartBeat(tokenData)

					if err != nil {
						fmt.Println("HeartBeat error:", err)
						if err, ok := err.(*api.NmsError); ok {
							if err.ServerCrashed() { // 服务器崩溃
								serverState = server_crashed
							} else if err.TokenExpired() { // token 过期 计划登录
								clientState = TokenExpired
								exp <- true
							}
						}
					}

					hbTimer = time.After(tokenExpTime)
				}
			case server_crashed:
				time.Sleep(crashedTime)
				serverState = server_running
			}
		}
	}()

	// submit
	for {
		switch serverState {
		case server_running:
			if clientState == TokenNotExpired {
				switch submissionState {
				case sub_getting:
					submission, err = api.GetSubmission(tokenData)

					if err != nil {
						fmt.Println("GetSub error:", err)
						if err, ok := err.(*api.NmsError); ok {
							if err.ServerCrashed() { // 服务器崩溃
								serverState = server_crashed
							} else if err.TokenExpired() { // token 过期 计划登录
								clientState = TokenExpired
							}
						}
					} else { // 服务器正常 内容获取成功 计划提交
						submissionState = sub_got
					}

				case sub_got:
					// time.Sleep(100 * time.Second)
					_, err := api.SubmitCode(tokenData, submission)

					if err != nil {
						fmt.Println("SubmitCode error:", err)
						if err, ok := err.(*api.NmsError); ok {
							if err.ServerCrashed() { // 服务器崩溃
								serverState = server_crashed
							} else if err.TokenExpired() { // token 过期
								clientState = TokenExpired
							} else if err.SubmissionExpired() { // 提交过期 计划重新获取
								submissionState = sub_getting
							}
						}
					} else { // 服务器正常 提交成功 计划下次获取
						submissionState = sub_getting
						time.Sleep(subExpTime)
					}
				}
			}
		}
	}
}
