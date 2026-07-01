package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Course struct {
	Name       string
	Unit       int
	Coursework float64
	Exam       float64
	Total      float64
	Grade      string // A, B, C, D, E, F
	PointValue float64
	GradePoint float64 // Unit * PointValue
}

// Semester represents a single term/session containing multiple courses
type Semester struct {
	Name       string
	Courses    []Course
	TotalUnits int
	TotalGP    float64
	GPA        float64
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	printHeader("WELCOME TO THE STUDENT RESULT PROCESSOR & TRANSCRIPT GENERATOR")

	fmt.Println("This program processes academic results with the following grading scale (5.0 GP Max):")
	fmt.Println(" - Coursework is graded over 40")
	fmt.Println(" - Exam is graded over 60")
	fmt.Println(" - Grade scale: 70+ [A = 5], 60-69 [B = 4], 50-59 [C = 3], 45-49 [D = 2], 40-44 [E = 1], <40 [F = 0]")
	fmt.Println()

	for {
		fmt.Println("Please choose a calculation mode:")
		fmt.Println("1. GPA (Single Session / Semester)")
		fmt.Println("2. CPA (Cumulative Point Average - Multiple Sessions / Semesters)")
		fmt.Println("3. Exit Program")
		fmt.Print("Enter your choice (1-3): ")

		choiceStr := readLine(scanner)
		choiceStr = strings.TrimSpace(choiceStr)

		if choiceStr == "3" {
			fmt.Println("\nThank you for using the Result Processor. Goodbye!")
			break
		}

		if choiceStr == "1" {
			processSingleSemester(scanner)
		} else if choiceStr == "2" {
			processMultipleSemesters(scanner)
		} else {
			fmt.Println("\n[!] Invalid choice. Please enter 1, 2, or 3.")
			printLine()
		}
	}
}

// processSingleSemester handles inputting and processing for a single semester GPA
func processSingleSemester(scanner *bufio.Scanner) {
	printHeader("GPA GENERATION - SINGLE SEMESTER")

	semName := readRequiredString(scanner, "Enter Semester/Session Name (e.g., Year 1 Semester 1): ")
	semester := inputSemesterData(scanner, semName)

	if len(semester.Courses) == 0 {
		fmt.Println("\n[!] No courses entered. Returning to main menu.")
		return
	}

	displayTranscript([]Semester{semester}, semester.GPA)
}

// processMultipleSemesters handles inputting and processing for multiple semesters CPA
func processMultipleSemesters(scanner *bufio.Scanner) {
	printHeader("CPA GENERATION - MULTIPLE SEMESTERS")

	var semesters []Semester
	semCount := 1

	for {
		fmt.Printf("\n--- Adding Semester #%d ---\n", semCount)
		semDefaultName := fmt.Sprintf("Semester %d", semCount)
		semName := readRequiredString(scanner, fmt.Sprintf("Enter Name for Semester %d [or Enter for '%s']: ", semCount, semDefaultName))
		if semName == "" {
			semName = semDefaultName
		}

		semester := inputSemesterData(scanner, semName)
		if len(semester.Courses) > 0 {
			semesters = append(semesters, semester)
		} else {
			fmt.Println("[!] No courses added for this semester.")
		}

		fmt.Print("\nDo you want to add another semester? (y/n): ")
		addAnother := strings.ToLower(readLine(scanner))
		if addAnother != "y" && addAnother != "yes" {
			break
		}
		semCount++
	}

	if len(semesters) == 0 {
		fmt.Println("\n[!] No semester records completed. Returning to main menu.")
		return
	}

	// Calculate overall CPA
	var grandTotalUnits int
	var grandTotalGP float64

	for _, sem := range semesters {
		grandTotalUnits += sem.TotalUnits
		grandTotalGP += sem.TotalGP
	}

	var cpa float64 = 0.0
	if grandTotalUnits > 0 {
		cpa = grandTotalGP / float64(grandTotalUnits)
	}

	displayTranscript(semesters, cpa)
}

// inputSemesterData prompts the user to enter courses and scores for a specific semester
func inputSemesterData(scanner *bufio.Scanner, semesterName string) Semester {
	var courses []Course
	var totalUnits int
	var totalGP float64

	fmt.Printf("\n>>> Inputting data for: %s <<<\n", semesterName)

	for {
		fmt.Println("\n--- New Course Entry ---")
		courseName := readRequiredString(scanner, "Enter Course Code/Name (e.g., MTH101): ")

		unit := readIntRange(scanner, "Enter Course Units (1-6): ", 1, 6)
		cw := readFloatRange(scanner, "Enter Coursework Score (0.0 to 40.0): ", 0.0, 40.0)
		exam := readFloatRange(scanner, "Enter Exam Score (0.0 to 60.0): ", 0.0, 60.0)

		total := cw + exam
		grade, pointVal := calculateGrade(total)
		gp := float64(unit) * pointVal

		course := Course{
			Name:       courseName,
			Unit:       unit,
			Coursework: cw,
			Exam:       exam,
			Total:      total,
			Grade:      grade,
			PointValue: pointVal,
			GradePoint: gp,
		}

		courses = append(courses, course)
		totalUnits += unit
		totalGP += gp

		fmt.Printf("\nAdded: %s | Total Score: %.2f | Grade: %s | Grade Points: %.2f\n", courseName, total, grade, gp)

		fmt.Print("Add another course to this semester? (y/n): ")
		again := strings.ToLower(readLine(scanner))
		if again != "y" && again != "yes" {
			break
		}
	}

	var gpa float64 = 0.0
	if totalUnits > 0 {
		gpa = totalGP / float64(totalUnits)
	}

	return Semester{
		Name:       semesterName,
		Courses:    courses,
		TotalUnits: totalUnits,
		TotalGP:    totalGP,
		GPA:        gpa,
	}
}

