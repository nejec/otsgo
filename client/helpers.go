package client

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

func fmtTableOutput(header []string, data [][]string) {
	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Borders:  tw.BorderNone,
			Settings: tw.Settings{Separators: tw.SeparatorsNone, Lines: tw.LinesNone},
		})),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{AutoWrap: tw.WrapNone, AutoFormat: tw.On},
				Alignment:  tw.CellAlignment{Global: tw.AlignLeft},
				Padding:    tw.CellPadding{Global: tw.PaddingNone},
			},
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{AutoWrap: tw.WrapNone},
				Alignment:  tw.CellAlignment{Global: tw.AlignLeft},
				Padding:    tw.CellPadding{Global: tw.Padding{Right: "\t"}},
			},
		}),
	)
	table.Header(header)
	if err := table.Bulk(data); err != nil {
		fmt.Fprintln(os.Stderr, "table bulk:", err)
	}
	if err := table.Render(); err != nil {
		fmt.Fprintln(os.Stderr, "table render:", err)
	}
}

func loadHistory(limit int) (History, error) {
	var history History

	home, err := os.UserHomeDir()

	if err != nil {
		fmt.Printf("Error homedir history file: %s\n", err)
	}
	file, err := os.Open(filepath.Join(home, ".ots_history"))

	if err != nil {
		return nil, fmt.Errorf("opening history file: %v", err)
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line into timestamp and key
		parts := strings.Split(line, ";")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line format: %q", line)
		}

		entry := strings.TrimSpace(parts[1])
		if len(entry) == 0 {
			continue
		}
		history = append(history, entry)
		// to do check expired entries to not load
		// timeNow := time.Now()
		// strTimestamp, err := strconv.ParseInt(parts[0], 10, 64)
		// if err != nil {
		// 	fmt.Printf("Error converting string to int: %v\n", err)
		// }
		// timestamp := time.Unix(strTimestamp, 0)

		// expired := timeNow.Before(timestamp)

		// if !expired {
		// 	history = append(history, entry)
		// }

	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading history file: %v", err)
	}

	if len(history) >= limit {
		history = history[len(history)-limit:]
	}

	return history, nil

}

func writeHistory(entry string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting homedir: %s\n", err)
	}

	file, err := os.OpenFile(filepath.Join(home, ".ots_history"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening history file: %s", err)
	}
	defer file.Close()

	// Format the entry as a string
	entryString := fmt.Sprintf("%d;%s\n", time.Now().Unix(), entry)

	if _, err := file.WriteString(entryString); err != nil {
		return fmt.Errorf("error writing entry to file: %s", err)
	}

	return nil
}
