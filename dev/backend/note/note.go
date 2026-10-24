package note

type Note struct {
	Author string
	Text   string
}

func (n Note) Summary() string {
	return n.Author + ": " + n.Text
}
