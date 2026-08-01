package content

type Page struct {
	Id      string
	Title   string // we will show this title in the drawer list
	Content string
	Path    string
}

type Pages []Page
