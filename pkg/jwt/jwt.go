package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// 过期时间
const TokenExpireDuration = time.Hour * 2

// secret加密
var MySecret = []byte("nkinder")

// 保存数据
type MyClaim struct {
	Username string `json:"username"`
	UserId   int64  `json:"user_id"`
	jwt.StandardClaims
}

func GenToken(username string, userId int64) (string, error) {
	//自定义token结构体
	v := MyClaim{
		username,
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "bluebell",
		},
	}
	//按照某种签名方法生成token
	//注意：SigningMethodHS256和SigningMethodES256的区别就是后者传入的不是byte[]字节，而且也没有转换string的方法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, v)
	//返回带有自己签名的字符串
	return token.SignedString(MySecret)
}

// 解码token
func ParasToken(tokenString string) (*MyClaim, error) {
	claim := new(MyClaim)
	//解码方式是传进去你的token、你声明的结构体（header和overload部分）、你的最后的签名方式
	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return claim, nil
	}
	return nil, errors.New("token is valid")

}
