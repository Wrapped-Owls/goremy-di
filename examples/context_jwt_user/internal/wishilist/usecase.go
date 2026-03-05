package wishilist

type User struct {
	ID       string
	Username string
}

type UseCase struct {
	user User
	repo Repository
}

func NewWishlistUseCase(user User, wishlistRepo Repository) UseCase {
	return UseCase{user: user, repo: wishlistRepo}
}

func (uc UseCase) ReadMyWishlist() []string {
	return uc.repo.ListByUserID(uc.user.ID)
}
