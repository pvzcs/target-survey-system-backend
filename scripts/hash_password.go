package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: go run scripts/hash_password.go <密码>")
		fmt.Println("示例: go run scripts/hash_password.go mypassword123")
		os.Exit(1)
	}

	password := os.Args[1]

	// 使用 bcrypt 加密密码（与项目中的实现一致）
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("加密失败: %v\n", err)
		os.Exit(1)
	}

	hashedPassword := string(hashedBytes)
	fmt.Println("原始密码:", password)
	fmt.Println("加密密码:", hashedPassword)
	fmt.Println("\n你可以将加密后的密码直接插入数据库")
}
