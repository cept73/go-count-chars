/**
 * Создаем файл наполненный псевдослучайным контентом (a-z и пробелы)
 * Затем считаем количество вхождений букв и выводим на экран
 *
 * Сравнение времени подсчета (ms):
 * - С корутинами 133 ms
 * - Без корутин 450 ms
 * - PHP 300 ms
 */

package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"
)

const FILE_NAME = "data.txt"
const BUFFER_SIZE = 4096
const FILE_LENGTH = 1024768 * 10 // 10MB

/**
 * Write random content to file
 */
func writeToFile() {
	file, err := os.OpenFile(FILE_NAME, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}

	defer file.Close()
	defer fmt.Println("End write:", file.Name())

	var data = make([]byte, 1)
	charA := int('a')
	charsAlphaCount := int('z') - charA
	for i := 0; i < FILE_LENGTH; i++ {
		charOffset := int(rand.Uint32()) % (charsAlphaCount + 1)
		if charOffset == charsAlphaCount {
			file.WriteString(" ")
			continue
		}
		data[0] = byte(charA + charOffset)
		file.Write(data)
	}
}

/**
 * Read file with parallel running
 */
func readFromFileWithCoroutines() {
	file, err := os.Open(FILE_NAME)
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}

	defer file.Close()

	var result = map[string]uint32{}
	var Wg sync.WaitGroup
	var mutex sync.Mutex

	for {
		data := make([]byte, BUFFER_SIZE)

		loadedLen, err := file.Read(data)
		if err == io.EOF {
			break
		}

		Wg.Add(1)
		go func(data []byte) {
			defer Wg.Done()

			localResult := map[string]uint32{}

			for index, element := range data {
				if index >= loadedLen {
					break
				}

				key := string(element)
				localResult[key]++
			}

			mutex.Lock()
			defer mutex.Unlock()
			for char, count := range localResult {
				result[char] += count
			}
		}(data)
	}

	Wg.Wait()

	fmt.Println(result)
	printResultOverall(result)
}

/**
 * Read file without parallel running
 */
func readFromFileWithoutCoroutiens() {
	file, err := os.Open(FILE_NAME)
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}

	defer file.Close()

	var result = map[string]uint32{}

	for {
		data := make([]byte, BUFFER_SIZE)
		readLen, err := file.Read(data)
		if err == io.EOF {
			break
		}

		if readLen > 0 {
			for i := 0; i < readLen; i++ {
				char := string(data[i])
				result[char]++
			}
		}
	}

	fmt.Println(result)
	printResultOverall(result)
}

/**
 * Show overall count
 */
func printResultOverall(result map[string]uint32) {
	var overallCount uint32 = 0
	for _, value := range result {
		overallCount += value
	}
	fmt.Println("Overall:", overallCount)
}

/**
 * Call function and calculate execution time
 */
func callAndShowDuration(message string, function func()) {
	fmt.Println(message)
	startTimer := time.Now()
	function()
	fmt.Println(time.Since(startTimer))
}

// Main function
func main() {
	// If file doesn't exists, create and fill
	_, err := os.Stat(FILE_NAME)
	if err != nil {
		fmt.Println("Create test file with length", FILE_LENGTH)
		writeToFile()
	}

	fmt.Println("Read test file")

	callAndShowDuration("With coroutines", readFromFileWithCoroutines)
	fmt.Println()
	callAndShowDuration("Without coroutines", readFromFileWithoutCoroutiens)
}
