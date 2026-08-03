package domain

type ProfessorReview struct {
	ID             string
	UploadedFileID string
	CourseID       string
	Quality        float64
	Difficulty     float64
	ForCredit      *bool
	WouldTakeAgain *bool
	Grade          string
	Textbook       *bool
	Comment        string
	Professor      string
	Department     string
}
