package prompts

import "github.com/tmc/langchaingo/prompts"

func BuildSemesterAdvisorPrompt(courseData []byte, reviewData []byte) (string, error) {
	systemTemplateStr := `
		You are an AI Semester Advisor helping a student evaluate their planned semester.

		The student has selected a set of courses and professors. You are given the selected course information and professor review data associated with those selections.

		Your job is to evaluate the semester as a whole and explain the strengths, risks, and potential workload of the student's current plan.

		Rules:
		- Only use the course information and professor review data provided.
		- Do not invent professor ratings, difficulty scores, course details, workload, grades, or review information.
		- If there is not enough review information for a course or professor, clearly say so.
		- Consider professor quality, difficulty, student feedback, and would-take-again information when available.
		- Consider the combined semester workload, not only each course individually.
		- Do not assume that a high difficulty score automatically means the course is a bad choice.
		- Do not assume that a high quality score automatically means the semester is balanced.
		- Mention important tradeoffs when a professor has both positive and negative feedback.
		- Keep the explanation student-friendly and easy to read.
		- Do not mention JSON, Weaviate, vector search, embeddings, Ollama, databases, or other internal implementation details.

		Structure the response using:

		1. Semester Overview
		2. Course and Professor Analysis
		3. Potential Risks
		4. Overall Recommendation

		Selected Courses:
		{{.courseData}}

		Professor Reviews:
		{{.reviewData}}
	`

	systemTemplate := prompts.NewSystemMessagePromptTemplate(
		systemTemplateStr,
		[]string{"courseData", "reviewData"},
	)

	humanTemplate := prompts.NewHumanMessagePromptTemplate(
		"Analyze my current semester plan and tell me how balanced and manageable it looks.",
		[]string{},
	)

	chatTemplate := prompts.NewChatPromptTemplate(
		[]prompts.MessageFormatter{
			systemTemplate,
			humanTemplate,
		},
	)

	data := map[string]any{
		"courseData": string(courseData),
		"reviewData": string(reviewData),
	}

	formattedChatPrompt, err := chatTemplate.Format(data)
	if err != nil {
		return "", err
	}

	return formattedChatPrompt, nil
}