// calculateGrade assigns the letter grade and numeric points (5-point scale)
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

// displayTranscript prints the academic summary
func displayTranscript(semesters []Semester, finalAverage float64) {
	printHeader("OFFICIAL ACADEMIC TRANSCRIPT REPORT")

	isCPA := len(semesters) > 1

	for _, sem := range semesters {
		fmt.Printf("\nSEMESTER: %s\n", strings.ToUpper(sem.Name))
		fmt.Println("+------------------+--------+------------+------------+-------------+-------+-------------+")
		fmt.Printf("| %-16s | %-6s | %-10s | %-10s | %-11s | %-5s | %-11s |\n",
			"Course Code", "Units", "Coursework", "Exam Score", "Total Score", "Grade", "Grade Point")
		fmt.Println("+------------------+--------+------------+------------+-------------+-------+-------------+")

		for _, course := range sem.Courses {
			fmt.Printf("| %-16s | %-6d | %-10.2f | %-10.2f | %-11.2f | %-5s | %-11.2f |\n",
				course.Name, course.Unit, course.Coursework, course.Exam, course.Total, course.Grade, course.GradePoint)
		}

		fmt.Println("+------------------+--------+------------+------------+-------------+-------+-------------+")
		fmt.Printf("  Semester Summary: Total Units: %d | Total Grade Points Earned: %.2f | Semester GPA: %.2f\n",
			sem.TotalUnits, sem.TotalGP, sem.GPA)
		printLine()
	}

	if isCPA {
		fmt.Println("\n=========================================================================")
		fmt.Printf("   FINAL SUMMARY STATUS: CUMULATIVE POINT AVERAGE (CPA)\n")
		var totalSemestersUnits int
		var totalSemestersGP float64
		for _, sem := range semesters {
			totalSemestersUnits += sem.TotalUnits
			totalSemestersGP += sem.TotalGP
		}
		fmt.Printf("   Total Credit Units Undertaken: %d\n", totalSemestersUnits)
		fmt.Printf("   Total Grade Points Accumulated: %.2f\n", totalSemestersGP)
		fmt.Printf("   OVERALL CPA (CGPA): %.2f / 5.00\n", finalAverage)
		fmt.Println("=========================================================================")
	} else {
		fmt.Println("\n=========================================================================")
		fmt.Printf("   FINAL SUMMARY STATUS: SINGLE SEMESTER SUMMARY GPA\n")
		fmt.Printf("   OVERALL GPA: %.2f / 5.00\n", finalAverage)
		fmt.Println("=========================================================================")
	}

	fmt.Println("\nPress Enter to return to the main menu...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// Helpers for the Input

func readLine(scanner *bufio.Scanner) string {
	scanner.Scan()
	return scanner.Text()
}

func readRequiredString(scanner *bufio.Scanner, prompt string) string {
	for {
		fmt.Print(prompt)
		val := strings.TrimSpace(readLine(scanner))
		if val != "" {
			return val
		}
		// but we allow empty User input (like custom semester name default)
		if strings.Contains(prompt, "[or Enter for") {
			return ""
		}
		fmt.Println("[!] Input cannot be blank. Please try again.")
	}
}

func readFloatRange(scanner *bufio.Scanner, prompt string, min, max float64) float64 {
	for {
		fmt.Print(prompt)
		valStr := strings.TrimSpace(readLine(scanner))
		val, err := strconv.ParseFloat(valStr, 64)
		if err == nil && val >= min && val <= max {
			return val
		}
		fmt.Printf("[!] Invalid input. Please enter a decimal value between %.1f and %.1f.\n", min, max)
	}
}

func readIntRange(scanner *bufio.Scanner, prompt string, min, max int) int {
	for {
		fmt.Print(prompt)
		valStr := strings.TrimSpace(readLine(scanner))
		val, err := strconv.Atoi(valStr)
		if err == nil && val >= min && val <= max {
			return val
		}
		fmt.Printf("[!] Invalid input. Please enter an integer between %d and %d.\n", min, max)
	}
}

// Helpers for the Visual Render

func printHeader(title string) {
	width := len(title) + 6
	border := strings.Repeat("=", width)
	fmt.Println("\n" + border)
	fmt.Printf("   %s\n", title)
	fmt.Println(border)
}

func printLine() {
	fmt.Println(strings.Repeat("-", 80))
}
