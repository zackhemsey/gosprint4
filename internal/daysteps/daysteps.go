package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	mInKm      = 1000.0 // метров в километре
	stepLength = 0.65   // длина шага в метрах
)

// parsePackage парсит строку с данными активности
func parsePackage(data string) (int, time.Duration, error) {
	if data == "" {
		return 0, 0, errors.New("пустая строка данных")
	}

	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат данных: ожидается 'шаги,длительность'")
	}

	// Парсим шаги
	stepsStr := strings.TrimSpace(parts[0])
	if stepsStr == "" {
		return 0, 0, errors.New("отсутствует количество шагов")
	}

	// Убираем возможный плюс в начале, но проверяем пробелы
	stepsStr = strings.TrimPrefix(stepsStr, "+")

	// Проверяем на пробелы внутри числа - это ошибка
	if strings.ContainsAny(stepsStr, " \t\n") {
		return 0, 0, errors.New("неверный формат количества шагов: пробелы в числе")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат количества шагов: %w", err)
	}

	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть положительным")
	}

	// Парсим продолжительность
	durationStr := strings.TrimSpace(parts[1])
	if durationStr == "" {
		return 0, 0, errors.New("отсутствует продолжительность")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат продолжительности: %w", err)
	}

	if duration <= 0 {
		return 0, 0, errors.New("продолжительность должна быть положительной")
	}

	return steps, duration, nil
}

// DayActionInfo формирует информацию о дневной активности
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	// Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// Переводим в километры
	distanceKm := distanceMeters / mInKm

	// Вычисляем калории
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	// Форматируем результат
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories)
}
