package git

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

func PrintResults(stats []AuthorStats, format string, out io.Writer) error {
	switch format {
	case "tabular":
		return printTabular(stats, out)
	case "csv":
		return printCSV(stats, out)
	case "json":
		return printJSON(stats, out)
	case "json-lines":
		return printJSONLines(stats, out)
	default: // cheking in validate :0
		return fmt.Errorf("unknown format: %s", format)
	}
}

func printTabular(stats []AuthorStats, out io.Writer) error {
	w := tabwriter.NewWriter(out, 0, 1, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "Name\tLines\tCommits\tFiles")
	for _, s := range stats {
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\n", s.Name, s.Lines, s.Commits, s.Files)
	}
	return w.Flush()
}

func printCSV(stats []AuthorStats, out io.Writer) error {
	w := csv.NewWriter(out)
	if err := w.Write([]string{"Name", "Lines", "Commits", "Files"}); err != nil {
		return err
	}
	for _, s := range stats {
		record := []string{
			s.Name,
			fmt.Sprintf("%d", s.Lines),
			fmt.Sprintf("%d", s.Commits),
			fmt.Sprintf("%d", s.Files),
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func printJSON(stats []AuthorStats, out io.Writer) error {
	data, err := json.Marshal(stats)
	if err != nil {
		return err
	}
	_, err = out.Write(data)
	return err
}

func printJSONLines(stats []AuthorStats, out io.Writer) error {
	enc := json.NewEncoder(out)
	for _, s := range stats {
		if err := enc.Encode(s); err != nil {
			return err
		}
	}
	return nil
}
