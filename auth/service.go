package auth

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(userID int ) (string, error)
	ValidateToken(encodedToken string) (*jwt.Token, error) //encoded token adalah memasukan token yang udah di dapat
}

type jwtService struct {}

var SECRETKEY = []byte("secretkeycrowdfunding") 

func NewService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) GenerateToken(userID int) (string, error) {

	claim := jwt.MapClaims{}
	claim["user_id"] = userID //key user_id value nya adalah userID yang di dapat dari parameter
	

	// GENERATE TOKEN // BIKIN TOKEN  
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// SIGN TOKEN // TANDA TANGANI TOKEN
	signedToken, err := token.SignedString(SECRETKEY)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil

}

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {

	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC) // cara baca nya jika token methodnya hmac true

		if !ok { //jika bukan hmac
			return nil, errors.New("Invalid token")
		}

		return []byte(SECRETKEY), nil //return pengembalian func(token *jwt.Token)
		
	})

	if err != nil {
		return token, err
	}

	return token, nil


}

