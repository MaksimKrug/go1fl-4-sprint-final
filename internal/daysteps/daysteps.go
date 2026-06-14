package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	dataParsed := strings.Split(data, ",")
	if len(dataParsed) != 2 {
		return 0, 0, fmt.Errorf("data contatins %d elements, 2 expected", len(dataParsed))
	}

	stepsCount, err := strconv.Atoi(dataParsed[0])
	if err != nil {
		return 0, 0, err
	}
	if stepsCount <= 0 {
		return 0, 0, fmt.Errorf("steps should be greater than zero, received %d", stepsCount)
	}

	walkDuration, err := time.ParseDuration(dataParsed[1])
	if walkDuration <= 0 {
		return 0, 0, fmt.Errorf("duration should be greater than zero, received %v", walkDuration)
	}
	if err != nil {
		return 0, 0, err
	}

	return stepsCount, walkDuration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	stepsCount, walkDuration, err := parsePackage(data)
	if err != nil {
		log.Println("Error:", err)
		return ""
	}

	// 2. Проверить, чтобы количество шагов было больше 0. В противном случае вернуть пустую строку.
	// здесь нет смысла, т.к. parsePackage вернёт ошибку

	distanceM := float64(stepsCount) * stepLength
	distanceKm := distanceM / mInKm
	spentCalories, err := spentcalories.WalkingSpentCalories(stepsCount, weight, height, walkDuration)
	if err != nil {
		log.Println("Error:", err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", stepsCount, distanceKm, spentCalories)
}
