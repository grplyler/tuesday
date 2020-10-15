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

func PrintTotalProgress(tasks [][]string) {
	tp := CalcTotalProgress(tasks)
	tp_fmt := FormatTotalProgress(tp, 50, "|")
	fmt.Printf("|%s| %d%% Total |\n", tp_fmt, tp)
}

func PrintFooter(tasks [][]string) {
	fmt.Println("+==============================================================+")
	PrintTotalProgress(tasks)
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

func CalcTotalProgress(tasks [][]string) int {

	total := 0.0
	done := 0.0
	for _, v := range tasks {
		inc, err := strconv.Atoi(v[3])
		if err != nil {
			fmt.Println("Error converting progress string to integer.")
			log.Fatalln(err)
		}
		done += float64(inc)
		total += 100.0

	}

	totalProgress := float64(done) / float64(total)
	return int(math.Round(totalProgress * 100))

}

func FormatTotalProgress(tp int, width int, start string) string {
	nBlocks := math.Round((float64(tp) / 100) * float64(width))
	fBlocks := float64(width) - nBlocks

	// filled = round((int(percent) / 100) * width)
	// empty = width - filled
	// return f"{start}{fill * filled}{' ' * empty}{end}"

	// Scale
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

	fmt.Printf("|%s |%s|%s\n", id_fmt, progress_fmt, content)
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

	PrintFooter(lines)
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
	sort.Slice(lines, func(i, j int) bool {
		// Sort by week (latest week last)
		if lines[i][1] < lines[j][1] {
			return true
		}
		if lines[i][1] > lines[j][1] {
			return false
		}

		// Sort by class name (ASC)
		return lines[i][0] < lines[j][0]
	})

	return lines
}

func load_csv() [][]string {
	dpath := data_path()
	fmt.Println(dpath)

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
	sort_data(lines)

	// Add ID's
	counter := 0
	for i, _ := range lines {
		lines[i] = append(lines[i], fmt.Sprintf("%d", counter))
		counter += 1
	}
	return lines

}

func data_path() string {
	usr, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}
	dataPath := filepath.Join(usr.HomeDir, ".tuesday/", "main.csv")
	return dataPath
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
	usr, _ := user.Current()
	fmt.Println(usr.HomeDir)
	Boot()
	lines := load_csv()

	fmt.Println(len(os.Args))

	if len(os.Args) == 1 {
		// Print Latest Week
		week := GetLatestWeek(lines)
		PrintAll(lines, week)
	} else {
		// Print Specified Week
		week := os.Args[1]
		fmt.Println("week", week)
		PrintAll(lines, week)
	}

	fmt.Println("total Progress", CalcTotalProgress(lines))

}
