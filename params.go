package goweb

type param struct {
	key   string
	value string
}

type params []param

func (ps params) get(key string) string {
	for i := range ps {
		if ps[i].key == key {
			return ps[i].value
		}
	}
	return ""
}
