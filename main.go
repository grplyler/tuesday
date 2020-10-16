package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/logrusorgru/aurora"
)

/**------------------------------------------------------------------------
 *                             VARS / CONSTS
 *------------------------------------------------------------------------**/

const (
	FILLER_CHAR = "░"
	// PROGRESS_CHAR = "▓"
	PROGRESS_CHAR = "█"
	USE_COLOR     = true
	CSCALE_START  = uint8(214)
)

/**------------------------------------------------------------------------
 *                           Printing Functions
 *------------------------------------------------------------------------**/

func PrintHeader(week string) {
	fmt.Printf("╔══════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║  ID                        Week %s                ╔══════════╣\n", PadR(week, " ", 2))
}

func PrintClassHeader(class string) {
	fmt_class := PadR(class, " ", 9)
	fmt.Printf("╠═════╬══════════╬══════════════════════════════════╣ %s║\n", fmt_class)
}

func PrintClassHeaderFirst(class string) {
	fmt_class := PadR(class, " ", 9)
	fmt.Printf("╠═════╦══════════╦══════════════════════════════════╣ %s║\n", fmt_class)
}

func PrintTotalProgress(tasks [][]string, week string) {
	tp := CalcTotalProgress(tasks, week)
	tp_fmt := FormatTotalProgress(tp, 51, "║")
	tpi_fmt := PadL(strconv.Itoa(tp), " ", 2)
	fmt.Printf("║%s║ %s%% done ║\n", tp_fmt, tpi_fmt)
}

func PrintFooter(tasks [][]string, week string) {
	fmt.Println("╠═════╩══════════╩══════════════════════════════════╦══════════╗")
	PrintTotalProgress(tasks, week)
	fmt.Println("╚═══════════════════════════════════════════════════╩══════════╝")

}

func PrintTask(task []string, extra string) {
	content := task[2]

	// Format Content
	content_fmt := PadR(content, " ", 33)
	content_fmt += extra

	i_progress := task[3]
	id := task[5]

	id_fmt := PadL(id, " ", 4)
	// content_fmt := PadR(content, " ", 20)
	progress_fmt := FormatTaskProgress(i_progress)

	fmt.Printf("║%s ║%s║ %s\n", id_fmt, progress_fmt, content_fmt)
}

func PrintAll(lines [][]string, week string) {

	PrintHeader(week)

	last_class := ""
	lines_len := len(lines)
	for i, task := range lines {

		// Filter data by the week we want
		if task[1] == week {

			class := task[0]

			// If we are encountering a new class
			if class != last_class {

				// If this is the first task in the new class
				if last_class == "" {
					PrintClassHeaderFirst(class)
				} else {
					PrintClassHeader(class)
				}

				PrintTask(task, "╚══════════╝")

				// If we're still in the same class
			} else {

				// If next class will be different,
				// prevent out of range error
				if i == lines_len-1 {
					i -= 1
				}

				// Print tax with extra
				if lines[i+1][0] != lines[i][0] {
					PrintTask(task, "╔══════════╗")

				} else {
					PrintTask(task, "")

				}
			}

			// Update last class
			last_class = class
		}

	}

	PrintFooter(lines, week)
}

/**------------------------------------------------------------------------
 *                           FORMATTING HELPERS
 *------------------------------------------------------------------------**/

// Pad Left: Takes a string and pdds it left with specified char
func PadL(content string, fill string, width int) string {
	cwidth := len(content)
	fwidth := width - cwidth
	if fwidth < 0 {
		fwidth = -fwidth
	}
	return fmt.Sprintf("%s%s", strings.Repeat(fill, fwidth), content)
}

// Pad Right: Takes a string and pads it right with specified char
func PadR(content string, fill string, width int) string {
	cwidth := len(content)
	fwidth := width - cwidth
	if fwidth < 0 {
		fwidth = -fwidth
	}
	return fmt.Sprintf("%s%s", content, strings.Repeat(fill, fwidth))
}

func FormatTaskProgress(progress string) string {
	done, err := strconv.Atoi(progress)
	CheckError("progress is not an integer", err)
	n_blocks := done / 10.0
	blocks := strings.Repeat(PROGRESS_CHAR, n_blocks)
	filler := strings.Repeat(FILLER_CHAR, 10-n_blocks)
	fmt_progress := fmt.Sprintf("%s%s", blocks, filler)

	// Apply a color scale to the progress with aurora

	colored := ""

	if USE_COLOR {
		// 214-219
		step := 0
		color := CSCALE_START
		for _, v := range fmt_progress {
			if step == 2 {
				color += 1
				step = 0
			}
			if string(v) != FILLER_CHAR {
				colored += fmt.Sprintf("%s", aurora.Index(color, string(v)))

			} else {
				colored += string(v)
			}
			step += 1
		}

		return colored
	} else {
		return fmt_progress
	}

}

func FormatTotalProgress(tp int, width int, start string) string {
	fwidth := float64(width)
	ftp := float64(tp)

	nBlocks := math.Round((ftp / 100) * fwidth)
	fBlocks := fwidth - nBlocks

	filler := strings.Repeat(FILLER_CHAR, int(math.Round(fBlocks)))
	blocks := strings.Repeat(PROGRESS_CHAR, int(math.Round(nBlocks)))
	return fmt.Sprintf("%s%s", blocks, filler)
}

