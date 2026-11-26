package domain

import (
	"errors"
	"math"
	"strings"
)

type Weight float64

func NewWeight(v float64) (Weight, error) {
	if v <= 0 || v > 1 {
		return 0, errors.New("вес критерия должен быть больше 0 и не больше 1")
	}
	return Weight(v), nil
}

type Criterion struct {
	name     string
	maxScore int
	weight   Weight
}

func (c Criterion) Name() string {
	return c.name
}

func (c Criterion) MaxScore() int {
	return c.maxScore
}

func (c Criterion) Weight() Weight {
	return c.weight
}

func NewCriterion(name string, maxScore int, weight Weight) (Criterion, error) {
	if name == "" {
		return Criterion{}, errors.New("название критерия не может быть пустым")
	}
	if maxScore <= 0 {
		return Criterion{}, errors.New("максимальный балл должен быть больше 0")
	}
	return Criterion{name: name, maxScore: maxScore, weight: weight}, nil
}

type Rubric struct {
	criteria []Criterion
}

type RubricBuilder struct {
	criteria []Criterion
}

func (r Rubric) Criteria() []Criterion {
	res := make([]Criterion, len(r.criteria))
	copy(res, r.criteria)
	return res
}

func (b *RubricBuilder) AddCriterion(name string, maxScore int, weight Weight) error {
	criterion, err := NewCriterion(name, maxScore, weight)
	if err != nil {
		return err
	}

	for _, existingCriterion := range b.criteria {
		if existingCriterion.Name() == name {
			return errors.New("критерий с таким названием уже существует")
		}
	}
	b.criteria = append(b.criteria, criterion)
	return nil
}

func (b *RubricBuilder) Build() (Rubric, error) {
	if len(b.criteria) == 0 {
		return Rubric{}, errors.New("нельзя создать рубрику без критериев")
	}

	sum := 0.0
	for _, criterion := range b.criteria {
		sum += float64(criterion.Weight())
	}
	if math.Abs(sum-1.0) > 1e-9 {
		return Rubric{}, errors.New("сумма весов всех критериев должна быть равна 1")
	}

	criteriaCopy := make([]Criterion, len(b.criteria))
	copy(criteriaCopy, b.criteria)
	return Rubric{criteria: criteriaCopy}, nil
}

type Email string

func NewEmail(email string) (Email, error) {
	if email == "" {
		return "", errors.New("email не может быть пустым")
	} else if !strings.Contains(email, "@") {
		return "", errors.New("email должен содержать символ '@'")
	}
	return Email(email), nil
}

type Student struct {
	name  string
	email Email
}

func (s Student) Name() string {
	return s.name
}

func (s Student) Email() Email {
	return s.email
}

func NewStudent(name string, email Email) (Student, error) {
	if name == "" {
		return Student{}, errors.New("имя студента не может быть пустым")
	}
	return Student{name: name, email: email}, nil
}

type Evaluation struct {
	studentEmail  Email
	criterionName string
	score         int
}

func (e Evaluation) StudentEmail() Email {
	return e.studentEmail
}

func (e Evaluation) CriterionName() string {
	return e.criterionName
}

func (e Evaluation) Score() int {
	return e.score
}

func NewEvaluation(studentEmail Email, criterion Criterion, score int) (Evaluation, error) {
	if criterion.Name() == "" {
		return Evaluation{}, errors.New("название критерия не может быть пустым")
	}
	if score < 0 || score > criterion.MaxScore() {
		return Evaluation{}, errors.New("балл не может быть отрицательным или превышать максимальный балл")
	}
	return Evaluation{
		studentEmail:  studentEmail,
		criterionName: criterion.Name(),
		score:         score,
	}, nil
}
