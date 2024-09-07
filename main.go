package main

import (
	"fmt"
	"netExperiment/api"
	"time"
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

	time.Sleep(2 * time.Second)

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
