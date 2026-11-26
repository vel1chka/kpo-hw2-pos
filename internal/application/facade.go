package application

import (
	"fmt"
	"grading-system/internal/domain"
)

type AssessmentFacade struct {
	studentRepo    domain.StudentRepository
	evaluationRepo domain.EvaluationRepository
	rubricRepo     domain.RubricRepository
}

func NewAssessmentFacade(studentRepo domain.StudentRepository, evaluationRepo domain.EvaluationRepository, rubricRepo domain.RubricRepository) *AssessmentFacade {
	return &AssessmentFacade{
		studentRepo:    studentRepo,
		evaluationRepo: evaluationRepo,
		rubricRepo:     rubricRepo,
	}
}

func (f *AssessmentFacade) GradeStudent(email, criterionName string, score int) error {
	newEmail, err := domain.NewEmail(email)
	if err != nil {
		return err
	}

	_, err = f.studentRepo.FindByEmail(newEmail)
	if err != nil {
		return err
	}

	rubric, err := f.rubricRepo.Get()
	if err != nil {
		return err
	}

	var foundCriterion domain.Criterion
	for _, criterion := range rubric.Criteria() {
		if criterion.Name() == criterionName {
			foundCriterion = criterion
			break
		}
	}
	if foundCriterion.Name() == "" {
		return fmt.Errorf("критерий %s не найден в рубрике", criterionName)
	}

	evaluation, err := domain.NewEvaluation(newEmail, foundCriterion, score)
	if err != nil {
		return err
	}
	return f.evaluationRepo.Save(evaluation)
}

type StudentResult struct {
	Name       string
	Email      string
	FinalScore float64
}

func (f *AssessmentFacade) GetAllResults() ([]StudentResult, error) {
	students, err := f.studentRepo.GetAll()
	if err != nil {
		return nil, err
	}

	rubric, err := f.rubricRepo.Get()
	if err != nil {
		return nil, err
	}

	result := make([]StudentResult, len(students))
	for i, student := range students {
		evaluations, err := f.evaluationRepo.FindByStudent(student.Email())
		if err != nil {
			return nil, err
		}
		finalScore, err := domain.CalculateFinalScore(rubric, evaluations)
		if err != nil {
			return nil, err
		}
		result[i] = StudentResult{
			Name:       student.Name(),
			Email:      string(student.Email()),
			FinalScore: finalScore,
		}
	}

	return result, nil
}

func (f *AssessmentFacade) SetRubric(rubric domain.Rubric) error {
	if len(rubric.Criteria()) == 0 {
		return fmt.Errorf("нельзя сохранить пустую рубрику")
	}
	return f.rubricRepo.Save(rubric)
}

func (f *AssessmentFacade) RegisterStudent(name, email string) error {
	newEmail, err := domain.NewEmail(email)
	if err != nil {
		return err
	}

	student, err := domain.NewStudent(name, newEmail)
	if err != nil {
		return err
	}

	return f.studentRepo.Save(student)
}

func (f *AssessmentFacade) GetRubric() (domain.Rubric, error) {
	return f.rubricRepo.Get()
}

func (f *AssessmentFacade) GetStudents() ([]domain.Student, error) {
	return f.studentRepo.GetAll()
}

func (f *AssessmentFacade) GetStudentEvaluations(email string) ([]domain.Evaluation, error) {
	newEmail, err := domain.NewEmail(email)
	if err != nil {
		return nil, err
	}

	if _, err := f.studentRepo.FindByEmail(newEmail); err != nil {
		return nil, err
	}

	return f.evaluationRepo.FindByStudent(newEmail)
}

func (f *AssessmentFacade) ClearEvaluations() error {
	return f.evaluationRepo.Clear()
}
