package injopts

// ResolveConfOption configures how an injector resolves dependencies.
type ResolveConfOption uint8

const (
	ResolveOptDuckTyping ResolveConfOption = 1 << iota
	ResolveOptIsolated
	ResolveOptTracePath
	ResolveOptNone ResolveConfOption = 0
)

func (opt ResolveConfOption) Is(check ResolveConfOption) bool {
	return opt&check == check
}
