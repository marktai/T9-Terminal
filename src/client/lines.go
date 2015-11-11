package client

type lines struct {
	stringSlice []string
	start       int
}

func (l *lines) Add(s string) {
	l.stringSlice = append(l.stringSlice, s)
}

func (l *lines) Set(s string) {
	l.stringSlice = []string{s}
	if l.start >= len(l.stringSlice) {
		l.start = len(l.stringSlice) - 1
	}
}

func (l *lines) Clear() {
	l.stringSlice = []string{}
	l.start = 0
}

func (l *lines) Up() {
	l.start -= 1
	if l.start <= 0 {
		l.start = 0
	}
}

func (l *lines) Down() {
	l.start += 1
	if l.start >= len(l.stringSlice) {
		l.start = len(l.stringSlice) - 1
	}
	if l.start <= 0 {
		l.start = 0
	}
}

func (l *lines) String() string {
	retStr := ""
	// retStr += fmt.Sprint(l.start)
	for _, line := range l.stringSlice[l.start:] {
		retStr += line + "\n"
	}
	return retStr
}

func (l *lines) Length() int {
	return len(l.stringSlice) - l.start
}
