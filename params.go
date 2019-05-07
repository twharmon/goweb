package gorouter

// Param .
type Param struct {
	Key   string
	Value string
}

// Params .
type Params []*Param

// Get .
func (ps Params) Get(key string) string {
	for i := range ps {
		if ps[i].Key == key {
			return ps[i].Value
		}
	}
	return ""
}
