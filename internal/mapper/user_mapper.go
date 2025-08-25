package mapper

import (
	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/itsLeonB/cocoon/internal/entity"
)

func UserToAuthData(user entity.User) map[string]any {
	return map[string]any{
		appconstant.ContextUserID: user.ID,
	}
}
