package materi

import (
	"pintuniv-go/internal/database"
	"pintuniv-go/internal/models"
)

// ambil semua materi (tanpa lock logic)
func GetAllMateri() ([]models.Materi, error) {
	var list []models.Materi
	err := database.DB.
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}

// ambil materi + semua section (order by urutan)
func GetMateriBySlug(slug string) (*models.Materi, []models.MateriSection, error) {
	var materi models.Materi
	if err := database.DB.
		Where("slug = ?", slug).
		First(&materi).Error; err != nil {
		return nil, nil, err
	}

	var sections []models.MateriSection
	err := database.DB.
		Where("materi_id = ?", materi.ID).
		Order("`order` ASC").
		Find(&sections).Error
	if err != nil {
		return &materi, nil, err
	}

	return &materi, sections, nil
}

// cek apakah user PRO
// userID boleh nil (user belum login)
func IsUserPro(userID *uint) bool {
	if userID == nil {
		return false
	}

	var profile models.Profile
	if err := database.DB.
		Select("is_pro").
		Where("user_id = ?", *userID).
		First(&profile).Error; err != nil {
		return false
	}

	return profile.IsPro
}

// =====================
// ADMIN SERVICES
// =====================

func CreateMateri(m *models.Materi) error {
	return database.DB.Create(m).Error
}

func UpdateMateri(materiID uint, updates map[string]interface{}) error {
	return database.DB.
		Model(&models.Materi{}).
		Where("id = ?", materiID).
		Updates(updates).Error
}

func DeleteMateri(materiID uint) error {
	return database.DB.
		Where("id = ?", materiID).
		Delete(&models.Materi{}).Error
}

func GetMateriByID(materiID uint) (*models.Materi, error) {
	var materi models.Materi
	err := database.DB.First(&materi, materiID).Error
	return &materi, err
}

// list section by materi
func GetSectionsByMateriID(materiID uint) ([]models.MateriSection, error) {
	var sections []models.MateriSection
	err := database.DB.
		Where("materi_id = ?", materiID).
		Order("`order` ASC").
		Find(&sections).Error
	return sections, err
}

// create section
func CreateSection(section *models.MateriSection) error {
	return database.DB.Create(section).Error
}

// update section
func UpdateSection(sectionID uint, updates map[string]interface{}) error {
	return database.DB.
		Model(&models.MateriSection{}).
		Where("id = ?", sectionID).
		Updates(updates).Error
}

// delete section
func DeleteSection(sectionID uint) error {
	return database.DB.
		Where("id = ?", sectionID).
		Delete(&models.MateriSection{}).Error
}
