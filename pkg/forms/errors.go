package forms

type errors map[string][]string

func (e errors) Add(field, text string) {
	e[field] = append(e[field], text)
}

func (e errors) Get(field string) string {
	es, ok := e[field]
	if !ok || len(es) == 0 {
		return ""
	}
	return es[0]
}
