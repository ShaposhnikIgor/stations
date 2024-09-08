package generator

import (
	"fmt"
	"os"
	"strconv"
)

func F_Gener(filePath string) {
	// Создание файла для записи данных
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer file.Close()

	// Количество станций
	numStations := 10000

	// Запись станций в файл
	_, err = file.WriteString("stations:\n")
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return
	}

	for i := 1; i <= numStations; i++ {
		line := "st" + strconv.Itoa(i) + "," + strconv.Itoa(i) + "," + strconv.Itoa(i) + "\n"
		_, err = file.WriteString(line)
		if err != nil {
			fmt.Println("Ошибка при записи в файл:", err)
			return
		}
	}

	// Запись соединений в файл
	_, err = file.WriteString("connections:\n")
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return
	}

	// Мапа для хранения уже созданных соединений
	connectionsMap := make(map[string]bool)

	for i := 1; i < numStations; i++ {
		// Создаем соединение между текущей и следующей станцией
		connection1 := "st" + strconv.Itoa(i) + "-st" + strconv.Itoa(i+1)
		// Проверяем, существует ли уже такое соединение (либо в прямом, либо в обратном направлении)
		if !connectionsMap[connection1] && !connectionsMap["st"+strconv.Itoa(i+1)+"-st"+strconv.Itoa(i)] {
			_, err = file.WriteString(connection1 + "\n")
			if err != nil {
				fmt.Println("Ошибка при записи в файл:", err)
				return
			}
			// Добавляем соединение в карту
			connectionsMap[connection1] = true
		}

		// Создаем соединение между текущей и случайной станцией (для примера - i%100 + 1)
		connection2 := "st" + strconv.Itoa(i) + "-st" + strconv.Itoa(i%100+1)
		// Проверяем, существует ли уже такое соединение (либо в прямом, либо в обратном направлении)
		if !connectionsMap[connection2] && !connectionsMap["st"+strconv.Itoa(i%100+1)+"-st"+strconv.Itoa(i)] {
			_, err = file.WriteString(connection2 + "\n")
			if err != nil {
				fmt.Println("Ошибка при записи в файл:", err)
				return
			}
			// Добавляем соединение в карту
			connectionsMap[connection2] = true
		}
	}

	// Добавим несколько дополнительных соединений для последней станции
	connection1 := "st" + strconv.Itoa(numStations) + "-st1"
	if !connectionsMap[connection1] && !connectionsMap["st1-st"+strconv.Itoa(numStations)] {
		_, err = file.WriteString(connection1 + "\n")
		if err != nil {
			fmt.Println("Ошибка при записи в файл:", err)
			return
		}
		connectionsMap[connection1] = true
	}

	connection2 := "st" + strconv.Itoa(numStations) + "-st" + strconv.Itoa(numStations/2)
	if !connectionsMap[connection2] && !connectionsMap["st"+strconv.Itoa(numStations/2)+"-st"+strconv.Itoa(numStations)] {
		_, err = file.WriteString(connection2 + "\n")
		if err != nil {
			fmt.Println("Ошибка при записи в файл:", err)
			return
		}
		connectionsMap[connection2] = true
	}

	fmt.Println("Данные успешно сохранены в файл", filePath)
}
