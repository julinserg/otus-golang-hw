package pack09fortest

type User09 struct {
	ID   string
	name string `validate:"len:5"`
}

func Main() {
	st := &User09{ID: "1", name: "12"}
	_ = st.name
}
