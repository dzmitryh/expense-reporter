package main

import (
	"fmt"
	"os"
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

var targetStopsStr = getEnv("REPORTING_BUS_STOPS",
	"Schiphol Airport | Schiphol-Rijk, Boeingavenue | Schiphol-Rijk, Beechavenue | Schiphol, Schipholgebouw | Schiphol, P12/Vrachtgebouw")
var targetStops = strings.Split(targetStopsStr, "|")

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		log.Fatal("Arguments are requried: 1) input file path 2) output dir path")
		return
	}
	inputFilePath := args[0]
	input, inputErr := os.Open(inputFilePath)
	checkError("Cannot open file", inputErr)
	defer input.Close()

	outputDirPath := args[1]
	if _, err := os.Stat(outputDirPath); os.IsNotExist(err) {
		log.Fatalf("output dir %[1]s doesn't exist", outputDirPath)
		return
	}
	now := time.Now()
	currMonth := now.Month()
	currYear := now.Year()
	output, outputErr := os.Create(fmt.Sprintf(
		"%[1]s/transportation-compensation-bus-%[2]s-%[3]d.csv", outputDirPath, currMonth, currYear))
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
		if inTargetStops(targetStops, departure, destination) && strPrice != "" {
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

func inTargetStops(targetStops []string, departure, destination string) bool {
	return contains(targetStops, departure) || contains(targetStops, destination)
}

func contains(array []string, str string) bool {
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

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
