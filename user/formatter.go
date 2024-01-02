package user

// fungsi struct ini agar mengubah menjadi format json

type UserFormatter struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Token      string `json:"token"`
	ImageURL   string `jsn:"image_url"`
}

func FormatUser(user User, token string) UserFormatter {
	formatter := UserFormatter{
		ID:         user.ID,
		Name:       user.Name,
		Occupation: user.Occupation,
		Token:      token,
		ImageURL:   user.AvatarFileName,
	}

	return formatter

}