package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	stepLength = 0.65
	mInKm      = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных")
	}

	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || steps <= 0 {
		return 0, 0, fmt.Errorf("неверный формат шагов")
	}

	duration, err := time.ParseDuration(strings.TrimSpace(parts[1]))
	if err != nil || duration <= 0 {
		return 0, 0, fmt.Errorf("неверный формат продолжительности")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		return fmt.Sprintf("Ошибка: %v", err)
	}

	if weight <= 0 || height <= 0 {
		return "Ошибка: некорректные параметры пользователя"
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	hours := duration.Hours()
	calories := 0.0
	if hours > 0 {

		calories = 0.5 * distanceKm * weight
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps, distanceKm, calories)
}
