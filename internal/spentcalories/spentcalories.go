package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	mInKm                      = 1000.0 // метров в километре
	minInH                     = 60.0   // минут в часе
	stepLengthCoefficient      = 0.45   // коэффициент длины шага
	walkingCaloriesCoefficient = 0.5    // коэффициент калорий при ходьбе
)

// parseTraining парсит строку с данными тренировки
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных: ожидается 'шаги,тип,длительность'")
	}

	// Парсим шаги
	stepsStr := strings.TrimSpace(parts[0])
	if stepsStr == "" {
		return 0, "", 0, errors.New("отсутствует количество шагов")
	}

	// Убираем возможный плюс в начале
	stepsStr = strings.TrimPrefix(stepsStr, "+")

	// Проверяем на пробелы внутри числа - это ошибка
	if strings.ContainsAny(stepsStr, " \t\n") {
		return 0, "", 0, errors.New("неверный формат количества шагов: пробелы в числе")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверный формат количества шагов: %w", err)
	}

	// ВАЖНО: здесь мы НЕ проверяем steps на положительность!
	// Это проверяется позже в функциях TrainingInfo, RunningSpentCalories, WalkingSpentCalories

	// Парсим тип тренировки
	trainingType := strings.TrimSpace(parts[1])
	if trainingType == "" {
		return 0, "", 0, errors.New("отсутствует тип тренировки")
	}

	// Парсим продолжительность
	durationStr := strings.TrimSpace(parts[2])
	if durationStr == "" {
		return 0, "", 0, errors.New("отсутствует продолжительность")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверный формат продолжительности: %w", err)
	}

	// ВАЖНО: здесь мы НЕ проверяем duration на положительность!
	// Это проверяется позже в функциях TrainingInfo, RunningSpentCalories, WalkingSpentCalories

	return steps, trainingType, duration, nil
}

// distance вычисляет дистанцию в километрах
func distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 {
		return 0
	}

	// Вычисляем длину шага
	stepLength := height * stepLengthCoefficient

	// Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// Переводим в километры
	return distanceMeters / mInKm
}

// meanSpeed вычисляет среднюю скорость
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 || steps <= 0 || height <= 0 {
		return 0
	}

	// Вычисляем дистанцию
	dist := distance(steps, height)

	// Переводим продолжительность в часы
	durationHours := duration.Hours()

	if durationHours == 0 {
		return 0
	}

	// Вычисляем среднюю скорость
	return dist / durationHours
}

// RunningSpentCalories вычисляет потраченные калории при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка входных параметров
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

	// Вычисляем среднюю скорость
	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0, errors.New("невозможно рассчитать скорость")
	}

	// Переводим продолжительность в минуты
	durationMinutes := duration.Minutes()

	// Вычисляем калории по формуле: (вес * скорость * продолжительность в минутах) / минут в часе
	calories := (weight * speed * durationMinutes) / minInH

	return calories, nil
}

// WalkingSpentCalories вычисляет потраченные калории при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка входных параметров
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

	// Вычисляем среднюю скорость
	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0, errors.New("невозможно рассчитать скорость")
	}

	// Переводим продолжительность в минуты
	durationMinutes := duration.Minutes()

	// Вычисляем базовые калории по формуле: (вес * скорость * продолжительность в минутах) / минут в часе
	baseCalories := (weight * speed * durationMinutes) / minInH

	// Применяем корректирующий коэффициент для ходьбы
	calories := baseCalories * walkingCaloriesCoefficient

	return calories, nil
}

// TrainingInfo формирует информацию о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Проверяем корректность данных
	if steps <= 0 {
		return "", errors.New("количество шагов должно быть положительным")
	}
	if duration <= 0 {
		return "", errors.New("продолжительность должна быть положительной")
	}

	var calories float64
	var calcErr error

	// В зависимости от типа тренировки вычисляем калории
	switch trainingType {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, calcErr = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", calcErr
	}

	// Вычисляем дистанцию и скорость
	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	durationHours := duration.Hours()

	// Форматируем результат
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		trainingType, durationHours, dist, speed, calories), nil
}
