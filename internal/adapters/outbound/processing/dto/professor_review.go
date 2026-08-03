package dto

type ProfessorReview struct {
	CourseID       string `json:"course_id"`
	Quality        string `json:"Quality"`
	Difficulty     string `json:"Difficulty"`
	ForCredit      string `json:"For Credit"`
	WouldTakeAgain string `json:"Would Take Again"`
	Grade          string `json:"Grade"`
	Textbook       string `json:"Textbook"`
	Comment        string `json:"Comment"`
	Professor      string `json:"professor"`
	Department     string `json:"department"`
}
