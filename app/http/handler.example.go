package http

func (h Handler) getExample(path struct {
	ID int `json:"id,string"`
}, query struct {
	IDK []string `json:"idk"`
}, req map[string]string) (id []string, err error) {
	println(query.IDK)
	return query.IDK, nil
}
