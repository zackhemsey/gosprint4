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

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid data format: expected 'steps,duration'")
	}

	// Парсим шаги сразу
	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid steps: %w", err)
	}

	if steps <= 0 {
		return 0, 0, errors.New("steps must be positive")
	}

	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть положительным")
	}

	// Парсим продолжительность
	duration, err := time.ParseDuration(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid duration format: %w", err)
	}

	if duration <= 0 {
		return 0, 0, errors.New("duration must be positive")
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
