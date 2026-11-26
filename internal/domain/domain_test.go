package domain

import "testing"

func TestNewEmail(t *testing.T) {
	if _, err := NewEmail("anna@edu.ru"); err != nil {
		t.Errorf("ожидали валидный email, получили ошибку: %v", err)
	}
	if _, err := NewEmail("invalid-email"); err == nil {
		t.Error("ожидали ошибку для некорректного email, но получили nil")
	}
}

func TestRubricBuilder(t *testing.T) {
	w1, _ := NewWeight(0.4)
	w2, _ := NewWeight(0.3) // сумма 0.7 != 1.0

	var rb RubricBuilder
	_ = rb.AddCriterion("Домашки", 20, w1)
	_ = rb.AddCriterion("Экзамен", 10, w2)

	if _, err := rb.Build(); err == nil {
		t.Error("ожидали ошибку при сумме весов не равной 1.0")
	}
}

func TestCalculateFinalScore(t *testing.T) {
	w1, _ := NewWeight(0.4)
	w2, _ := NewWeight(0.6)

	var rb RubricBuilder
	_ = rb.AddCriterion("Домашки", 20, w1)
	_ = rb.AddCriterion("Экзамен", 10, w2)
	rubric, _ := rb.Build()

	email, _ := NewEmail("anna@edu.ru")
	c1 := rubric.Criteria()[0]
	c2 := rubric.Criteria()[1]

	// 1. Без оценок итог должен быть 0
	emptyScore, err := CalculateFinalScore(rubric, nil)
	if err != nil || emptyScore != 0 {
		t.Errorf("ожидали 0 без оценок, получили %v (err: %v)", emptyScore, err)
	}

	// 2. Половина баллов по каждому критерию должна дать ровно половину от максимума
	evalHalf1, _ := NewEvaluation(email, c1, 10)
	evalHalf2, _ := NewEvaluation(email, c2, 5)
	halfScore, _ := CalculateFinalScore(rubric, []Evaluation{evalHalf1, evalHalf2})

	evalFull1, _ := NewEvaluation(email, c1, 20)
	evalFull2, _ := NewEvaluation(email, c2, 10)
	fullScore, _ := CalculateFinalScore(rubric, []Evaluation{evalFull1, evalFull2})

	if fullScore <= 0 || fullScore != halfScore*2 {
		t.Errorf("неверный расчёт весов: half=%v, full=%v", halfScore, fullScore)
	}
}

func TestNewEvaluation(t *testing.T) {
	w, _ := NewWeight(1.0)
	var rb RubricBuilder
	_ = rb.AddCriterion("Экзамен", 10, w)
	rubric, _ := rb.Build()

	email, _ := NewEmail("anna@edu.ru")
	criterion := rubric.Criteria()[0]

	// 1. Валидные граничные значения (0 и максимум)
	if _, err := NewEvaluation(email, criterion, 0); err != nil {
		t.Errorf("ожидали валидную оценку 0, получили ошибку: %v", err)
	}
	if _, err := NewEvaluation(email, criterion, 10); err != nil {
		t.Errorf("ожидали валидную оценку 10, получили ошибку: %v", err)
	}

	// 2. Отрицательный балл
	if _, err := NewEvaluation(email, criterion, -1); err == nil {
		t.Error("ожидали ошибку при отрицательном балле, но получили nil")
	}

	// 3. Превышение максимума
	if _, err := NewEvaluation(email, criterion, 11); err == nil {
		t.Error("ожидали ошибку при превышении максимума, но получили nil")
	}
}
