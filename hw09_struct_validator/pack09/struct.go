package pack09

type User09 struct {
	ID   string
	name string `validate:"len:5"`
}
