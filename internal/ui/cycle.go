package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/misham/linear-cli/internal/api"
)

var (
	cycleActiveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	cycleNumberStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
)

func cycleStatus(c *api.Cycle) string {
	switch {
	case c.IsActive:
		return "Active"
	case c.IsNext:
		return "Next"
	case c.IsPast:
		return "Past"
	default:
		return "Future"
	}
}

// PrintCycles formats and prints a list of cycles.
func PrintCycles(ios *IOStreams, cycles []api.Cycle) error {
	if len(cycles) == 0 {
		_, err := fmt.Fprintln(ios.Out, "No cycles found")
		return err
	}
	if ios.IsPlain() {
		return printCyclesPlain(ios, cycles)
	}
	return printCyclesStyled(ios, cycles)
}

func printCyclesPlain(ios *IOStreams, cycles []api.Cycle) error {
	for _, c := range cycles {
		if _, err := fmt.Fprintf(ios.Out, "%d\t%s\t%s\t%s - %s\t%s\n",
			c.Number,
			c.Name,
			cycleStatus(&c),
			c.StartsAt.Format("Jan 02"),
			c.EndsAt.Format("Jan 02"),
			fmt.Sprintf("%.0f%%", c.Progress*100),
		); err != nil {
			return err
		}
	}
	return nil
}

func printCyclesStyled(ios *IOStreams, cycles []api.Cycle) error {
	numW, nameW, statusW, dateW := len("NUMBER"), len("NAME"), len("STATUS"), len("DATES")
	for _, c := range cycles {
		if l := len(fmt.Sprintf("%d", c.Number)); l > numW {
			numW = l
		}
		if l := len(c.Name); l > nameW {
			nameW = l
		}
		if l := len(cycleStatus(&c)); l > statusW {
			statusW = l
		}
		dateStr := c.StartsAt.Format("Jan 02") + " - " + c.EndsAt.Format("Jan 02")
		if l := len(dateStr); l > dateW {
			dateW = l
		}
	}

	const gap = "  "
	ew := &errWriter{w: ios.Out}
	ew.printf("%s%s%s%s%s%s%s%s%s\n",
		padRight(headerStyle.Render("NUMBER"), len("NUMBER"), numW), gap,
		padRight(headerStyle.Render("NAME"), len("NAME"), nameW), gap,
		padRight(headerStyle.Render("STATUS"), len("STATUS"), statusW), gap,
		padRight(headerStyle.Render("DATES"), len("DATES"), dateW), gap,
		headerStyle.Render("PROGRESS"),
	)
	for _, c := range cycles {
		numStr := fmt.Sprintf("%d", c.Number)
		status := cycleStatus(&c)
		dateStr := c.StartsAt.Format("Jan 02") + " - " + c.EndsAt.Format("Jan 02")
		progressStr := fmt.Sprintf("%.0f%%", c.Progress*100)

		styledStatus := status
		if c.IsActive {
			styledStatus = cycleActiveStyle.Render(status)
		}

		ew.printf("%s%s%s%s%s%s%s%s%s\n",
			padRight(cycleNumberStyle.Render(numStr), len(numStr), numW), gap,
			padRight(c.Name, len(c.Name), nameW), gap,
			padRight(styledStatus, len(status), statusW), gap,
			padRight(dateStr, len(dateStr), dateW), gap,
			progressStr,
		)
	}
	return ew.err
}

// PrintCycleDetail formats and prints a single cycle with full detail.
func PrintCycleDetail(ios *IOStreams, cycle *api.Cycle) error {
	if ios.IsPlain() {
		return printCycleDetailPlain(ios, cycle)
	}
	return printCycleDetailStyled(ios, cycle)
}

func printCycleDetailPlain(ios *IOStreams, c *api.Cycle) error {
	ew := &errWriter{w: ios.Out}

	ew.printf("#%d  %s\n", c.Number, c.Name)
	ew.printf("Status:    %s\n", cycleStatus(c))
	ew.printf("Team:      %s\n", c.Team.Name)
	ew.printf("Dates:     %s - %s\n", c.StartsAt.Format("Jan 02"), c.EndsAt.Format("Jan 02"))
	ew.printf("Progress:  %.0f%%\n", c.Progress*100)

	if c.Description != "" {
		ew.printf("\n%s\n", c.Description)
	}

	ew.printf("Created:   %s\n", c.CreatedAt.Format("2006-01-02 15:04"))
	ew.printf("Updated:   %s\n", c.UpdatedAt.Format("2006-01-02 15:04"))

	if c.CompletedAt != nil {
		ew.printf("Completed: %s\n", c.CompletedAt.Format("2006-01-02 15:04"))
	}

	return ew.err
}

func printCycleDetailStyled(ios *IOStreams, c *api.Cycle) error {
	ew := &errWriter{w: ios.Out}

	ew.printf("%s  %s\n", cycleNumberStyle.Render(fmt.Sprintf("#%d", c.Number)), lipgloss.NewStyle().Bold(true).Render(c.Name))

	status := cycleStatus(c)
	styledStatus := status
	if c.IsActive {
		styledStatus = cycleActiveStyle.Render(status)
	}
	ew.printf("%s  %s\n", headerStyle.Render("Status:"), styledStatus)
	ew.printf("%s  %s\n", headerStyle.Render("Team:"), c.Team.Name)
	ew.printf("%s  %s - %s\n", headerStyle.Render("Dates:"), c.StartsAt.Format("Jan 02"), c.EndsAt.Format("Jan 02"))
	ew.printf("%s  %.0f%%\n", headerStyle.Render("Progress:"), c.Progress*100)

	if c.Description != "" {
		ew.printf("\n%s\n", c.Description)
	}

	ew.printf("%s  %s\n", headerStyle.Render("Created:"), c.CreatedAt.Format("2006-01-02 15:04"))
	ew.printf("%s  %s\n", headerStyle.Render("Updated:"), c.UpdatedAt.Format("2006-01-02 15:04"))

	if c.CompletedAt != nil {
		ew.printf("%s  %s\n", headerStyle.Render("Completed:"), c.CompletedAt.Format("2006-01-02 15:04"))
	}

	return ew.err
}
