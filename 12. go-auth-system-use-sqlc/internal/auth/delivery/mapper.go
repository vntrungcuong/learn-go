package delivery

import (
	"time"

	db "go-auth-system/internal/db/sqlc"
)

// toUserData chuyển đổi từ db.User (SQLC) sang UserData (DTO)
func toUserData(user db.User) UserData {
	return UserData{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.Fullname,
		CreatedAt: parseTime(user.CreatedAt),
		UpdatedAt: parseTime(user.UpdatedAt),
	}
}

// toAuthData đóng gói bộ đôi token vào DTO
func toAuthData(accessToken, refreshToken string) AuthData {
	return AuthData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

// parseTime là helper xử lý an toàn việc ép kiểu từ interface{} sang time.Time
func parseTime(raw interface{}) time.Time {
	if raw == nil {
		return time.Time{}
	}

	t, ok := raw.(time.Time)
	if !ok {
		// Nếu SQL trả về chuỗi (string) thay vì object time.Time,
		// có thể thực hiện parse thủ công tại đây nếu cần.
		return time.Time{}
	}

	return t
}
