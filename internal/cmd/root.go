package cmd

import (
	"fmt"
	"io"

	"github.com/1NepuNep1/gitfame/internal/git"
	"github.com/spf13/cobra"
)

type Options struct {
	Repository   string
	Revision     string
	UseCommitter bool

	OrderBy string
	Format  string

	Extensions []string
	Languages  []string
	Exclude    []string
	RestrictTo []string
}

var opts Options

func NewRootCommand(in io.Reader, out, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gitfame",
		Short: "gitfame shows git contributor statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Validate(opts); err != nil {
				return fmt.Errorf("invalid options: %w", err)
			}

			files, err := git.ListFiles(opts.Repository, opts.Revision)
			if err != nil {
				return fmt.Errorf("listing files: %w", err)
			}

			filteredFiles, err := git.FilterFiles(files, opts.Exclude, opts.RestrictTo, opts.Languages, opts.Extensions)
			if err != nil {
				fmt.Fprintf(errOut, "Warning: %v\n", err)
			}

			blamedResult, err := git.BlameFiles(filteredFiles, opts.Repository, opts.Revision, opts.UseCommitter, errOut)
			if err != nil {
				return fmt.Errorf("blaming files: %w", err)
			}

			agrResult := git.AggregateAndSort(blamedResult, opts.OrderBy)

			if err := git.PrintResults(agrResult, opts.Format, out); err != nil {
				return fmt.Errorf("printing results: %w", err)
			}

			return nil
		},
	}

	cmd.SetIn(in)
	cmd.SetOut(out)
	cmd.SetErr(errOut)

	cmd.Flags().StringVar(&opts.Repository, "repository", ".", "Path to git repo (default: current directory)")
	cmd.Flags().StringVar(&opts.Revision, "revision", "HEAD", "Pointer to commit (default: HEAD)")
	cmd.Flags().BoolVar(&opts.UseCommitter, "use-committer", false, "Use committer instead of author")
	cmd.Flags().StringVar(&opts.OrderBy, "order-by", "lines", "Sort by lines, commits, or files")
	cmd.Flags().StringVar(&opts.Format, "format", "tabular", "Output format: tabular, csv, json, json-lines")
	cmd.Flags().StringSliceVar(&opts.Extensions, "extensions", nil, "Limit files by extensions (comma-separated)")
	cmd.Flags().StringSliceVar(&opts.Languages, "languages", nil, "Limit files by languages")
	cmd.Flags().StringSliceVar(&opts.Exclude, "exclude", nil, "Glob patterns to exclude")
	cmd.Flags().StringSliceVar(&opts.RestrictTo, "restrict-to", nil, "Glob patterns to include")

	return cmd
}
