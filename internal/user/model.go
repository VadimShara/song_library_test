package user

type Song struct {
	ID          string `json:"id" `
	Group       string `json:"group" example:"Pizza"`
	Song        string `json:"song" example:"Lift"`
	ReleaseDate string `json:"releaseDate" example:"24.10.2014"`
	Text        string `json:"text" example:"Скорее, минуты летите, чтобы я вас не заметил\nИ поспешите на третий, свободы мне принесите\nСкорее, минуты летите, чтобы я вас не заметил\nИ поспешите на третий, откройте, освободите\n"`
	Link        string `json:"link" example:"https://www.youtube.com/watch?v=Eyp3bnl5Cng"`
}
