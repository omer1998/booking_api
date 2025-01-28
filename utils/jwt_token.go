package utils

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func CreateJWTToken(email string, id string, validForHour int) (*string, error) {
	// claims := &jwt.RegisteredClaims{
	// 	ExpiresAt: &jwt.NumericDate{Time: validFor,},

	// }
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyClaims{
		Id:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour * time.Duration(validForHour))),
		},
	})

	ss, err := token.SignedString([]byte("omerfaris"))
	if err != nil {
		return nil, err
	}
	return &ss, nil
}

func ParseJWTToken(tokenString string) (*jwt.Token, error) {
	MyCustomClaims := new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, MyCustomClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte("omerfaris"), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	} else if _, ok := token.Claims.(*MyClaims); ok { // this is just a type assertion,.(*MyClaims) is used to assert that the Claims field of the token is of the specific type *MyClaims. The *MyClaims type is a pointer to a struct defined as MyClaims.
		exp, err := token.Claims.GetExpirationTime()
		if err != nil {
			return nil, err
		}
		if time.Now().After(exp.Time) {
			return nil, errors.New("token is expired")
		}
		return token, nil
	} else {
		log.Println("unknown claims type, cannot proceed")
		return nil, errors.New("unknown claims type, cannot proceed")
	}

}
