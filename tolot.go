package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Tolot struct {
	Cookies  []*http.Cookie
	Contacts []*Contact
}

// Loginはtolot.comにログインし、API実行に必要なcookieを取得する
func Login(mailAddress string, password string) (tolot *Tolot, err error) {
	type Param struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}

	param := new(Param)
	param.UserName = mailAddress
	param.Password = password

	paramJson, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	url := "https://api.tolot.com/api/member/login"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(paramJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code is not correct: %d", res.StatusCode)
	}
	return &Tolot{Cookies: res.Cookies()}, nil
}

type ListResponse struct {
	Result struct {
		Contacts []Contact `json:"contacts"`
	} `json:"result"`
}

// List は連絡先の一覧を取得する
func (tolot *Tolot) List() (contacts *[]Contact, err error) {
	url := "https://api.tolot.com/api/contact/list?page=1&paginate_by=50"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	for _, c := range tolot.Cookies {
		req.AddCookie(c)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code is not correct: %d", res.StatusCode)
	}
	body := new(bytes.Buffer)
	_, err = io.Copy(body, res.Body)
	if err != nil {
		return nil, err
	}

	var listResponse ListResponse
	err = json.Unmarshal(body.Bytes(), &listResponse)
	if err != nil {
		return nil, err
	}

	return &listResponse.Result.Contacts, nil
}

// Add は連絡帳を1件インポートする
// フォーマットがContactではなく、ImportContactであることに注意する
func (tolot *Tolot) Import(contact *ImportContact) (err error) {
	paramJson, err := json.Marshal(contact)
	if err != nil {
		return err
	}

	url := "https://api.tolot.com/api/contact/create"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(paramJson))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	for _, c := range tolot.Cookies {
		req.AddCookie(c)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Status code is not correct: %d", res.StatusCode))
	}
	return nil
}
