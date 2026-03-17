// =============================================================================
// Encrypt - 加密和哈希工具
// =============================================================================
// 提供常用的加密和哈希功能，包括 MD5 哈希和 bcrypt 密码加密。
//
// 功能特性:
//   - MD5 哈希：快速生成数据摘要
//   - bcrypt 密码加密：安全的密码存储
//   - bcrypt 密码验证：验证密码是否正确
//
// 使用场景:
//   - 用户密码加密存储
//   - 用户登录密码验证
//   - 数据完整性校验（MD5）
//   - 文件指纹生成（MD5）
//
// 设计思路:
//   - MD5 用于快速哈希，不用于密码存储（不安全）
//   - bcrypt 用于密码加密，自动加盐，防止彩虹表攻击
//   - 提供简单易用的 API，隐藏底层实现细节
//
// 项目中的应用:
//   - User 服务：用户注册时加密密码
//   - User 服务：用户登录时验证密码
//   - 文件服务：生成文件 MD5 校验和
//
// 第三方依赖:
//   - golang.org/x/crypto/bcrypt: bcrypt 密码加密库
//
// 安全建议:
//   - 密码存储必须使用 bcrypt，不要使用 MD5
//   - MD5 仅用于非安全场景（如文件校验）
//   - 不要在客户端进行密码加密，应在服务端加密
//
// =============================================================================
package encrypt

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

// Md5 计算数据的 MD5 哈希值
// 将字节数组转换为 32 位十六进制字符串
//
// 参数:
//   - str: 需要计算哈希的字节数组
//
// 返回:
//   - string: 32 位十六进制 MD5 哈希值
//
// 使用场景:
//   - 文件完整性校验
//   - 生成唯一标识符
//   - 数据去重（非安全场景）
//
// 示例:
//   hash := encrypt.Md5([]byte("hello world"))
//   // 输出: 5eb63bbbe01eeed093cb22bb8f5acdc3
//
// 注意:
//   - MD5 已被证明不安全，不要用于密码存储
//   - 仅用于非安全场景，如文件校验、缓存键生成
func Md5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

// GenPasswordHash 生成密码的 bcrypt 哈希值
// 使用 bcrypt 算法加密密码，自动加盐，防止彩虹表攻击
//
// 参数:
//   - password: 明文密码字节数组
//
// 返回:
//   - []byte: bcrypt 哈希值（包含盐值）
//   - error: 加密过程中的错误
//
// 使用场景:
//   - 用户注册时加密密码
//   - 用户修改密码时加密新密码
//
// 示例:
//   hash, err := encrypt.GenPasswordHash([]byte("password123"))
//   if err != nil {
//       return err
//   }
//   // 将 hash 存储到数据库
//
// 注意:
//   - 使用默认成本（bcrypt.DefaultCost = 10）
//   - 成本越高，加密越慢，但安全性越高
//   - 每次加密同一密码，结果都不同（因为盐值不同）
func GenPasswordHash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// ValidatePasswordHash 验证密码是否正确
// 使用 bcrypt 算法验证明文密码与哈希值是否匹配
//
// 参数:
//   - password: 明文密码字符串
//   - hashed: bcrypt 哈希值字符串（从数据库中读取）
//
// 返回:
//   - bool: true 表示密码正确，false 表示密码错误
//
// 使用场景:
//   - 用户登录时验证密码
//   - 修改密码时验证旧密码
//
// 示例:
//   // 用户登录
//   user, _ := userModel.FindByPhone(phone)
//   if !encrypt.ValidatePasswordHash(inputPassword, user.Password) {
//       return errors.New("密码错误")
//   }
//
// 注意:
//   - bcrypt 验证会自动提取盐值进行比对
//   - 验证失败不会返回错误，只返回 false
func ValidatePasswordHash(password string, hashed string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return false
	}
	return true
}
