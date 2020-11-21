package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	// settings.yamlから設定を読み込み
	viper.SetConfigName("settings") // 設定ファイル名を拡張子抜きで指定する
	viper.AddConfigPath(".")        // 現在のワーキングディレクトリを探索することもできる
	err := viper.ReadInConfig()     // 設定ファイルを探索して読み取る
	if err != nil {                 // 設定ファイルの読み取りエラー対応
		panic(fmt.Errorf("設定ファイル読み込みエラー: %s", err))
	}
	var userID = viper.GetString("userID")
	var password = viper.GetString("password")
	var inFileName = viper.GetString("addressFileName")

	contacts, err := Open(inFileName)
	if err != nil {
		panic(err)
	}

	// ログイン
	tolot, err := Login(userID, password)
	if err != nil {
		panic(err)
	}
	// 連絡帳一覧を取得
	Contacts, err := tolot.List()
	if err != nil {
		panic(err)
	}
	// 削除
	for _, c := range *Contacts {
		fmt.Print(*c.NameLast + " " + *c.NameFirst + "...")
		err := c.Delete(tolot.Cookies)
		if err != nil {
			panic(err)
		}
		fmt.Println("\tdeleted.")
	}

	// 追加
	for _, c := range *contacts {
		fmt.Print(c.Name)
		err = tolot.Import(c)
		if err != nil {
			panic(fmt.Errorf("error on add contact: %s, name: %s", err.Error(), c.Name))
		}
		fmt.Println("\tadded.")
	}
}
