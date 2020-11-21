package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Contact はTolotの一件分の連絡先
type Contact struct {
	Code              *string      `json:"code"`
	CreateDate        *string      `json:"create_date"`
	ModifyDate        *string      `json:"modify_date"`
	Honorific         *string      `json:"honorific"`
	CompanyName       *string      `json:"company_name"`
	CompanyDepartment *string      `json:"company_department"`
	CompanyZip        *string      `json:"company_zip"`
	CompanyState      *string      `json:"company_state"`
	CompanyAddress1   *string      `json:"company_address1"`
	CompanyAddress2   *string      `json:"company_address2"`
	CompanyTel        *string      `json:"company_tel"`
	Memo              *string      `json:"memo"`
	JointNames        *[]JointName `json:"joint_names"`
	Name              *string      `json:"name"`
	Zip               *string      `json:"zip"`
	Country           *string      `json:"country"`
	State             *string      `json:"state"`
	City              *string      `json:"city"`
	Address1          *string      `json:"address1"`
	Address2          *string      `json:"address2"`
	Address3          *string      `json:"address3"`
	Area              *string      `json:"area"`
	NameLast          *string      `json:"name_last"`
	NameFirst         *string      `json:"name_first"`
	Kana              *string      `json:"kana"`
	Tel               *string      `json:"tel"`
	KanaLast          *string      `json:"kana_last"`
	KanaFirst         *string      `json:"kana_first"`
	Birthday          *string      `json:"birthday"`
	Email             *string      `json:"email"`
}

// JointName は連名
type JointName struct {
	Name      string `json:"name,omitempty"`
	NameLast  string `json:"name_last"`
	NameFirst string `json:"name_first"`
	Honorific string `json:"honorific"`
}

// Delete は連絡先を削除する
func (c *Contact) Delete(cookies []*http.Cookie) error {
	type Param struct {
		Code string `json:"code"`
	}

	param := new(Param)
	param.Code = *c.Code

	paramJson, err := json.Marshal(param)
	if err != nil {
		return err
	}

	url := "https://api.tolot.com/api/contact/delete"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(paramJson))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	for _, c := range cookies {
		req.AddCookie(c)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("status code is not correct: %d", res.StatusCode)
	}
	return nil
}
