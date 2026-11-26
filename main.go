package main

import (
	"grading-system/internal/adapters"
	"grading-system/internal/application"
)

func main() {
	studentRepo := adapters.NewInMemoryStudentRepo()
	evaluationRepo := adapters.NewInMemoryEvaluationRepo()
	rubricRepo := adapters.NewInMemoryRubricRepo()

	facade := application.NewAssessmentFacade(
		studentRepo,
		evaluationRepo,
		rubricRepo,
	)

	app := adapters.NewConsoleApp(facade)
	app.Run()
}
