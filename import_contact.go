package main

import (
	"encoding/csv"
	"io"
	"os"
	"strings"
)

// ImportContact はインポートする連絡先
// インポート（画面表示）と実際の連絡先の肩が異なることに注意する
type ImportContact struct {
	Name       string            `json:"name"` // 氏名は半角スペースで区切る
	Kana       string            `json:"kana"` // 氏名は半角スペースで区切る
	Honorific  string            `json:"honorific"`
	Zip        string            `json:"zip"`
	State      string            `json:"state"`
	Address1   string            `json:"address1"`
	Address2   string            `json:"address2"`
	Tel        string            `json:"tel"`
	Birthday   string            `json:"birthday"`
	Email      string            `json:"email"`
	JointNames []ImportJointName `json:"joint_names,omitempty"`
}

// ImportContact はインポートする連名
type ImportJointName struct {
	NameLast  string `json:"name_last"`
	NameFirst string `json:"name_first"`
	Honorific string `json:"honorific"`
}

// ImportContact はインポートする連絡帳
type ImportContacts []*ImportContact

// 連絡帳のCSVファイルを開き、ImportContacts形式を返す
func Open(inFileName string) (*ImportContacts, error) {
	// アドレス帳を開く
	// 性	名	連名1	連名2	連名3	郵便番号	都道府県	市町村	それ以降	備考
	file, err := os.Open(inFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true // ダブルクオートを厳密にチェックしない

	var contacts ImportContacts
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		var contact ImportContact
		contact.Name = record[0] + " " + record[1]
		var jointNames []ImportJointName
		for _, jN := range record[2:5] {
			if jN != "" {
				var jointName ImportJointName
				jointName.NameFirst = jN
				jointNames = append(jointNames, jointName)
			}
		}
		contact.JointNames = jointNames
		zip := strings.Replace(record[5], "-", "", -1)
		contact.Zip = zip
		contact.State = record[6]
		contact.Address1 = record[7]
		contact.Address2 = record[8]
		contacts = append(contacts, &contact)
	}
	return &contacts, nil
}
