package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// type jsonData interface {
// 	resolveJWT(j []byte)
// }

type Account struct {
	Password string `json:"password"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

type TokenData struct {
	Token     string `json:"token"`
	ValidTime time.Time
	Message   string `json:"message"`
}

type Submission struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SubResp struct {
	Message string `json:"message"`
}

func (a *Account) resolveJWT(j []byte) bool {
	err := json.Unmarshal(j, &a)
	if err != nil {
		fmt.Println("err in resolveJson:", err)
		return false
	}
	if a.Message != "" {
		fmt.Println("Message:", a.Message)
		return false
	}
	fmt.Println("Password:", a.Password)
	fmt.Println("Username:", a.Username)
	return true
}

func (t *TokenData) resolveJWT(j []byte) bool {
	err := json.Unmarshal(j, &t)
	if err != nil {
		fmt.Println("err in resolveJson:", err)
		return false
	}

	//divide
	parts := strings.Split(t.Token, ".")
	if len(parts) != 3 {
		fmt.Println("Invalid JWT token:", t.Token)
		return false
	}

	// head
	// payload
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Println("err in resolveJson:", err)
		return false
	}

	type Vt struct {
		Name string `json:"name"`
		Exp  int64  `json:"exp"`
	}
	var vt Vt
	err = json.Unmarshal(payload, &vt)
	if err != nil {
		fmt.Println("err in resolveJson:", err)
		return false
	}
	t.ValidTime = time.Unix(vt.Exp, 0)
	return true
}

func (s *Submission) resolveJWT(j []byte) bool {
	err := json.Unmarshal(j, &s)
	if err != nil {
		fmt.Println("err in resolveJson:", err)
		return false
	}
	// fmt.Println("Code:", s.Code)
	return true
}

func (s *SubResp) resolveJWT(j []byte) bool {
	err := json.Unmarshal(j, &s)
	if err != nil {
		// fmt.Println("err in resolveJson:", err)
		return false
	}
	// fmt.Println("Message:", s.Message)
	return true
}

func SignUp(usrName string) (Account, error) {
	bodyText, err := HttpRequest("POST", "signup", `username=`+usrName, "")
	if err != nil {
		err = fmt.Errorf("err in SignUp: %w", err)
		return Account{}, err
	}
	fmt.Printf("signup: %s\n", bodyText)

	var a Account
	if !a.resolveJWT(bodyText) {
		err = fmt.Errorf("err in SignUp: " + a.Message)
		return Account{}, err
	}
	return a, nil
}

func Login(a Account) (TokenData, error) {
	bodyText, err := HttpRequest("POST", "login", `username=`+a.Username+`&password=`+a.Password, "")
	if err != nil {
		err = fmt.Errorf("err in Login: %w", err)
		return TokenData{}, err
	}

	var t TokenData
	if !t.resolveJWT(bodyText) {
		err = fmt.Errorf("err in Login: " + t.Message)
		if t.Message == "Bad Gateway" {
			err = NewNmsError("err in Login: server crashed", ServerCrashErr)
		}
		return TokenData{}, err
	}
	if t.Message != "" {
		err = NewNmsError("err in Login: token expired", TokenExpErr)
		return TokenData{}, err
	}
	fmt.Println("Token:", t.Token)
	fmt.Println("expTime:", t.ValidTime)
	fmt.Println("nowTime:", time.Now())
	fmt.Println()
	return t, nil
}

func HeartBeat(t TokenData) (TokenData, error) {
	bodyText, err := HttpRequest("GET", "api/heartbeat", "", t.Token)
	if err != nil {
		err = fmt.Errorf("err in HeartBeat: %w", err)
		return TokenData{}, err
	}
	// fmt.Printf("heartbeat: %s\n", bodyText)

	var t_new TokenData
	if !t_new.resolveJWT(bodyText) {
		err = fmt.Errorf("err in HeartBeat: " + t_new.Message)
		if t.Message == "Bad Gateway" {
			err = NewNmsError("err in Login: server crashed", ServerCrashErr)
		}
		return TokenData{}, err
	}
	if t.Message != "" {
		err = NewNmsError("err in HeartBeat: token expired", TokenExpErr)
		return TokenData{}, err
	}
	fmt.Println("Token:", t_new.Token)
	fmt.Println("expTime:", t_new.ValidTime)
	fmt.Println("nowTime:", time.Now())
	fmt.Println()
	return t_new, nil
}

func GetSubmission(t TokenData) (Submission, error) {
	bodyText, err := HttpRequest("GET", "api/info", "", t.Token)
	if err != nil {
		err = fmt.Errorf("err in GetSubmission: %w", err)
		return Submission{}, err
	}
	// fmt.Printf("submission: %s\n", bodyText)

	var s Submission
	if !s.resolveJWT(bodyText) {
		err = fmt.Errorf("err in GetSubmission: " + s.Message)
		return Submission{}, err
	}
	if s.Message != "" {
		err = NewNmsError("err in GetSubmission: token expired", TokenExpErr)
		return Submission{}, err
	}
	fmt.Println("getSub:", s.Code)
	return s, nil
}

func SubmitCode(t TokenData, s Submission) ([]byte, error) {
	bodyText, err := HttpRequest("POST", "api/validate", `code=`+s.Code, t.Token)
	if err != nil {
		err = fmt.Errorf("err in SubmitCode: %w", err)
		return []byte{}, err
	}

	// submission expired
	if string(bodyText) == "\"NOPE\"\n" {
		err = NewNmsError("err in SubmitCode: submission expired", SubExpErr)
		return bodyText, err
	}

	// token expired
	var subResp SubResp
	if subResp.resolveJWT(bodyText) {
		if subResp.Message != "" {
			fmt.Println("Message:", subResp.Message)
			err = NewNmsError("err in SubmitCode: token expired", TokenExpErr)
			return []byte{}, err
		}
	}

	// fmt.Printf("submit: %s\n", string(bodyText))
	fmt.Println("submit:", string(bodyText))
	return bodyText, nil
}
