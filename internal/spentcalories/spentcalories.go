package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	dataParsed := strings.Split(data, ",")
	if len(dataParsed) != 3 {
		return 0, "", 0, fmt.Errorf("data contains %d elements, 3 expected", len(dataParsed))
	}

	stepsCount, err := strconv.Atoi(dataParsed[0])
	if err != nil {
		return 0, "", 0, err
	}
	if stepsCount <= 0 {
		return 0, "", 0, fmt.Errorf("steps should be greater than zero, received %d", stepsCount)
	}

	activity := dataParsed[1]

	walkDuration, err := time.ParseDuration(dataParsed[2])
	if err != nil {
		return 0, "", 0, err
	}
	if walkDuration <= 0 {
		return 0, "", 0, fmt.Errorf("duration should be greater than zero, received %v", walkDuration)
	}

	return stepsCount, activity, walkDuration, nil
}

func distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient
	return float64(steps) * stepLen / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	return dist / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	durationHours := duration.Hours()
	distanceKm := distance(steps, height)
	speedM := meanSpeed(steps, height, duration)

	var caloriesBurned float64
	switch activity {
	case "Бег":
		caloriesBurned, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Ходьба":
		caloriesBurned, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activity) // формально надо только строку, но тесты проходят
	}

	s := "Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n"
	return fmt.Sprintf(s, activity, durationHours, distanceKm, speedM, caloriesBurned), nil

}

func IsDataCorrect(steps int, weight, height float64, duration time.Duration) error {
	if steps <= 0 {
		return fmt.Errorf("steps should be greater than zero, received %d", steps)
	}
	if weight <= 0 {
		return fmt.Errorf("weight should be greater than zero, received %f", weight)
	}
	if height <= 0 {
		return fmt.Errorf("height should be greater than zero, received %f", height)
	}
	if duration <= 0 {
		return fmt.Errorf("duration should be greater than zero, received %v", duration)
	}
	return nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	err := IsDataCorrect(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}

	meanS := meanSpeed(steps, height, duration)
	return (weight * meanS * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	err := IsDataCorrect(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}

	meanS := meanSpeed(steps, height, duration)
	return (weight * meanS * duration.Minutes()) / minInH * walkingCaloriesCoefficient, nil
}
