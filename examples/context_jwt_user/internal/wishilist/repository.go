package wishilist

type Repository interface {
	ListByUserID(userID string) []string
}

type InMemoryWishlistRepo struct {
	rows map[string][]string
}

func NewInMemoryWishlistRepo() InMemoryWishlistRepo {
	return InMemoryWishlistRepo{
		rows: map[string][]string{
			"u-alice": {"Mechanical keyboard", "Ergonomic chair"},
			"u-bob":   {"Go mug", "Noise-cancelling headphones"},
		},
	}
}

func (repo InMemoryWishlistRepo) ListByUserID(userID string) []string {
	items := repo.rows[userID]
	copyItems := make([]string, len(items))
	copy(copyItems, items)
	return copyItems
}
