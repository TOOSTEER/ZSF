package crossword

// Question is a single clue belonging to a crossword task.
// NOTE: crossword_questions is an addition on top of the schema given in
// the spec (which only defines crossword_answers). It is needed so the
// frontend can render the clue list / word list before the real content
// is supplied. // TODO: уточнить окончательную структуру, когда придут реальные кроссворды
type Question struct {
	ID             int    `json:"id"`
	TaskID         int    `json:"task_id"`
	QuestionNumber int    `json:"question_number"`
	QuestionText   string `json:"question_text"`
	Direction      string `json:"direction,omitempty"` // "across" | "down", classic type only
	AnswerLength   int    `json:"answer_length,omitempty"`
}

// StructureResponse is returned by GET /api/tasks/:id/crossword.
type StructureResponse struct {
	TaskID    int        `json:"task_id"`
	Title     string     `json:"title"`
	Subtype   string     `json:"crossword_subtype"`
	Questions []Question `json:"questions"`
}

// CheckAnswerInput is one submitted answer.
type CheckAnswerInput struct {
	QuestionNumber int    `json:"question_number"`
	Answer         string `json:"answer"`
}

// CheckRequest is the body of POST /api/crossword/:task_id/check.
type CheckRequest struct {
	Answers []CheckAnswerInput `json:"answers"`
}

// CheckResponse is the reply for POST /api/crossword/:task_id/check.
type CheckResponse struct {
	Success      bool  `json:"success"`
	Awarded      bool  `json:"awarded"`
	WrongNumbers []int `json:"wrong_numbers,omitempty"`
}
