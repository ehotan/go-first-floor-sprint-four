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
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	dataSlices := strings.Split(data, ",")
	if len(dataSlices) != 2 {
		return 0, 0, errors.New("ошибка ввода данных, дневная активность: неверное количество параметров")
	}
	stepsCount, err := strconv.Atoi(dataSlices[0])
	if err != nil {
		return 0, 0, err
	}
	if stepsCount <= 0 {
		return 0, 0, errors.New("ошибка ввода данных, дневная активность: количество шагов должно быть больше нуля")
	}
	dayActionsDuration, err := time.ParseDuration(dataSlices[1])
	if err != nil {
		return 0, 0, err
	}
	if dayActionsDuration <= 0 {
		return 0, 0, errors.New("ошибка ввода данных, дневная активность: продолжительность должна быть больше нуля")
	}
	return stepsCount, dayActionsDuration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	stepsCount, dayActionsDuration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	if stepsCount <= 0 {
		return ""
	}
	distanceM := float64(stepsCount) * stepLength
	distanceKm := distanceM / mInKm
	calories, _ := spentcalories.WalkingSpentCalories(stepsCount, weight, height, dayActionsDuration)
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", stepsCount, distanceKm, calories)
}
