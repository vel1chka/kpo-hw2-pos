package adapters

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"grading-system/internal/application"
	"grading-system/internal/domain"
	"os"
)

type ExportPipeline struct {
	FilePath string
	Fetch    func() ([]application.StudentResult, error)
	Format   func(data []application.StudentResult) ([]byte, error)
}

func (p ExportPipeline) Run() error {
	results, err := p.Fetch()
	if err != nil {
		return err
	}

	bytes, err := p.Format(results)
	if err != nil {
		return err
	}

	return os.WriteFile(p.FilePath, bytes, 0644)
}

type InMemoryStudentRepo struct {
	students map[domain.Email]domain.Student
}

func NewInMemoryStudentRepo() *InMemoryStudentRepo {
	return &InMemoryStudentRepo{
		students: make(map[domain.Email]domain.Student),
	}
}

func (r *InMemoryStudentRepo) Save(student domain.Student) error {
	r.students[student.Email()] = student
	return nil
}

func (r *InMemoryStudentRepo) FindByEmail(email domain.Email) (domain.Student, error) {
	student, ok := r.students[email]
	if !ok {
		return domain.Student{}, fmt.Errorf("студент с email %s не найден", email)
	}
	return student, nil
}

func (r *InMemoryStudentRepo) GetAll() ([]domain.Student, error) {
	var students []domain.Student
	for _, student := range r.students {
		students = append(students, student)
	}
	return students, nil
}

type InMemoryEvaluationRepo struct {
	evaluations map[domain.Email]map[string]domain.Evaluation
}

func NewInMemoryEvaluationRepo() *InMemoryEvaluationRepo {
	return &InMemoryEvaluationRepo{
		evaluations: make(map[domain.Email]map[string]domain.Evaluation),
	}
}

func (r *InMemoryEvaluationRepo) Save(evaluation domain.Evaluation) error {
	email := evaluation.StudentEmail()
	if _, ok := r.evaluations[email]; !ok {
		r.evaluations[email] = make(map[string]domain.Evaluation)
	}
	r.evaluations[email][evaluation.CriterionName()] = evaluation
	return nil
}

func (r *InMemoryEvaluationRepo) FindByStudent(email domain.Email) ([]domain.Evaluation, error) {
	evalsMap, ok := r.evaluations[email]
	if !ok {
		return nil, nil
	}
	var evaluations []domain.Evaluation
	for _, eval := range evalsMap {
		evaluations = append(evaluations, eval)
	}
	return evaluations, nil
}

func (r *InMemoryEvaluationRepo) GetAll() ([]domain.Evaluation, error) {
	var evaluations []domain.Evaluation
	for _, evalsMap := range r.evaluations {
		for _, eval := range evalsMap {
			evaluations = append(evaluations, eval)
		}
	}
	return evaluations, nil
}

func (r *InMemoryEvaluationRepo) Clear() error {
	r.evaluations = make(map[domain.Email]map[string]domain.Evaluation)
	return nil
}

type InMemoryRubricRepo struct {
	rubric domain.Rubric
}

func NewInMemoryRubricRepo() *InMemoryRubricRepo {
	return &InMemoryRubricRepo{}
}

func (r *InMemoryRubricRepo) Save(rubric domain.Rubric) error {
	r.rubric = rubric
	return nil
}

func (r *InMemoryRubricRepo) Get() (domain.Rubric, error) {
	if len(r.rubric.Criteria()) == 0 {
		return domain.Rubric{}, fmt.Errorf("рубрика не найдена")
	}
	return r.rubric, nil
}

func FormatJSON(data []application.StudentResult) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

func FormatCSV(data []application.StudentResult) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	if err := writer.Write([]string{"Name", "Email", "FinalScore"}); err != nil {
		return nil, err
	}

	for _, r := range data {
		row := []string{
			r.Name,
			r.Email,
			fmt.Sprintf("%.2f", r.FinalScore),
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
