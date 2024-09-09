package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

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

func (a *Account) resolveJWT(j []byte) error {
	err := json.Unmarshal(j, &a)
	if err != nil {
		return err
	}
	return nil
}

func (t *TokenData) resolveJWT(j []byte) error {
	err := json.Unmarshal(j, &t)
	if err != nil {
		return err
	}

	//divide
	parts := strings.Split(t.Token, ".")
	if len(parts) != 3 {
		fmt.Print("Invalid JWT token:", t.Token, string(j))
		return nil
	}

	// head
	// payload
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return err
	}

	type Vt struct {
		Name string `json:"name"`
		Exp  int64  `json:"exp"`
	}
	var vt Vt
	err = json.Unmarshal(payload, &vt)
	if err != nil {
		return err
	}
	t.ValidTime = time.Unix(vt.Exp, 0)
	return nil
}

func (s *Submission) resolveJWT(j []byte) error {
	err := json.Unmarshal(j, &s)
	if err != nil {
		return err
	}
	// fmt.Println("Code:", s.Code)
	return nil
}

func (s *SubResp) resolveJWT(j []byte) error {
	err := json.Unmarshal(j, &s)
	if err != nil {
		// fmt.Println("err in resolveJson:", err)
		return err
	}
	// fmt.Println("Message:", s.Message)
	return nil
}

func SignUp(usrName string) (Account, error) {
	bodyText, err := HttpRequest("POST", "signup", `username=`+usrName, "")
	if err != nil {
		return Account{}, err
	}
	fmt.Printf("signup: %s", bodyText)

	var a Account
	err = a.resolveJWT(bodyText)
	err = PendingError(a, err)
	if err != nil {
		return Account{}, err
	}
	fmt.Println("Username:", a.Username)
	fmt.Println("Password:", a.Password)
	return a, nil
}

func Login(a Account) (TokenData, error) {
	bodyText, err := HttpRequest("POST", "login", `username=`+a.Username+`&password=`+a.Password, "")
	if err != nil {
		return TokenData{}, err
	}

	var t TokenData
	err = t.resolveJWT(bodyText)
	err = PendingError(t, err)
	if err != nil {
		return TokenData{}, err
	}
	fmt.Println("Login token:", t.Token)
	return t, nil
}

func HeartBeat(t TokenData) (TokenData, error) {
	bodyText, err := HttpRequest("GET", "api/heartbeat", "", t.Token)
	if err != nil {
		return t, err
	}
	// fmt.Printf("heartbeat: %s\n", bodyText)

	var t_new TokenData
	err = t_new.resolveJWT(bodyText)
	err = PendingError(t_new, err)
	// if t_new.Token == "" {
	// 	err = NewNmsError("token expired", TokenExpErr)
	// }
	if err != nil {
		return t, err
	}
	fmt.Println("Heartbeat token:", t_new.Token)
	return t_new, nil
}

func GetSubmission(t TokenData) (Submission, error) {
	bodyText, err := HttpRequest("GET", "api/info", "", t.Token)
	if err != nil {
		return Submission{}, err
	}
	// fmt.Printf("submission: %s\n", bodyText)

	var s Submission
	err = s.resolveJWT(bodyText)
	err = PendingError(s, err)
	if err != nil {
		return Submission{}, err
	}
	fmt.Println("getSub:", s.Code)
	return s, nil
}

func SubmitCode(t TokenData, s Submission) ([]byte, error) {
	bodyText, err := HttpRequest("POST", "api/validate", `code=`+s.Code, t.Token)
	if err != nil {
		return []byte{}, err
	}

	// submission expired
	if string(bodyText) == "\"NOPE\"\n" {
		err = NewNmsError("err in SubmitCode: submission expired", SubExpErr)
		return bodyText, err
	}

	// token expired
	var subResp SubResp
	err = subResp.resolveJWT(bodyText)
	err = PendingError(subResp, err)
	if err != nil {
		fmt.Println("err in SubmitCode:", err)
	}

	// fmt.Printf("submit: %s\n", string(bodyText))
	fmt.Println("submit:", string(bodyText))
	return bodyText, nil
}

func GenToken(t TokenData) (TokenData, error) {
	parts := strings.Split(t.Token, ".")
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return t, err
	}
	type Vt struct {
		Name string `json:"name"`
		Exp  int64  `json:"exp"`
	}
	var vt Vt
	err = json.Unmarshal(payload, &vt)
	if err != nil {
		return t, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": vt.Name,
		"exp":      time.Now().Add(time.Second * 20).Unix(),
	})
	tokenString, err := token.SignedString([]byte("Hello_new_Byrs_1234123412341234"))
	if err != nil {
		return t, err
	}
	t.Token = tokenString
	fmt.Println("genToken:", t.Token)
	return t, nil
}
