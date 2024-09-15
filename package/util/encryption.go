package util

import "golang.org/x/crypto/bcrypt"

// BcryptHash 使用 bcrypt 进行加密，返回加密后的哈希值
func BcryptHash(password string) string {
	// 实际上就是一个哈希计算
	// 应该每个密码使用不同的盐（需要写入DB）: 增加安全性，防止彩虹表攻击
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值是否一致
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
