package spentcalories

import (
	"errors"
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
	dataSlices := strings.Split(data, ",")
	if len(dataSlices) != 3 {
		err := errors.New("ошибка ввода данных, тренировки: неверное количество параметров")
		return 0, "", 0, err
	}
	stepsCount, err := strconv.Atoi(dataSlices[0])
	if err != nil {
		return 0, "", 0, err
	}
	if stepsCount <= 0 {
		return 0, "", 0, errors.New("ошибка ввода данных, тренировки: количество шагов должно быть положительно")
	}
	trainigDuration, err := time.ParseDuration(dataSlices[2])
	if err != nil {
		return 0, "", 0, err
	}
	if trainigDuration <= 0 {
		return 0, "", 0, errors.New("ошибка ввода данных, тренировки: продолжительность должна быть положительной")
	}
	return stepsCount, dataSlices[1], trainigDuration, nil
}

func distance(steps int, height float64) float64 {
	oneStepLen := height * stepLengthCoefficient
	distanceM := oneStepLen * float64(steps)
	return distanceM / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	stepsCount, trainingType, trainingDuration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
	}
	var trainingOutput string
	switch trainingType {
	case "Ходьба":
		walkingDistance := distance(stepsCount, height)
		walkingMeanSpeed := meanSpeed(stepsCount, height, trainingDuration)
		calories, err := WalkingSpentCalories(stepsCount, weight, height, trainingDuration)
		if err != nil {
			return "", err
		}
		trainingOutput = fmt.Sprintf("Тип тренировки: Ходьба\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			trainingDuration.Hours(), walkingDistance, walkingMeanSpeed, calories)
	case "Бег":
		runningDistance := distance(stepsCount, height)
		runningMeanSpeed := meanSpeed(stepsCount, height, trainingDuration)
		calories, err := RunningSpentCalories(stepsCount, weight, height, trainingDuration)
		if err != nil {
			return "", err
		}
		trainingOutput = fmt.Sprintf("Тип тренировки: Бег\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			trainingDuration.Hours(), runningDistance, runningMeanSpeed, calories)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	return trainingOutput, err
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if duration <= 0 {
		return 0, errors.New("ошибка ввода данных, тренировки: продолжительность должна быть положительной")
	}
	if steps <= 0 {
		return 0, errors.New("ошибка ввода данных, тренировки: количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, errors.New("ошибка ввода данных, тренировки: вес должен быть положительным")
	}
	if height <= 0 {
		return 0, errors.New("ошибка ввода данных, тренировки: рост должен быть положительным")
	}
	runningMeanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	return (weight * runningMeanSpeed * durationInMinutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: обработать ошибки, разбить мат. действия на строки
	if duration <= 0 {
		return 0, errors.New("ошибка ввода данных, тренировки: продолжительность должна быть положительной")
	}
	if steps <= 0 {
		return 0, errors.New("ошибка ввода данных, тренировки: количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, errors.New("ошибка ввода данных, тренировки: вес должен быть положительным")
	}
	if height <= 0 {
		return 0, errors.New("ошибка ввода данных, тренировки: рост должен быть положительным")
	}
	walkingMeanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	return ((weight * walkingMeanSpeed * durationInMinutes) / minInH) * walkingCaloriesCoefficient, nil
}
