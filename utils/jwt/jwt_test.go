package jwt

import (
	"testing"
	"time"
)

func TestCreateJWT(t *testing.T) {

	jwt, err := CreateJWT(map[string]interface{}{
		"email": "sdfsdf@qq.com",
		"time":  time.Now().UnixMilli(),
	})
	if err != nil {
		t.Log(err.Error())
	}

	t.Log(jwt)

	t.Fail()
}

func TestValidJWT(t *testing.T) {
	claims, err := ValidateJWT("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InNkZnNkZkBxcS5jb20iLCJ0aW1lIjoxNzA4MTc2NzgyfQ.UTHFEfBkQRMMHoWxXczxDfxGr28-xlmtuqiQT211YTs")
	if err != nil {
		t.Log(err.Error())
	}

	t.Log(claims)

	t.Fail()
}
