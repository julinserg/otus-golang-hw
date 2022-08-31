package pack09fortest

type User09 struct {
	ID   string
	name string `validate:"len:5"`
}

//lint:ignore U1000 Ignore unused function temporarily for test
func main() {
	st := &User09{ID: "1", name: "12"}
	_ = st.name
}
