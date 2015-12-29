package client

type lines struct {
	stringSlice []string
	start       int
	size        int
}

func NewLines() *lines {
	retLine := &lines{[]string{}, 0, 0}
	return retLine
}

func (l *lines) Add(s string) {
	l.stringSlice = append(l.stringSlice, s)
	l.size += 1
}

func (l *lines) Set(s string) {
	l.stringSlice = []string{s}
	if l.start >= len(l.stringSlice) {
		l.start = len(l.stringSlice) - 1
	}
	l.size = 1
}

func (l *lines) Clear() {
	l.stringSlice = []string{}
	l.start = 0
	l.size = 0
}

func (l *lines) Up() {
	l.start -= 1
	if l.start <= 0 {
		l.start = 0
	}
}

func (l *lines) Down() {
	l.start += 1
	if l.start >= l.size {
		l.start = l.size - 1
	}
	if l.start <= 0 {
		l.start = 0
	}
}

func (l *lines) String() string {
	retStr := ""
	// retStr += fmt.Sprint(l.start)
	start := true
	for _, line := range l.stringSlice[l.start:] {
		if !start {
			retStr += "\n"
		} else {
			start = false
		}
		retStr += line
	}
	return retStr
}

func (l *lines) Length() int {
	return l.size - l.start
}

func (l *lines) CalcHeight(width int) int {
	totalHeight := 0
	for _, line := range l.stringSlice[l.start:] {
		//stackoverflow.com/questions/2745074/fast-ceiling-of-an-integer-division-in-c-c
		totalHeight += (len(line) + width - 1) / width
	}
	return totalHeight
}
