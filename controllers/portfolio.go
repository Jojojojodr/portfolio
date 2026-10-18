package controllers

import (
	"log"
	"net/http"

	"github.com/Jojojojodr/portfolio"
	"github.com/Jojojojodr/portfolio/database/models"

	"github.com/gin-gonic/gin"
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
