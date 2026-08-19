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

func BuildSemesterPlannerPrompt(courseData []byte, reviewData []byte) (string, error) {
	systemTemplateStr := `
		You are an AI Semester Planner helping a student improve their current semester plan.

		The student has already selected a set of courses and professors. You are given the current course information and professor review data associated with those selections.

		Your job is to identify which parts of the semester plan could be improved and explain what better alternatives should look like.

		Rules:
		- Only use the course information and professor review data provided.
		- Do not invent professors, courses, ratings, difficulty scores, grades, or review information.
		- If there is not enough information to recommend a specific alternative professor or course, say so clearly.
		- Consider professor quality, difficulty, would-take-again information, student comments, and the overall semester workload.
		- Look for courses or professor selections that may create unnecessary difficulty or workload.
		- Preserve strong course and professor selections when there is no clear reason to replace them.
		- Do not recommend replacing a course only because it is difficult.
		- Explain why a selection should be kept or reconsidered.
		- Focus on creating a more balanced and manageable semester.
		- Keep recommendations practical, student-friendly, and easy to understand.
		- Do not mention JSON, Weaviate, vector search, embeddings, Ollama, databases, or other internal implementation details.

		Structure the response using:

		1. Current Plan Assessment
		2. Courses to Keep
		3. Courses or Professors to Reconsider
		4. Recommended Improvements
		5. Suggested Semester Strategy

		Current Semester Plan:
		{{.courseData}}

		Professor Reviews:
		{{.reviewData}}
	`

	systemTemplate := prompts.NewSystemMessagePromptTemplate(
		systemTemplateStr,
		[]string{"courseData", "reviewData"},
	)

	humanTemplate := prompts.NewHumanMessagePromptTemplate(
		"Review my current semester plan and help me find better alternatives where improvements are needed.",
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
