package adapters

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"grading-system/internal/application"
	"grading-system/internal/domain"
)

type ConsoleApp struct {
	facade  *application.AssessmentFacade
	scanner *bufio.Scanner
}

func NewConsoleApp(facade *application.AssessmentFacade) *ConsoleApp {
	return &ConsoleApp{
		facade:  facade,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (c *ConsoleApp) Run() {
	fmt.Println("=== Модуль ПОС: Преподавательская оценка студентов ===")

	for {
		c.printMenu()
		choice, ok := c.readLine("\nВыберите пункт: ")
		if !ok {
			fmt.Println("\nЗавершение работы.")
			return
		}

		switch choice {
		case "1":
			c.runDemoScenario()
		case "2":
			c.handleShowRubric()
		case "3":
			c.handleSetRubric()
		case "4":
			c.handleShowStudents()
		case "5":
			c.handleRegisterStudent()
		case "6":
			c.handleGradeStudent()
		case "7":
			c.handleShowStudentEvaluations()
		case "8":
			c.handleShowResults()
		case "9":
			c.handleExport("results.csv", FormatCSV)
		case "10":
			c.handleExport("results.json", FormatJSON)
		case "11":
			c.handleClearEvaluations()
		case "0":
			fmt.Println("До свидания!")
			return
		default:
			fmt.Println("Неизвестная команда, введите число от 0 до 10.")
		}
	}
}

func (c *ConsoleApp) printMenu() {
	fmt.Println("\n--- Меню ---")
	fmt.Println("1.  Запустить сквозной демо-сценарий (рубрика -> студенты -> оценки -> итог -> CSV и JSON)")
	fmt.Println("2.  Показать текущую рубрику")
	fmt.Println("3.  Настроить рубрику")
	fmt.Println("4.  Показать список студентов")
	fmt.Println("5.  Добавить студента")
	fmt.Println("6.  Выставить / исправить оценку")
	fmt.Println("7.  Показать оценки студента")
	fmt.Println("8.  Показать итоговую ведомость")
	fmt.Println("9.  Экспортировать ведомость в CSV")
	fmt.Println("10. Экспортировать ведомость в JSON")
	fmt.Println("11. Очистить все оценки")
	fmt.Println("0.  Выход")
}

func (c *ConsoleApp) runDemoScenario() {
	fmt.Println("\n[Демо] Настраиваем рубрику через RubricBuilder...")
	w1, err := domain.NewWeight(0.4)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	w2, err := domain.NewWeight(0.6)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	var rb domain.RubricBuilder
	if err := rb.AddCriterion("Домашки", 20, w1); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	if err := rb.AddCriterion("Экзамен", 10, w2); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	rubric, err := rb.Build()
	if err != nil {
		fmt.Printf("Ошибка сборки рубрики: %v\n", err)
		return
	}
	if err := c.facade.SetRubric(rubric); err != nil {
		fmt.Printf("Ошибка сохранения рубрики: %v\n", err)
		return
	}

	fmt.Println("[Демо] Регистрируем студентов и выставляем оценки...")
	demoSteps := []func() error{
		func() error { return c.facade.RegisterStudent("Анна Смирнова", "anna@edu.ru") },
		func() error { return c.facade.RegisterStudent("Иван Петров", "ivan@edu.ru") },
		func() error { return c.facade.GradeStudent("anna@edu.ru", "Домашки", 18) },
		func() error { return c.facade.GradeStudent("anna@edu.ru", "Экзамен", 9) },
		func() error { return c.facade.GradeStudent("ivan@edu.ru", "Домашки", 15) },
	}
	for _, step := range demoSteps {
		if err := step(); err != nil {
			fmt.Printf("Ошибка демо-сценария: %v\n", err)
			return
		}
	}

	fmt.Println("\n[Демо] Итоговая ведомость:")
	c.handleShowResults()

	fmt.Println("\n[Демо] Экспорт через ExportPipeline:")
	c.handleExport("results.csv", FormatCSV)
	c.handleExport("results.json", FormatJSON)
}

func (c *ConsoleApp) handleShowRubric() {
	rubric, err := c.facade.GetRubric()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	PrintRubric(rubric)
}

func (c *ConsoleApp) handleSetRubric() {
	var rb domain.RubricBuilder
	totalWeight := 0.0

	fmt.Println("Добавление критериев (сумма весов должна составить ровно 1.0).")
	for {
		fmt.Printf("Текущая сумма весов: %.2f / 1.00\n", totalWeight)
		name, ok := c.readLine("Название критерия (пустая строка — завершить сборку): ")
		if !ok || name == "" {
			break
		}

		maxScoreStr, ok := c.readLine("Максимальный балл (целое число > 0): ")
		if !ok {
			return
		}
		maxScore, err := strconv.Atoi(maxScoreStr)
		if err != nil {
			fmt.Println("Ошибка: введите целое число.")
			continue
		}

		weightStr, ok := c.readLine("Вес критерия (например, 0.4): ")
		if !ok {
			return
		}
		weightVal, err := strconv.ParseFloat(strings.ReplaceAll(weightStr, ",", "."), 64)
		if err != nil {
			fmt.Println("Ошибка: введите дробное число.")
			continue
		}

		weight, err := domain.NewWeight(weightVal)
		if err != nil {
			fmt.Printf("Ошибка домена: %v\n", err)
			continue
		}

		if err := rb.AddCriterion(name, maxScore, weight); err != nil {
			fmt.Printf("Ошибка домена: %v\n", err)
			continue
		}
		totalWeight += weightVal
		fmt.Printf("Критерий %q добавлен.\n", name)
	}

	rubric, err := rb.Build()
	if err != nil {
		fmt.Printf("Не удалось собрать рубрику: %v\n", err)
		return
	}
	if err := c.facade.SetRubric(rubric); err != nil {
		fmt.Printf("Ошибка сохранения рубрики: %v\n", err)
		return
	}
	fmt.Println("Рубрика успешно сохранена.")
}

func (c *ConsoleApp) handleShowStudents() {
	students, err := c.facade.GetStudents()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	if len(students) == 0 {
		fmt.Println("Список студентов пуст.")
		return
	}
	PrintStudents(students)
}

func (c *ConsoleApp) handleRegisterStudent() {
	name, ok := c.readLine("Имя студента: ")
	if !ok {
		return
	}
	email, ok := c.readLine("Email студента: ")
	if !ok {
		return
	}

	if err := c.facade.RegisterStudent(name, email); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	fmt.Println("Студент успешно добавлен.")
}

func (c *ConsoleApp) handleGradeStudent() {
	email, ok := c.readLine("Email студента: ")
	if !ok {
		return
	}
	criterion, ok := c.readLine("Название критерия: ")
	if !ok {
		return
	}
	scoreStr, ok := c.readLine("Балл (целое число): ")
	if !ok {
		return
	}
	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		fmt.Println("Ошибка: балл должен быть целым числом.")
		return
	}

	if err := c.facade.GradeStudent(email, criterion, score); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	fmt.Println("Оценка сохранена.")
}

func (c *ConsoleApp) handleShowStudentEvaluations() {
	email, ok := c.readLine("Email студента: ")
	if !ok {
		return
	}

	evals, err := c.facade.GetStudentEvaluations(email)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	if len(evals) == 0 {
		fmt.Println("У этого студента пока нет оценок.")
		return
	}
	PrintEvaluations(evals)
}

func (c *ConsoleApp) handleShowResults() {
	results, err := c.facade.GetAllResults()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	if len(results) == 0 {
		fmt.Println("Нет студентов для расчёта ведомости.")
		return
	}
	PrintTable(results)
}

func (c *ConsoleApp) handleExport(filePath string, format func([]application.StudentResult) ([]byte, error)) {
	pipeline := ExportPipeline{
		FilePath: filePath,
		Fetch:    c.facade.GetAllResults,
		Format:   format,
	}
	if err := pipeline.Run(); err != nil {
		fmt.Printf("Ошибка экспорта в %s: %v\n", filePath, err)
		return
	}
	fmt.Printf("Ведомость сохранена в %s\n", filePath)
}

func (c *ConsoleApp) readLine(prompt string) (string, bool) {
	fmt.Print(prompt)
	if !c.scanner.Scan() {
		return "", false
	}
	return strings.TrimSpace(c.scanner.Text()), true
}

func PrintTable(results []application.StudentResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "СТУДЕНТ\tEMAIL\tИТОГОВЫЙ БАЛЛ")
	fmt.Fprintln(w, "-------\t-----\t-------------")
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\t%.2f\n", r.Name, r.Email, r.FinalScore)
	}
	w.Flush()
}

