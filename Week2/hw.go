package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	p := Person{"Hoang", "Engineer", 1976}
	fmt.Println(p.workFit())
	fmt.Println(characterCountMap("Engineer")['e'])
	for _, p := range readFile("Week2/a.txt") {
		fmt.Println(p)
	}
}

type Person struct {
	Name       string
	Occupation string
	YOB        int
}

func (p Person) workFit() bool {
	return p.YOB%len(p.Name) == 0
}

func characterCountMap(str string) map[rune]int {
	m := make(map[rune]int)
	for _, char := range strings.ToLower(str) {
		m[char] = m[char] + 1
	}
	return m
}

func twoSum(nums []int, target int) []int {
	var result []int
	m := make(map[int][]int)
	for i, v := range nums {
		// Slice of index so we can handle dupes
		m[v] = append(m[v], i)
	}
	for i, v := range nums {
		term2 := target - v
		val, ok := m[term2]
		if ok && len(val) > 0 {
			// if duplicate means we have 2 val with different index
			if i == val[len(val)-1] && len(val) > 1 {
				result = append(result, i, val[len(val)-2])
				break
			} else if i != val[len(val)-1] {
				result = append(result, i, val[len(val)-1])
				break
			}

		}
	}
	return result
}

func readFile(fileName string) []Person {
	var people []Person
	f, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		person, _ := convertStringToPerson(scanner.Text())
		people = append(people, person)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return people
}

func convertStringToPerson(str string) (Person, error) {
	tokens := strings.Split(str, "|")
	yob, err := strconv.Atoi(str)
	if err != nil {
		// ... handle error
		errors.New("YOB isn't a valid integer")
	}
	if len(tokens) != 3 {
		errors.New("File is not properly format, missing delimeter")
	}
	p := Person{strings.ToUpper(tokens[0]), strings.ToLower(tokens[1]), yob}
	return p, nil
}
