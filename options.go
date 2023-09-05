package gopdfgen

type OptionParser interface {
	Parse() []string
}

type BodyURLOption struct {
	bodyURL string
}

func (bo *BodyURLOption) Set(bodyURL string) {
	bo.bodyURL = bodyURL
}

func (bo *BodyURLOption) Parse() []string {
	args := make([]string, 0)
	if bo.bodyURL == "" {
		return args
	}
	args = append(args, bo.bodyURL)
	return args
}

type StringOption struct {
	option string
	value  string
}

func (so *StringOption) Set(value string) {
	so.value = value
}

func (so *StringOption) Parse() []string {
	args := make([]string, 0)
	if so.value == "" {
		return args
	}
	args = append(args, so.option)
	args = append(args, so.value)
	return args
}

type PasswordOption struct {
	password string
}

func (po *PasswordOption) Set(password string) {
	po.password = password
}
