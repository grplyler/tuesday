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
)

func FormatTaskProgress(progress string) string {
	done, _ := strconv.Atoi(progress)
	n_blocks := done / 10.0
	blocks := strings.Repeat("█", n_blocks)
	filler := strings.Repeat(" ", 10-n_blocks)
	return fmt.Sprintf("%s%s", blocks, filler)
}

func PrintHeader(week string) {
	fmt.Printf("+--------------------------------------------------------------+\n")
	fmt.Printf("|                            Week %s                           \n", week)
}

func PrintClassHeader(class string) {
	fmt.Printf("+-----+----------+--------------------------------- [ %s ]\n", class)
}

func PrintTotalProgress(tasks [][]string, week string) {
	tp := CalcTotalProgress(tasks, week)
	tp_fmt := FormatTotalProgress(tp, 50, "|")
	fmt.Printf("|%s| %d%% Total |\n", tp_fmt, tp)
}

func PrintFooter(tasks [][]string, week string) {
	fmt.Println("+==============================================================+")
	PrintTotalProgress(tasks, week)
	fmt.Println("+==============================================================+")

}

func Filll(content string, fill string, width int) string {
	cwidth := len(content)
	fwidth := width - cwidth
	if fwidth < 0 {
		fwidth = -fwidth
	}
	return fmt.Sprintf("%s%s", strings.Repeat(fill, fwidth), content)
}

func Fillr(content string, fill string, width int) string {
	cwidth := len(content)
	fwidth := width - cwidth
	if fwidth < 0 {
		fwidth = -fwidth
	}
	return fmt.Sprintf("%s%s", content, strings.Repeat(fill, fwidth))
}

func FormatTaskID() {

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

func FormatTotalProgress(tp int, width int, start string) string {
	nBlocks := math.Round((float64(tp) / 100) * float64(width))
	fBlocks := float64(width) - nBlocks

	filler := strings.Repeat(" ", int(math.Round(fBlocks)))
	blocks := strings.Repeat("█", int(math.Round(nBlocks)))
	return fmt.Sprintf("%s%s", blocks, filler)
}

func PrintTask(task []string) {
	// class := task[0]
	content := task[2]
	i_progress := task[3]
	id := task[5]

	id_fmt := Filll(id, " ", 4)
	// content_fmt := Fillr(content, " ", 20)
	progress_fmt := FormatTaskProgress(i_progress)

	fmt.Printf("|%s |%s| %s\n", id_fmt, progress_fmt, content)
}

func PrintAll(lines [][]string, week string) {

	PrintHeader(week)

	last_class := ""
	for _, task := range lines {

		// If its the week we want
		if task[1] == week {

			class := task[0]
			// Start a new class header
			if class != last_class {
				PrintClassHeader(class)
				PrintTask(task)
				// Continue old class header
			} else {
				PrintTask(task)
			}

			last_class = class
		}

	}

	PrintFooter(lines, week)
}

func GetLatestWeek(lines [][]string) string {
	greatest := 0
	for _, v := range lines {
		currentWeek, _ := strconv.Atoi(v[1])
		if currentWeek > greatest {
			greatest = currentWeek
		}
	}

	return fmt.Sprintf("%d", greatest)
}

func sort_data(lines [][]string) [][]string {
	sort.SliceStable(lines, func(i, j int) bool {
		// Sort by week (latest week last)
		// if lines[i][1] < lines[j][1] {
		// 	return true
		// }
		// if lines[i][1] > lines[j][1] {
		// 	return false
		// }

		// Sort by class name (ASC)
		return lines[i][0] > lines[j][0]
	})

	return lines
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
	lines = sort_data(lines)

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
	fmt.Println("Added Task:", split)
	SaveCSV(lines)
}

func Boot() {
	// Creates ~/.tuesday/main.csv if it doesn't already exist
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dataFolderPath := filepath.Join(usr.HomeDir, ".tuesday/")
	folderInfo, err := os.Stat(dataFolderPath)
	if folderInfo.IsDir() && os.IsNotExist(err) {
		os.Mkdir(dataFolderPath, os.ModePerm)
	}

}

func main() {
	Boot()
	lines := LoadCSV()

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

	// Add Todo
	if os.Args[1] == "a" {
		AddTask(lines, os.Args[2:])
	}

}
