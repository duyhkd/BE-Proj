package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	employee := getEmployees()
	// for _, p := range employee {
	// 	fmt.Println(p)
	// }
	getSalaryByAge(employee)
}

type Employee struct {
	Id           int    `json:"id"`
	EmployeeName string `json:"employee_name"`
	Salary       int    `json:"employee_salary"`
	Age          int    `json:"employee_age"`
	ProfileImage string `json:"profile_image"`
}

type Response struct {
	Status  string
	Data    []Employee
	Message string
}

func getEmployees() []Employee {
	var response Response

	resp, err := http.Get("https://dummy.restapiexample.com/api/v1/employees")
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	body, err := io.ReadAll(resp.Body)
	errY := json.Unmarshal(body, &response)
	if errY != nil {
		fmt.Println(errY.Error())
	}
	return response.Data
}

type JobResult struct {
	Employee Employee
	Result   int
}

func worker(id int, employees <-chan Employee, results chan<- JobResult) {
	for employee := range employees {
		fmt.Println("worker", id, "started  job", employee.Id)
		res := employee.Salary / employee.Age
		fmt.Println("worker", id, "finished job", employee.Id)
		results <- JobResult{Employee: employee, Result: res}
	}
}

func getSalaryByAge(employees []Employee) {
	numJobs := len(employees)
	jobs := make(chan Employee, numJobs)
	results := make(chan JobResult, numJobs)

	// 3 workers
	for w := 0; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// Send jobs
	for j := 1; j < numJobs; j++ {
		jobs <- employees[j]
	}
	close(jobs)

	for a := 1; a < numJobs; a++ {
		r := <-results
		fmt.Println("Employee", r.Employee.EmployeeName, "salary by age is ", r.Result)
	}
	close(results)
}
