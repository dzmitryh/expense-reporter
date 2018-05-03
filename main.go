package main

import (
	"fmt"
	"os"
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"time"
)

var targetStops = [...]string{"Schiphol Airport", "Schiphol-Rijk, Boeingavenue", "Schiphol-Rijk, Beechavenue"}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("input file path argument is required")
		// do what you want
		return
	}
	inputFilePath := args[0]
	input, inputErr := os.Open(inputFilePath)
	now := time.Now()
	currMonth := now.Month()
	currYear := now.Year()
	output, outputErr := os.Create(fmt.Sprintf("transportation-compensation-bus-%[1]s-%[2]d.csv", currMonth, currYear))

	checkError("Cannot open file", inputErr)
	// automatically call Close() at the end of current method
	defer input.Close()
	checkError("Cannot create file", outputErr)
	defer output.Close()

	writer := csv.NewWriter(output)
	writer.Comma = ';'
	defer writer.Flush()

	reader := csv.NewReader(input)
	reader.Comma = ';'
	lineCount := 0
	sum := 0.0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Record ", lineCount, record, "and has", len(record), "fields")
		if lineCount == 0 {
			writeInFile(writer, record)
		}
		departure := record[2]
		destination := record[4]
		strPrice := record[5]
		if inTargetStops(&targetStops, departure, destination) && strPrice != "" {
			writeInFile(writer, record)
			if floatPrice, err := strconv.ParseFloat(strPrice, 32); err == nil {
				sum += floatPrice
			}
		}
		lineCount += 1
	}
	totalSum := strconv.FormatFloat(sum, 'f', 2, 32)
	fmt.Println("Total expenses = " + totalSum)
	writeInFile(writer, []string {"Total expenses:", "", "", "", "", totalSum, "", "", "", ""})
}

func writeInFile(writer *csv.Writer, record []string) {
	err := writer.Write(record)
	checkError("Cannot write to file", err)
}

func inTargetStops(targetStops *[3]string, departure, destination string) bool {
	return contains(targetStops, departure) || contains(targetStops, destination)
}

func contains(array *[3]string, str string) bool {
	for _, value := range array {
		if value == str {
			return true
		}
	}
	return false
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
