package injopts

type CacheConfOption uint8

const (
	CacheOptAllowOverride CacheConfOption = 1 << iota
	CacheOptReturnAll
	CacheOptNone CacheConfOption = 0
)

func (opt CacheConfOption) Is(check CacheConfOption) bool {
	return opt&check == check
}
