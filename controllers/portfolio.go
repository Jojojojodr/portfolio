package controllers

import (
	"log"
	"net/http"

	"github.com/Jojojojodr/portfolio"
	"github.com/Jojojojodr/portfolio/database/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetPortfolio() (*models.Portfolio, error) {
	var p models.Portfolio
	err := portfolio.Data.GetDB().
		Preload("Skills").
		Preload("SoftSkills").
		Preload("Languages").
		Preload("Education.Tasks").
		Preload("Hobbies").
		First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetSkills() ([]models.Skill, error) {
	var skills []models.Skill
	err := portfolio.Data.GetDB().Find(&skills).Error
	if err != nil {
		return nil, err
	}
	return skills, nil
}

func GetSoftSkills() ([]models.SoftSkill, error) {
	var softSkills []models.SoftSkill
	err := portfolio.Data.GetDB().Find(&softSkills).Error
	if err != nil {
		return nil, err
	}
	return softSkills, nil
}

func GetLanguages() ([]models.Language, error) {
	var languages []models.Language
	err := portfolio.Data.GetDB().Find(&languages).Error
	if err != nil {
		return nil, err
	}
	return languages, nil
}

func GetEducation() ([]models.Education, error) {
	var education []models.Education
	err := portfolio.Data.GetDB().Preload("Tasks").Find(&education).Error
	if err != nil {
		return nil, err
	}
	return education, nil
}

func GetHobbies() ([]models.Hobby, error) {
	var hobbies []models.Hobby
	err := portfolio.Data.GetDB().Find(&hobbies).Error
	if err != nil {
		return nil, err
	}
	return hobbies, nil
}

func LanguagesAndFrameworks(skills []models.Skill) []models.Skill {
	var filtered []models.Skill
	for _, skill := range skills {
		if skill.Language || skill.Framework {
			filtered = append(filtered, skill)
		}
	}
	return filtered
}

func PlatformsAndTools(skills []models.Skill) []models.Skill {
	var filtered []models.Skill
	for _, skill := range skills {
		if skill.Platform || skill.Tool {
			filtered = append(filtered, skill)
		}
	}
	return filtered
}

func GetPortfolioAPI(c *gin.Context) {
	p, err := GetPortfolio()
	if err != nil {
		log.Printf("Error fetching portfolio: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func GetSkillsAPI(c *gin.Context) {
	skills, err := GetSkills()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch skills"})
		return
	}
	c.JSON(http.StatusOK, skills)
}

func GetSoftSkillsAPI(c *gin.Context) {
	softSkills, err := GetSoftSkills()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch soft skills"})
		return
	}
	c.JSON(http.StatusOK, softSkills)
}

func GetLanguagesAPI(c *gin.Context) {
	languages, err := GetLanguages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch languages"})
		return
	}
	c.JSON(http.StatusOK, languages)
}

func GetEducationAPI(c *gin.Context) {
	education, err := GetEducation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch education"})
		return
	}
	c.JSON(http.StatusOK, education)
}

func GetHobbiesAPI(c *gin.Context) {
	hobbies, err := GetHobbies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch hobbies"})
		return
	}
	c.JSON(http.StatusOK, hobbies)
}

func GetPortfolioInfo() (*models.Portfolio, error) {
	var p models.Portfolio
	if err := portfolio.Data.GetDB().First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func SavePortfolioInfo(p *models.Portfolio) error {
	return portfolio.Data.GetDB().Omit(clause.Associations).Save(p).Error
}

func SaveSkill(s *models.Skill) error {
	ensureID(&s.ID, "skill")
	return portfolio.Data.GetDB().Omit(clause.Associations).Save(s).Error
}

func SaveSoftSkill(s *models.SoftSkill) error {
	ensureID(&s.ID, "soft-skill")
	return portfolio.Data.GetDB().Omit(clause.Associations).Save(s).Error
}

func SaveLanguage(l *models.Language) error {
	ensureID(&l.ID, "language")
	return portfolio.Data.GetDB().Omit(clause.Associations).Save(l).Error
}

func SaveEducation(e *models.Education) error {
	ensureID(&e.ID, "education")
	return portfolio.Data.GetDB().Omit(clause.Associations).Save(e).Error
}

func SaveEducationTask(t *models.EducationTask) error {
	ensureID(&t.ID, "education-task")
	return portfolio.Data.GetDB().Omit(clause.Associations).Save(t).Error
}

func SaveHobby(h *models.Hobby) error {
	ensureID(&h.ID, "hobby")
	return portfolio.Data.GetDB().Omit(clause.Associations).Save(h).Error
}

// DeletePortfolioItem deletes a skill, soft skill, language or hobby of the given portfolio.
func DeletePortfolioItem[T any](portfolioID, id string) error {
	return portfolio.Data.GetDB().Where("id = ? AND portfolio_id = ?", id, portfolioID).Delete(new(T)).Error
}

func GetEducationByID(portfolioID, id string) (*models.Education, error) {
	var e models.Education
	if err := portfolio.Data.GetDB().Where("id = ? AND portfolio_id = ?", id, portfolioID).First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func DeleteEducation(portfolioID, id string) error {
	return portfolio.Data.GetDB().Transaction(func(tx *gorm.DB) error {
		owned := tx.Model(&models.Education{}).Select("id").Where("id = ? AND portfolio_id = ?", id, portfolioID)
		if err := tx.Where("education_id IN (?)", owned).Delete(&models.EducationTask{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ? AND portfolio_id = ?", id, portfolioID).Delete(&models.Education{}).Error
	})
}

func DeleteEducationTask(educationID, id string) error {
	return portfolio.Data.GetDB().Where("id = ? AND education_id = ?", id, educationID).Delete(&models.EducationTask{}).Error
}
