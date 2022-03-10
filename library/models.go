package library

type Album struct {
	Title  string
	Artist string
}

type Track struct {
	Artist string
	Title  string
	Album  *Album
}
