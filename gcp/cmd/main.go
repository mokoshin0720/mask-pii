package main

import (
	"fmt"

	"github.com/mokoshin0720/mask-pii/gcp"
	"github.com/mokoshin0720/mask-pii/gcp/config"
)

func main() {
	if err := config.Init(); err != nil {
		fmt.Printf("設定の初期化に失敗しました: %v\n", err)
		return
	}

	result, err := gcp.Mask("my-project-id", "My SSN is 111222333", []string{"US_SOCIAL_SECURITY_NUMBER"}, "+", 6)
	if err != nil {
		fmt.Printf("エラーが発生しました: %v\n", err)
		return
	}
	fmt.Printf("マスク処理結果: %s\n", result)
}
