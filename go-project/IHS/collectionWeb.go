package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// ---------- Same grading logic as CLI ----------
type Course struct {
	Name       string
	Unit       int
	Coursework float64
	Exam       float64
	Total      float64
	Grade      string
	PointValue float64
	GradePoint float64
}

type Semester struct {
	Name       string
	Courses    []Course
	TotalUnits int
	TotalGP    float64
	GPA        float64
}

func calculateGrade(score float64) (string, float64) {
	switch {
	case score >= 70.0:
		return "A", 5.0
	case score >= 60.0:
		return "B", 4.0
	case score >= 50.0:
		return "C", 3.0
	case score >= 45.0:
		return "D", 2.0
	case score >= 40.0:
		return "E", 1.0
	default:
		return "F", 0.0
	}
}

func processSemester(name string, courses []Course) Semester {
	var totalUnits int
	var totalGP float64
	for i := range courses {
		c := &courses[i]
		c.Total = c.Coursework + c.Exam
		c.Grade, c.PointValue = calculateGrade(c.Total)
		c.GradePoint = float64(c.Unit) * c.PointValue
		totalUnits += c.Unit
		totalGP += c.GradePoint
	}
	var gpa float64
	if totalUnits > 0 {
		gpa = totalGP / float64(totalUnits)
	}
	return Semester{
		Name:       name,
		Courses:    courses,
		TotalUnits: totalUnits,
		TotalGP:    totalGP,
		GPA:        gpa,
	}
}

// ---------- HTTP Handlers ----------
var tmpl = template.Must(template.ParseFiles("templates/index_cpa.html"))

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// Render the form with no result yet
	tmpl.Execute(w, nil)
}

func calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// We expect multiple courses: names[], units[], cw[], exam[]
	names := r.Form["name[]"]
	unitsStr := r.Form["unit[]"]
	cwStr := r.Form["cw[]"]
	examStr := r.Form["exam[]"]

	if len(names) == 0 {
		http.Error(w, "No courses submitted", http.StatusBadRequest)
		return
	}

	courses := make([]Course, 0, len(names))
	for i := range names {
		name := strings.TrimSpace(names[i])
		if name == "" {
			continue
		}
		unit, _ := strconv.Atoi(unitsStr[i])
		cw, _ := strconv.ParseFloat(cwStr[i], 64)
		exam, _ := strconv.ParseFloat(examStr[i], 64)
		// Clamp values just in case
		if cw < 0 {
			cw = 0
		}
		if cw > 40 {
			cw = 40
		}
		if exam < 0 {
			exam = 0
		}
		if exam > 60 {
			exam = 60
		}
		courses = append(courses, Course{
			Name:       name,
			Unit:       unit,
			Coursework: cw,
			Exam:       exam,
		})
	}

	if len(courses) == 0 {
		http.Error(w, "No valid courses", http.StatusBadRequest)
		return
	}

	// Get semester name
	semName := strings.TrimSpace(r.FormValue("semester_name"))
	if semName == "" {
		semName = "Semester"
	}

	// Process the semester
	semester := processSemester(semName, courses)

	// Render result
	tmpl.Execute(w, semester)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/calculate", calculate)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
