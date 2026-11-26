package domain

type StudentRepository interface {
	Save(student Student) error
	FindByEmail(email Email) (Student, error)
	GetAll() ([]Student, error)
}

type EvaluationRepository interface {
	Save(e Evaluation) error
	FindByStudent(email Email) ([]Evaluation, error)
	GetAll() ([]Evaluation, error)
	Clear() error
}

type RubricRepository interface {
	Save(rubric Rubric) error
	Get() (Rubric, error)
}
