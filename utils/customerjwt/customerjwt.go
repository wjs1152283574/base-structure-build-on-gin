package customerjwt

import (
	"errors"
	"goweb/utils/response"
	"goweb/utils/statuscode"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// TokenExpired 一些常量
var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey          string = "newtrekWang"
)

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")
		var token = ""
		if tokenStr == "" {
			response.ReturnJSON(c, http.StatusOK, statuscode.FailToken.Code, statuscode.FailToken.Msg, nil)
			c.Abort()
			return
		} else {
			// 自定义token拼接格式 -- 前端请求头需添加同样格式(请求拦截器里面设置)
			if len(strings.Split(tokenStr, " ")) > 1 && strings.Split(tokenStr, " ")[0] == "GOJWT" {
				token = strings.Split(tokenStr, " ")[1]
			}
		}
		if token == "" {
			response.ReturnJSON(c, http.StatusOK, statuscode.FailToken.Code, statuscode.FailToken.Msg, nil)
			c.Abort()
			return
		}

		log.Print("cost: ", time.Now().Format("2006-01-02 15:04:05"))

		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				response.ReturnJSON(c, http.StatusOK, statuscode.ExprieToken.Code, statuscode.ExprieToken.Msg, nil)
				c.Abort()
				return
			}
			response.ReturnJSON(c, http.StatusOK, statuscode.FailToken.Code, statuscode.FailToken.Msg, nil)
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("userID", claims.ID)
	}
}

// CustomClaims 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	ID int64 `json:"id"` // 用户ID
	jwt.StandardClaims
}

// NewJWT 新建一个jwt实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

// GetSignKey 获取signKey
func GetSignKey() string {
	return SignKey
}

// SetSignKey 这是SignKey
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

// CreateToken 生成一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	// 设置自定义token过期时间
	claims.StandardClaims.ExpiresAt = time.Now().Add(time.Minute * 60).Unix() // 测试时60S过期 : time.Second * 60
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析Tokne
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// RefreshToken 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(12 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
