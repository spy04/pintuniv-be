package profile

import (
	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"
)

func GetProfileByUserID(userID uint) (*models.Profile, error) {
	var profile models.Profile
	err := database.DB.Where("user_id = ?", userID).First(&profile).Error
	return &profile, err
}

func UpdateProfile(userID uint, req UpdateProfileRequest) error {
	updates := map[string]interface{}{}

	if req.FullName != nil {
		updates["full_name"] = *req.FullName
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}

	if len(updates) == 0 {
		return nil
	}

	return database.DB.
		Model(&models.Profile{}).
		Where("user_id = ?", userID).
		Updates(updates).Error
}

func UpdateAvatar(userID uint, avatarURL string) error {
	return database.DB.
		Model(&models.Profile{}).
		Where("user_id = ?", userID).
		Update("avatar_url", avatarURL).Error
}

func GetAvatarByUserID(userID uint) (string, error) {
	var profile models.Profile
	err := database.DB.
		Select("avatar_url").
		Where("user_id = ?", userID).
		First(&profile).Error

	return profile.AvatarURL, err
}
