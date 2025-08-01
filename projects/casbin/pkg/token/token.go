package token

import (
	"casbin_demo/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const SECRET_KEY = "sadhkljfasj;dflk"

func GenerateToken(user *models.User) string {
	// 创建一个Token对象
	token := jwt.New(jwt.SigningMethodHS256)

	// 设置Token的自定义声明
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["id"] = user.Id
	claims["tenant_id"] = user.TenantId
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 设置Token的过期时间

	// 使用密钥对Token进行签名，生成最终的Token字符串
	tokenString, _ := token.SignedString([]byte(SECRET_KEY))

	return tokenString
}

func ParseToken(tokenString string) (*models.User, error) {
	// 解析Token字符串
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证Token的签名方法是否有效
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("无效的签名方法：%v", token.Header["alg"])
	}

	// 返回Token中的声明部分
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("无效的Token")
	}

	userId, ok := claims["id"].(float64)
	tenantId, ok := claims["tenant_id"].(float64)
	username, ok := claims["username"].(string)
	var user models.User
	user = models.User{
		CommonField: models.CommonField{Id: uint64(userId)},
		Username:    username,
		TenantId:    uint64(tenantId),
	}

	return &user, nil
}
