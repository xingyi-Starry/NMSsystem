package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type jsonData interface {
	resolveJWT(j []byte)
}

type Account struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type TokenData struct {
	Token     string `json:"token"`
	validTime time.Time
}

type Submission struct {
	Code string `json:"code"`
}

func (a *Account) resolveJWT(j []byte) {
	err := json.Unmarshal(j, &a)
	if err != nil {
		fmt.Println("err in resolveJson:", err)
		return
	}
	fmt.Println("Password:", a.Password)
	fmt.Println("Username:", a.Username)
}

func (t *TokenData) resolveJWT(j []byte) {
	err := json.Unmarshal(j, &t)
	if err != nil {
		fmt.Println("err in resolveJson:", err)
		return
	}
	//divide
	parts := strings.Split(t.Token, ".")
	if len(parts) != 3 {
		fmt.Println("Invalid JWT token")
		return
	}

	// head
	// payload
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Println("err in resolveJson:", err)
		return
	}

	type Vt struct {
		Name string `json:"name"`
		Exp  int64  `json:"exp"`
	}
	var vt Vt
	err = json.Unmarshal(payload, &vt)
	if err != nil {
		fmt.Println("err in resolveJson:", err)
		return
	}
	t.validTime = time.Unix(vt.Exp, 0)
}

func (s *Submission) resolveJWT(j []byte) {
	err := json.Unmarshal(j, &s)
	if err != nil {
		fmt.Println("err in resolveJson:", err)
		return
	}
	fmt.Println("Code:", s.Code)
}

func SignUp(usrName string) Account {
	client := &http.Client{}
	var data = strings.NewReader(`username=` + usrName)
	req, err := http.NewRequest("POST", "http://localhost:1323/signup", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("signin: %s\n", bodyText)
	var a Account
	a.resolveJWT(bodyText)
	return a
}

func Login(a Account) TokenData {
	client := &http.Client{}
	var data = strings.NewReader(`username=` + a.Username + `&password=` + a.Password)
	req, err := http.NewRequest("POST", "http://localhost:1323/login", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("login: %s\n", bodyText)
	var t TokenData
	t.resolveJWT(bodyText)
	return t
}