func PrintStudents(students []domain.Student) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "СТУДЕНТ\tEMAIL")
	fmt.Fprintln(w, "-------\t-----")
	for _, s := range students {
		fmt.Fprintf(w, "%s\t%s\n", s.Name(), s.Email())
	}
	w.Flush()
}

func PrintRubric(rubric domain.Rubric) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "КРИТЕРИЙ\tМАКС. БАЛЛ\tВЕС")
	fmt.Fprintln(w, "--------\t----------\t---")
	for _, cr := range rubric.Criteria() {
		fmt.Fprintf(w, "%s\t%d\t%.2f\n", cr.Name(), cr.MaxScore(), float64(cr.Weight()))
	}
	w.Flush()
}

func PrintEvaluations(evals []domain.Evaluation) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "КРИТЕРИЙ\tБАЛЛ")
	fmt.Fprintln(w, "--------\t----")
	for _, e := range evals {
		fmt.Fprintf(w, "%s\t%d\n", e.CriterionName(), e.Score())
	}
	w.Flush()
}

func (c *ConsoleApp) confirm(prompt string) bool {
	ans, ok := c.readLine(prompt + " (yes/no): ")
	if !ok {
		return false
	}
	ans = strings.ToLower(ans)
	return ans == "yes" || ans == "y"
}

func (c *ConsoleApp) handleClearEvaluations() {
	if !c.confirm("Точно стереть все выставленные оценки?") {
		fmt.Println("Отменено.")
		return
	}
	if err := c.facade.ClearEvaluations(); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	fmt.Println("Все оценки очищены.")
}
