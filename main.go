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
	TokenUpdating          // token 更新中
	GettingSub             // 获取内容
	Submitting             // 提交内容
)

// server state
const (
	server_crashed = iota
	server_running
)

var clientState = SigningUp
var serverState = server_running

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
	var hbTimer <-chan time.Time
	exp := make(chan bool, 1)
	exp <- true

	// fmt.Println(acount)

	for {
		switch serverState {
		case server_running:
			select {
			case <-exp: // token 过期
				tokenData, err = api.Login(acount)
				if err != nil {
					fmt.Println("Login error:", err)
				}

				clientState = TokenNotExpired
				hbTimer = time.After(10 * time.Second)
			case <-hbTimer: // token 更新
				tokenData, err = api.HeartBeat(tokenData)
				if err != nil {
					fmt.Println("HeartBeat error:", err)
				}

				hbTimer = time.After(10 * time.Second)
			}
		case server_crashed:

		}
	}
}
