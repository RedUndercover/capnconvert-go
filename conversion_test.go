package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func hasSameLines(file1, file2 string) (bool, error) {
	readAndSortLines := func(filename string) ([]string, error) {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		sort.Strings(lines)
		return lines, nil
	}

	lines1, err := readAndSortLines(file1)
	if err != nil {
		return false, err
	}

	lines2, err := readAndSortLines(file2)
	if err != nil {
		return false, err
	}

	if len(lines1) != len(lines2) {
		return false, nil
	}

	for i := range lines1 {
		if lines1[i] != lines2[i] {
			return false, nil
		}
	}

	return true, nil
}

func TestConversion(t *testing.T) {
	// Define paths to the files
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current directory: %v", err)
	}

	inputFilePath := filepath.Join(cwd, "./testdata/test_file.go")
	expectedOutputPath := filepath.Join(cwd, "./testdata/expected_output.capnp")
	outputDirPath := filepath.Join(cwd, "./testdata/temp")
	outputFilePath := filepath.Join(outputDirPath, "test_output.capnp")

	// Ensure the output directory exists
	if err := os.MkdirAll(outputDirPath, 0755); err != nil {
		t.Fatalf("Error creating output directory: %v", err)
	}

	// Step 1: Run the conversion function
	schema, err := Convert(inputFilePath)
	if err != nil {
		t.Fatalf("Error running conversion function: %v", err)
	}
	schemaStr := schema.String()

	// Step 2: Read the expected output
	expected, err := ioutil.ReadFile(expectedOutputPath)
	if err != nil {
		t.Fatalf("Error reading expected output: %v", err)
	}

	// Step 3: Write the output to the file
	err = ioutil.WriteFile(outputFilePath, []byte(schemaStr), 0644)
	if err != nil {
		t.Fatalf("Error writing to file: %v", err)
	}

	// Step 4: Compare the files
	result, err := hasSameLines(outputFilePath, expectedOutputPath)
	if err != nil {
		t.Fatalf("Error comparing files: %v", err)
	}

	if !result {
		t.Errorf("Test failed!\nExpected:\n%s\nGot:\n%s", string(expected), schemaStr)
	}
}
