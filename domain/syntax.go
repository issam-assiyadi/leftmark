package domain

// CommentSyntax describes how a language spells comments, for the purpose of
// finding and tagging a marker on a single line. This is deliberately not a
// parser: it only needs enough to locate a comment opener and, for
// block-comment-only languages, its closer on the same line.
type CommentSyntax struct {
	// LinePrefixes are tokens that open a single-line comment, e.g. "//".
	LinePrefixes []string
	// BlockOpen/BlockClose are the block-comment delimiters, e.g. "/*" and
	// "*/". Empty if the language has no block form.
	BlockOpen  string
	BlockClose string
}

var syntaxByExt = map[string]CommentSyntax{
	".go":   {LinePrefixes: []string{"//"}, BlockOpen: "/*", BlockClose: "*/"},
	".js":   {LinePrefixes: []string{"//"}, BlockOpen: "/*", BlockClose: "*/"},
	".jsx":  {LinePrefixes: []string{"//"}, BlockOpen: "/*", BlockClose: "*/"},
	".ts":   {LinePrefixes: []string{"//"}, BlockOpen: "/*", BlockClose: "*/"},
	".tsx":  {LinePrefixes: []string{"//"}, BlockOpen: "/*", BlockClose: "*/"},
	".java": {LinePrefixes: []string{"//"}, BlockOpen: "/*", BlockClose: "*/"},
	".c":    {LinePrefixes: []string{"//"}, BlockOpen: "/*", BlockClose: "*/"},
	".h":    {LinePrefixes: []string{"//"}, BlockOpen: "/*", BlockClose: "*/"},
	".cpp":  {LinePrefixes: []string{"//"}, BlockOpen: "/*", BlockClose: "*/"},
	".hpp":  {LinePrefixes: []string{"//"}, BlockOpen: "/*", BlockClose: "*/"},
	".cs":   {LinePrefixes: []string{"//"}, BlockOpen: "/*", BlockClose: "*/"},
	".rs":   {LinePrefixes: []string{"//"}, BlockOpen: "/*", BlockClose: "*/"},

	".py":   {LinePrefixes: []string{"#"}},
	".rb":   {LinePrefixes: []string{"#"}},
	".sh":   {LinePrefixes: []string{"#"}},
	".bash": {LinePrefixes: []string{"#"}},
	".yml":  {LinePrefixes: []string{"#"}},
	".yaml": {LinePrefixes: []string{"#"}},
	".toml": {LinePrefixes: []string{"#"}},

	".sql": {LinePrefixes: []string{"--"}},

	".html": {BlockOpen: "<!--", BlockClose: "-->"},
	".htm":  {BlockOpen: "<!--", BlockClose: "-->"},
	".xml":  {BlockOpen: "<!--", BlockClose: "-->"},
	".vue":  {BlockOpen: "<!--", BlockClose: "-->"},

	".css":  {BlockOpen: "/*", BlockClose: "*/"},
	".scss": {BlockOpen: "/*", BlockClose: "*/"},
}

// SyntaxForExt looks up the comment syntax for a file extension (including
// the leading dot, as returned by filepath.Ext).
func SyntaxForExt(ext string) (CommentSyntax, bool) {
	s, ok := syntaxByExt[ext]
	return s, ok
}
