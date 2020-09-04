package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const (
	pageTop = `<!DOCTYPE HTML><html><head>
<style>.error{color:#FF0000;}</style></head><title>Статистика</title>
<body><h3>Статистика</h3>
<p>Вычисляет базовые статистические величины для данного ряда чисел</p>`
	form = `<form action="/" method="POST">
<label for="numbers">Числа (должны быть разделены запятыми или пробелами):</label><br />
<input type="text" name="numbers" size="30"><br />
<input type="submit" value="Calculate">
</form>`
	pageBottom = `</body></html>`
	anError    = `<p class="error">%s</p>`
)

type statistics struct {
	numbers   []float64
	mean      float64
	median    float64
	mode      []float64
	deviation float64
}

func main() {
	http.HandleFunc("/", homePage)
	if err := http.ListenAndServe(":9001", nil); err != nil {
		log.Fatal("failed to start server", err)
	}
}

func homePage(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm() // Must be called before writing response
	fmt.Fprint(writer, pageTop, form)
	if err != nil {
		fmt.Fprintf(writer, anError, err)
	} else {
		if numbers, message, ok := processRequest(request); ok {
			stats := getStats(numbers)
			fmt.Fprint(writer, formatStats(stats))
		} else if message != "" {
			fmt.Fprintf(writer, anError, message)
		}
	}
	fmt.Fprint(writer, pageBottom)
}

func processRequest(request *http.Request) ([]float64, string, bool) {
	var numbers []float64
	if slice, found := request.Form["numbers"]; found && len(slice) > 0 {
		text := strings.Replace(slice[0], ",", " ", -1)
		for _, field := range strings.Fields(text) {
			if x, err := strconv.ParseFloat(field, 64); err != nil {
				return numbers, "'" + field + "' is invalid", false
			} else {
				numbers = append(numbers, x)
			}
		}
	}
	if len(numbers) == 0 {
		return numbers, "", false // no data first time form is shown
	}
	return numbers, "", true
}

func formatStats(stats statistics) string {
	return fmt.Sprintf(`<table border="1">
<tr><th colspan="2">Результат</th></tr>
<tr><td>Numbers</td><td>%v</td></tr>
<tr><td>Count</td><td>%d</td></tr>
<tr><td>Mean</td><td>%f</td></tr>
<tr><td>Median</td><td>%f</td></tr>
<tr><td>Mode</td><td>%f</td></tr>
<tr><td>Стандартное отклонение</td><td>%f</td></tr>
</table>`, stats.numbers, len(stats.numbers), stats.mean, stats.median, stats.mode, stats.deviation)
}

func getStats(numbers []float64) (stats statistics) {
	stats.numbers = numbers
	sort.Float64s(stats.numbers)
	stats.mean = sum(numbers) / float64(len(numbers))
	stats.median = median(numbers)
	stats.mode = mode(numbers)
	stats.deviation = deviation(numbers)
	return stats
}

func sum(numbers []float64) (total float64) {
	for _, x := range numbers {
		total += x
	}
	return total
}

func median(numbers []float64) float64 {
	middle := len(numbers) / 2
	result := numbers[middle]
	if len(numbers)%2 == 0 {
		result = (result + numbers[middle-1]) / 2
	}
	return result
}

func mode(numbers []float64) []float64 {
	result := make([]float64, 0)
	mostFrequent := 0
	freqs := make(map[float64]int)
	for _, x := range numbers {
		freqs[x]++
		if freqs[x] > mostFrequent {
			mostFrequent = freqs[x]
		}
	}
	for idx, x := range freqs {
		if x == mostFrequent {
			result = append(result, idx)
		}
	}
	return result
}

func deviation(numbers []float64) float64 {
	result := math.NaN()
	if len(numbers)-1 > 0 {
		numberm1 := float64(len(numbers) - 1)
		mean := sum(numbers) / float64(len(numbers))
		sumOfSquares := 0.0
		for _, x := range numbers {
			sumOfSquares += (x - mean) * (x - mean)
		}
		result = math.Sqrt(sumOfSquares / numberm1)
	}
	return result
}
