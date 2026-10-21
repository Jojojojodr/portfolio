package views

import (
	"net/http"
	"net/mail"
	"net/url"
	"strings"

	"github.com/Jojojojodr/portfolio/controllers"
	"github.com/Jojojojodr/portfolio/database/models"
	"github.com/Jojojojodr/portfolio/frontend/admin"

	"github.com/gin-gonic/gin"
)

func HandleAdminPortfolioPage(c *gin.Context) {
	renderTempl(c, 200, admin.PortfolioProfilePage(c, adminPortfolio(), c.Query("error"), c.Query("success")))
}

func HandleAdminPortfolioSkillsPage(c *gin.Context) {
	renderTempl(c, 200, admin.PortfolioSkillsPage(c, adminPortfolio(), c.Query("error"), c.Query("success")))
}

func HandleAdminPortfolioEducationPage(c *gin.Context) {
	renderTempl(c, 200, admin.PortfolioEducationPage(c, adminPortfolio(), c.Query("error"), c.Query("success")))
}

func HandleAdminPortfolioHobbiesPage(c *gin.Context) {
	renderTempl(c, 200, admin.PortfolioHobbiesPage(c, adminPortfolio(), c.Query("error"), c.Query("success")))
}

func HandleAdminSavePortfolioProfile(c *gin.Context) {
	const page = "/portfolio"
	p, ok := portfolioForForm(c, page)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.PostForm("name"))
	email := strings.TrimSpace(c.PostForm("email"))
	github := strings.TrimSpace(c.PostForm("github"))
	linkedIn := strings.TrimSpace(c.PostForm("linkedin"))

	switch {
	case name == "":
		redirectPortfolio(c, page, "error", "Name is required")
		return
	case email != "" && !isEmail(email):
		redirectPortfolio(c, page, "error", "Invalid email")
		return
	case github != "" && !isURL(github):
		redirectPortfolio(c, page, "error", "Invalid GitHub URL")
		return
	case linkedIn != "" && !isURL(linkedIn):
		redirectPortfolio(c, page, "error", "Invalid LinkedIn URL")
		return
	}

	p.Name = name
	p.Function = strings.TrimSpace(c.PostForm("function"))
	p.Place = strings.TrimSpace(c.PostForm("place"))
	p.Email = email
	p.Github = github
	p.LinkedIn = linkedIn
	p.Birthday = strings.TrimSpace(c.PostForm("birthday"))
	p.ImgPath = strings.TrimSpace(c.PostForm("img_path"))

	if err := controllers.SavePortfolioInfo(p); err != nil {
		redirectPortfolio(c, page, "error", "Failed to save profile")
		return
	}
	redirectPortfolio(c, page, "success", "Profile saved")
}

func HandleAdminSavePortfolioAbout(c *gin.Context) {
	const page = "/portfolio#about"
	p, ok := portfolioForForm(c, page)
	if !ok {
		return
	}

	about := strings.TrimSpace(c.PostForm("about"))
	if about == "" {
		redirectPortfolio(c, page, "error", "About me is required")
		return
	}

	p.About = about
	if err := controllers.SavePortfolioInfo(p); err != nil {
		redirectPortfolio(c, page, "error", "Failed to save about me")
		return
	}
	redirectPortfolio(c, page, "success", "About me saved")
}

func HandleAdminSaveSkill(c *gin.Context) {
	const page = "/portfolio/skills#skills"
	p, ok := portfolioForForm(c, page)
	if !ok {
		return
	}

	skill := models.Skill{
		ID:          c.PostForm("id"),
		PortfolioID: p.ID,
		Skill:       strings.TrimSpace(c.PostForm("skill")),
		Language:    c.PostForm("language") == "on",
		Framework:   c.PostForm("framework") == "on",
		Platform:    c.PostForm("platform") == "on",
		Tool:        c.PostForm("tool") == "on",
	}
	if skill.Skill == "" {
		redirectPortfolio(c, page, "error", "Skill is required")
		return
	}
	if err := controllers.SaveSkill(&skill); err != nil {
		redirectPortfolio(c, page, "error", "Failed to save skill")
		return
	}
	redirectPortfolio(c, page, "success", "Skill '"+skill.Skill+"' saved")
}

func HandleAdminDeleteSkill(c *gin.Context) {
	deletePortfolioItem[models.Skill](c, "/portfolio/skills#skills", "Skill")
}

func HandleAdminSaveSoftSkill(c *gin.Context) {
	const page = "/portfolio/skills#soft-skills"
	p, ok := portfolioForForm(c, page)
	if !ok {
		return
	}

	softSkill := models.SoftSkill{
		ID:          c.PostForm("id"),
		PortfolioID: p.ID,
		Skill:       strings.TrimSpace(c.PostForm("skill")),
	}
	if softSkill.Skill == "" {
		redirectPortfolio(c, page, "error", "Soft skill is required")
		return
	}
	if err := controllers.SaveSoftSkill(&softSkill); err != nil {
		redirectPortfolio(c, page, "error", "Failed to save soft skill")
		return
	}
	redirectPortfolio(c, page, "success", "Soft skill '"+softSkill.Skill+"' saved")
}

func HandleAdminDeleteSoftSkill(c *gin.Context) {
	deletePortfolioItem[models.SoftSkill](c, "/portfolio/skills#soft-skills", "Soft skill")
}

func HandleAdminSaveLanguage(c *gin.Context) {
	const page = "/portfolio/skills#languages"
	p, ok := portfolioForForm(c, page)
	if !ok {
		return
	}

	language := models.Language{
		ID:          c.PostForm("id"),
		PortfolioID: p.ID,
		Language:    strings.TrimSpace(c.PostForm("language")),
		Proficiency: strings.TrimSpace(c.PostForm("proficiency")),
	}
	if language.Language == "" {
		redirectPortfolio(c, page, "error", "Language is required")
		return
	}
	if err := controllers.SaveLanguage(&language); err != nil {
		redirectPortfolio(c, page, "error", "Failed to save language")
		return
	}
	redirectPortfolio(c, page, "success", "Language '"+language.Language+"' saved")
}

