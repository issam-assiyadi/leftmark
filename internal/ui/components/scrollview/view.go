package scrollview

type Config struct {
	BaseName       string
	Title          string
	ScrollbarWidth int
}

type View struct {
	wrapperViewName string
	contentViewName string
	scrollViewName  string
	title           string
	scrollbarWidth  int
}

func New(cfg Config) *View {
	scrollbarWidth := cfg.ScrollbarWidth
	if scrollbarWidth < 2 {
		scrollbarWidth = 2
	}

	return &View{
		wrapperViewName: cfg.BaseName + "-wrapper",
		contentViewName: cfg.BaseName,
		scrollViewName:  cfg.BaseName + "-scroll",
		title:           cfg.Title,
		scrollbarWidth:  scrollbarWidth,
	}
}
