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
const tokenExpTime = 15 * time.Second
const subExpTime = 30 * time.Second

var clientState = SigningUp
var serverState = server_running
var submissionState = sub_getting

func main() {
	// sign up
	acount, err := api.SignUp(`byr`)
	if err != nil {
		for err != nil {
			time.Sleep(10 * time.Second)
			fmt.Println(err)
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
							if err.ServerCrashed() { // 服务器崩溃
								serverState = server_crashed
							}
						}
					}

					clientState = TokenNotExpired
					hbTimer = time.After(tokenExpTime)

				case <-hbTimer: // token 更新
					tokenData, err = api.HeartBeat(tokenData)

					if err != nil {
						fmt.Println("HeartBeat error:", err)
						if err, ok := err.(*api.NmsError); ok {
							if err.ServerCrashed() { // 服务器崩溃
								serverState = server_crashed
							} else if err.TokenExpired() { // token 过期
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
						fmt.Println("GetSubmission error:", err)
						if err, ok := err.(*api.NmsError); ok {
							if err.ServerCrashed() { // 服务器崩溃
								serverState = server_crashed
							} else if err.TokenExpired() { // token 过期
								clientState = TokenExpired
								exp <- true
							}
						}
					}

					submissionState = sub_got

				case sub_got:
					_, err := api.SubmitCode(tokenData, submission)

					if err != nil {
						fmt.Println("SubmitCode error:", err)
						if err, ok := err.(*api.NmsError); ok {
							if err.ServerCrashed() { // 服务器崩溃
								serverState = server_crashed
							} else if err.TokenExpired() { // token 过期
								clientState = TokenExpired
								exp <- true
							} else if err.SubmissionExpired() { // 提交过期
								submissionState = sub_getting
								continue
							}
						}
					}
					submissionState = sub_getting
					time.Sleep(subExpTime)
				}
			}
		}
	}
}
