package models

type Portfolio struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Function string `json:"function"`
	Place    string `json:"place"`
	Email    string `json:"email"`
	Github   string `json:"github"`
	LinkedIn string `json:"linkedin"`
	Birthday string `json:"birthday"`
	ImgPath  string `json:"img_path"`
	About    string `json:"about"`

	Skills     []Skill     `json:"skills" gorm:"foreignKey:PortfolioID"`
	SoftSkills []SoftSkill `json:"soft_skills" gorm:"foreignKey:PortfolioID"`
	Languages  []Language  `json:"languages" gorm:"foreignKey:PortfolioID"`
	Education  []Education `json:"education" gorm:"foreignKey:PortfolioID"`
	Hobbies    []Hobby     `json:"hobbies" gorm:"foreignKey:PortfolioID"`
}

type Skill struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	PortfolioID string    `json:"portfolio_id" gorm:"not null"`
	Portfolio   Portfolio `json:"-" gorm:"foreignKey:PortfolioID"`
	Skill       string    `json:"skill"`
	Language    bool      `json:"language"`
	Framework   bool      `json:"framework"`
	Platform    bool      `json:"platform"`
	Tool        bool      `json:"tool"`
}

type SoftSkill struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	PortfolioID string    `json:"portfolio_id" gorm:"not null"`
	Portfolio   Portfolio `json:"-" gorm:"foreignKey:PortfolioID"`
	Skill       string    `json:"skill"`
}

type Language struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	PortfolioID string    `json:"portfolio_id" gorm:"not null"`
	Portfolio   Portfolio `json:"-" gorm:"foreignKey:PortfolioID"`
	Language    string    `json:"language"`
	Proficiency string    `json:"proficiency"`
}

type Education struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	PortfolioID string    `json:"portfolio_id" gorm:"not null"`
	Portfolio   Portfolio `json:"-" gorm:"foreignKey:PortfolioID"`
	Education   string    `json:"education"`
	School      string    `json:"school"`
	Start       string    `json:"start"`
	End         string    `json:"end"`

	Tasks []EducationTask `json:"tasks" gorm:"foreignKey:EducationID"`
}

type EducationTask struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	EducationID string    `json:"education_id" gorm:"not null"`
	Education   Education `json:"-" gorm:"foreignKey:EducationID"`
	Task        string    `json:"task"`
}

type Hobby struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	PortfolioID string    `json:"portfolio_id" gorm:"not null"`
	Portfolio   Portfolio `json:"-" gorm:"foreignKey:PortfolioID"`
	Emoji       string    `json:"emoji"`
	Hobby       string    `json:"hobby"`
	Discription string    `json:"discription"`
}
