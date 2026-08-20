package scrollview

type Config struct {
	BaseName       string
	Title          string
	ScrollbarWidth int
	// HideScrollbar makes the view scrollable via keyboard/selection
	// without ever drawing a scrollbar column.
	HideScrollbar bool
}

type View struct {
	wrapperViewName string
	contentViewName string
	scrollViewName  string
	title           string
	scrollbarWidth  int
	hideScrollbar   bool
}

func New(cfg Config) *View {
	scrollbarWidth := cfg.ScrollbarWidth
	if !cfg.HideScrollbar && scrollbarWidth < 2 {
		scrollbarWidth = 2
	}
	if cfg.HideScrollbar {
		scrollbarWidth = 0
	}

	return &View{
		wrapperViewName: cfg.BaseName + "-wrapper",
		contentViewName: cfg.BaseName,
		scrollViewName:  cfg.BaseName + "-scroll",
		title:           cfg.Title,
		scrollbarWidth:  scrollbarWidth,
		hideScrollbar:   cfg.HideScrollbar,
	}
}

func (v *View) WrapperName() string { return v.wrapperViewName }

func (v *View) ContentViewName() string { return v.contentViewName }

func (v *View) ScrollViewName() string { return v.scrollViewName }