/**------------------------------------------------------------------------
 *                                  UTILS
 *------------------------------------------------------------------------**/

func CheckError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func DataPath() string {
	usr, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}
	dataPath := filepath.Join(usr.HomeDir, ".tuesday/", "main.csv")
	return dataPath
}

func IsDigit(num string) bool {
	if _, err := strconv.Atoi(num); err == nil {
		return true
	} else {
		return false
	}
}

func CalcTotalProgress(tasks [][]string, week string) int {

	total := 0.0
	done := 0.0
	for _, v := range tasks {
		inc, err := strconv.Atoi(v[3])
		if err != nil {
			fmt.Println("Error converting progress string to integer.")
			log.Fatalln(err)
		}
		if v[1] == week {
			done += float64(inc)
			total += 100.0
		}

	}

	totalProgress := float64(done) / float64(total)
	return int(math.Round(totalProgress * 100))

}

func GetLatestWeek(lines [][]string) string {
	greatest := 0
	for _, v := range lines {
		currentWeek, err := strconv.Atoi(v[1])
		CheckError("Week stored in csv data is not an integer", err)
		if currentWeek > greatest {
			greatest = currentWeek
		}
	}

	return fmt.Sprintf("%d", greatest)
}

func SortData(lines [][]string) [][]string {
	sort.SliceStable(lines, func(i, j int) bool {
		// Sort by week (latest week last)
		if lines[i][1] < lines[j][1] {
			return true
		}
		if lines[i][1] > lines[j][1] {
			return false
		}

		// Sort by class name (ASC)
		return lines[i][0] > lines[j][0]
	})

	return lines
}

/**------------------------------------------------------------------------
 *                           LOADING/SAVING
 *------------------------------------------------------------------------**/

func Boot() {
	// Creates ~/.tuesday/main.csv if it doesn't already exist
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dataFolderPath := filepath.Join(usr.HomeDir, ".tuesday/")
	_, err2 := os.Stat(dataFolderPath)
	if os.IsNotExist(err2) {
		fmt.Printf("INFO: Creating default data folder %s", dataFolderPath)
		os.Mkdir(dataFolderPath, os.ModePerm)
	}

}

func LoadCSV() [][]string {
	dpath := DataPath()

	csvFile, err := os.OpenFile(dpath, os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	r := csv.NewReader(csvFile)
	var lines [][]string

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		lines = append(lines, record)
	}

	// Sort by week then class title
	lines = SortData(lines)

	// Add ID's
	counter := 0
	for i, _ := range lines {
		lines[i] = append(lines[i], fmt.Sprintf("%d", counter))
		counter += 1
	}
	return lines

}

func SaveCSV(lines [][]string) {

	file, err := os.Create(DataPath())
	CheckError("Error opening datafile for saving", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range lines {
		if value[len(value)-1] != "assign" {
			value = value[:len(value)-1]
		}
		err := writer.Write(value)
		if err != nil {
			log.Fatalln(err)
		}
		CheckError("Cannot write to data file", err)
	}
}

/**------------------------------------------------------------------------
 *                          APPLICATION OPERATIONS
 *------------------------------------------------------------------------**/

func UpdateProgress(lines [][]string, id string, progress string) [][]string {

	iprog, err := strconv.Atoi(progress)
	CheckError("Progress must be an integer.", err)
	if iprog < 0 {
		iprog = 0
	}
	if iprog > 100 {
		iprog = 100
	}
	progress = strconv.Itoa(iprog)

	for i, v := range lines {
		if v[5] == id {
			lines[i][3] = progress
		}
	}

	SaveCSV(lines)
	return lines
}

func AddTask(lines [][]string, task []string) {
	joint := strings.Join(task, " ")
	split := strings.Split(joint, ":")
	if len(split) != 3 {
		fmt.Println("Error Assing Task: Need exactly 3 arguments")
		fmt.Println("Format: <class>:<week>:<content>")
		os.Exit(1)
	}

	lines = append(lines, []string{split[0], split[1], split[2], "0", "assign"})
	fmt.Println("Task Added:", split)
	SaveCSV(lines)
}

/**------------------------------------------------------------------------
 *                                 MAIN
 *------------------------------------------------------------------------**/

func main() {
	Boot()
	lines := LoadCSV()

	// Add Task
	if len(os.Args) >= 2 && os.Args[1] == "a" {
		AddTask(lines, os.Args[2:])
		return
	}

	// Exist if no arguments and no tasks created yet
	if lines == nil && len(os.Args) == 1 {
		fmt.Println("No tasks added yet.")
		fmt.Println("Add one with: t a <class>:<week>:<text>")
		os.Exit(1)
	}

	// Show Latest Week Summary
	if len(os.Args) == 1 {
		// Print Latest Week
		week := GetLatestWeek(lines)
		PrintAll(lines, week)
		return

		// Show specified week summary
	} else if len(os.Args) == 2 {
		// Print Specified Week
		week := os.Args[1]
		PrintAll(lines, week)
		return
	}

	// Update Progress
	if IsDigit(os.Args[1]) && (IsDigit(os.Args[2]) || os.Args[2] == "done") {
		if os.Args[2] == "done" {
			os.Args[2] = "100"
		}
		lines := UpdateProgress(lines, os.Args[1], os.Args[2])
		week := GetLatestWeek(lines)
		PrintAll(lines, week)
		return
	}

}
