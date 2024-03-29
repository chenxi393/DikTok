package main

import "golang.org/x/crypto/bcrypt"

// BcryptHash 使用 bcrypt 进行加密，返回加密后的哈希值
func bcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值是否一致
func bcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