func HandleAdminDeleteLanguage(c *gin.Context) {
	deletePortfolioItem[models.Language](c, "/portfolio/skills#languages", "Language")
}

func HandleAdminSaveEducation(c *gin.Context) {
	const page = "/portfolio/education"
	p, ok := portfolioForForm(c, page)
	if !ok {
		return
	}

	education := models.Education{
		ID:          c.PostForm("id"),
		PortfolioID: p.ID,
		Education:   strings.TrimSpace(c.PostForm("education")),
		School:      strings.TrimSpace(c.PostForm("school")),
		Start:       strings.TrimSpace(c.PostForm("start")),
		End:         strings.TrimSpace(c.PostForm("end")),
	}
	if education.Education == "" || education.School == "" {
		redirectPortfolio(c, page, "error", "Education and school are required")
		return
	}
	if err := controllers.SaveEducation(&education); err != nil {
		redirectPortfolio(c, page, "error", "Failed to save education")
		return
	}
	redirectPortfolio(c, page+"#"+education.ID, "success", "Education '"+education.Education+"' saved")
}

func HandleAdminDeleteEducation(c *gin.Context) {
	const page = "/portfolio/education"
	p, ok := portfolioForForm(c, page)
	if !ok {
		return
	}
	if err := controllers.DeleteEducation(p.ID, c.PostForm("id")); err != nil {
		redirectPortfolio(c, page, "error", "Failed to delete education")
		return
	}
	redirectPortfolio(c, page, "success", "Education deleted")
}

func HandleAdminSaveEducationTask(c *gin.Context) {
	education, page, ok := educationForForm(c)
	if !ok {
		return
	}

	task := models.EducationTask{
		ID:          c.PostForm("id"),
		EducationID: education.ID,
		Task:        strings.TrimSpace(c.PostForm("task")),
	}
	if task.Task == "" {
		redirectPortfolio(c, page, "error", "Task is required")
		return
	}
	if err := controllers.SaveEducationTask(&task); err != nil {
		redirectPortfolio(c, page, "error", "Failed to save task")
		return
	}
	redirectPortfolio(c, page, "success", "Task saved")
}

func HandleAdminDeleteEducationTask(c *gin.Context) {
	education, page, ok := educationForForm(c)
	if !ok {
		return
	}
	if err := controllers.DeleteEducationTask(education.ID, c.PostForm("id")); err != nil {
		redirectPortfolio(c, page, "error", "Failed to delete task")
		return
	}
	redirectPortfolio(c, page, "success", "Task deleted")
}

func HandleAdminSaveHobby(c *gin.Context) {
	const page = "/portfolio/hobbies"
	p, ok := portfolioForForm(c, page)
	if !ok {
		return
	}

	hobby := models.Hobby{
		ID:          c.PostForm("id"),
		PortfolioID: p.ID,
		Emoji:       strings.TrimSpace(c.PostForm("emoji")),
		Hobby:       strings.TrimSpace(c.PostForm("hobby")),
		Discription: strings.TrimSpace(c.PostForm("discription")),
	}
	if hobby.Hobby == "" {
		redirectPortfolio(c, page, "error", "Hobby is required")
		return
	}
	if err := controllers.SaveHobby(&hobby); err != nil {
		redirectPortfolio(c, page, "error", "Failed to save hobby")
		return
	}
	redirectPortfolio(c, page, "success", "Hobby '"+hobby.Hobby+"' saved")
}

func HandleAdminDeleteHobby(c *gin.Context) {
	deletePortfolioItem[models.Hobby](c, "/portfolio/hobbies", "Hobby")
}

func adminPortfolio() *models.Portfolio {
	p, err := controllers.GetPortfolio()
	if err != nil {
		return nil
	}
	return p
}

func portfolioForForm(c *gin.Context, page string) (*models.Portfolio, bool) {
	p, err := controllers.GetPortfolioInfo()
	if err != nil {
		redirectPortfolio(c, page, "error", "Portfolio not found")
		return nil, false
	}
	return p, true
}

// educationForForm loads the education the posted task belongs to and returns the page to go back to.
func educationForForm(c *gin.Context) (*models.Education, string, bool) {
	const page = "/portfolio/education"
	p, ok := portfolioForForm(c, page)
	if !ok {
		return nil, "", false
	}
	education, err := controllers.GetEducationByID(p.ID, c.PostForm("education_id"))
	if err != nil {
		redirectPortfolio(c, page, "error", "Education not found")
		return nil, "", false
	}
	return education, page + "#" + education.ID, true
}

func deletePortfolioItem[T any](c *gin.Context, page string, name string) {
	p, ok := portfolioForForm(c, page)
	if !ok {
		return
	}
	if err := controllers.DeletePortfolioItem[T](p.ID, c.PostForm("id")); err != nil {
		redirectPortfolio(c, page, "error", "Failed to delete "+strings.ToLower(name))
		return
	}
	redirectPortfolio(c, page, "success", name+" deleted")
}

// redirectPortfolio sends the admin back to page (which may end in a #fragment) with a message to show.
func redirectPortfolio(c *gin.Context, page string, kind string, message string) {
	path, fragment, _ := strings.Cut(page, "#")
	target := path + "?" + url.Values{kind: {message}}.Encode()
	if fragment != "" {
		target += "#" + fragment
	}
	c.Redirect(http.StatusSeeOther, target)
}

func isEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

func isURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	return err == nil && u.Host != ""
}
