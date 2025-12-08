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
	data = strings.TrimSpace(data)

	if data == "" {
		return 0, 0, fmt.Errorf("неверный формат данных")
	}

	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных")
	}

	stepsStr := strings.TrimSpace(parts[0])

	if strings.Contains(parts[0], " ") && stepsStr != parts[0] {
		return 0, 0, fmt.Errorf("неверный формат шагов")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil || steps <= 0 {
		return 0, 0, fmt.Errorf("неверный формат шагов")
	}

	durationStr := strings.TrimSpace(parts[1])

	if strings.Contains(parts[1], " ") && durationStr != parts[1] {
		return 0, 0, fmt.Errorf("неверный формат продолжительности")
	}

	duration, err := time.ParseDuration(durationStr)
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

	speed := distanceKm / hours
	calories := (0.035*weight + (speed*speed/height)*0.029*weight) * hours

	calories *= 16.46

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories)
}
