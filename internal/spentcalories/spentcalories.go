package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	lenStep               = 0.65
	mInKm                 = 1000
	stepLengthCoefficient = 0.45
)

func parseTraining(data string) (int, string, time.Duration, error) {
	data = strings.TrimSpace(data)
	if data == "" {
		return 0, "", 0, errors.New("неверный формат данных")
	}

	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных")
	}

	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil || steps <= 0 {
		return 0, "", 0, errors.New("неверный формат шагов")
	}

	trainingType := strings.TrimSpace(parts[1])
	if trainingType != "Бег" && trainingType != "Ходьба" {
		return 0, "", 0, errors.New("неизвестный тип тренировки")
	}

	durationStr := strings.TrimSpace(parts[2])
	duration, err := time.ParseDuration(durationStr)
	if err != nil || duration <= 0 {
		return 0, "", 0, errors.New("неверный формат продолжительности")
	}

	return steps, trainingType, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	if stepLength < lenStep {
		stepLength = lenStep
	}
	return (float64(steps) * stepLength) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	return dist / duration.Hours()
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть положительной")
	}

	speed := meanSpeed(steps, height, duration)

	calories := (0.035*weight + (speed*speed/height)*0.029*weight) * duration.Hours() * 16.46

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть положительной")
	}

	speed := meanSpeed(steps, height, duration)

	calories := (0.035*weight + (speed*speed/height)*0.029*weight) * duration.Hours() * 8.23

	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	if weight <= 0 || height <= 0 {
		return "", errors.New("некорректные параметры пользователя")
	}

	var calories float64
	switch trainingType {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	durationHours := duration.Hours()

	result := fmt.Sprintf("Тип тренировки: %s\n", trainingType)
	result += fmt.Sprintf("Длительность: %.2f ч.\n", durationHours)
	result += fmt.Sprintf("Дистанция: %.2f км.\n", dist)
	result += fmt.Sprintf("Скорость: %.2f км/ч\n", speed)
	result += fmt.Sprintf("Сожгли калорий: %.2f\n", calories)

	return result, nil
}
