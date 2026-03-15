package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/misham/linear-cli/internal/api"
)

var (
	projectStatusStartedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	projectNameStyle          = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
)

// PrintProjects formats and prints a list of projects.
func PrintProjects(ios *IOStreams, projects []api.Project) error {
	if len(projects) == 0 {
		_, err := fmt.Fprintln(ios.Out, "No projects found")
		return err
	}
	if ios.IsPlain() {
		return printProjectsPlain(ios, projects)
	}
	return printProjectsStyled(ios, projects)
}

func projectLeadName(p *api.Project) string {
	if p.Lead != nil {
		return p.Lead.DisplayName
	}
	return "-"
}

func printProjectsPlain(ios *IOStreams, projects []api.Project) error {
	for _, p := range projects {
		dates := p.StartDate + " - " + p.TargetDate
		if p.StartDate == "" && p.TargetDate == "" {
			dates = "-"
		}
		if _, err := fmt.Fprintf(ios.Out, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			p.ID,
			p.SlugID,
			p.Name,
			p.Status.Type,
			projectLeadName(&p),
			fmt.Sprintf("%.0f%%", p.Progress*100),
			dates,
		); err != nil {
			return err
		}
	}
	return nil
}

func printProjectsStyled(ios *IOStreams, projects []api.Project) error {
	idW, slugW, nameW, statusW, leadW, progW, dateW := len("ID"), len("SLUG"), len("NAME"), len("STATUS"), len("LEAD"), len("PROGRESS"), len("DATES")
	type row struct {
		id, slug, name, status, lead, progress, dates string
	}
	rows := make([]row, len(projects))
	for i, p := range projects {
		dates := p.StartDate + " — " + p.TargetDate
		if p.StartDate == "" && p.TargetDate == "" {
			dates = "-"
		}
		r := row{
			id:       p.ID,
			slug:     p.SlugID,
			name:     p.Name,
			status:   p.Status.Type,
			lead:     projectLeadName(&p),
			progress: fmt.Sprintf("%.0f%%", p.Progress*100),
			dates:    dates,
		}
		if l := len(r.id); l > idW {
			idW = l
		}
		if l := len(r.slug); l > slugW {
			slugW = l
		}
		if l := len(r.name); l > nameW {
			nameW = l
		}
		if l := len(r.status); l > statusW {
			statusW = l
		}
		if l := len(r.lead); l > leadW {
			leadW = l
		}
		if l := len(r.progress); l > progW {
			progW = l
		}
		if l := len(r.dates); l > dateW {
			dateW = l
		}
		rows[i] = r
	}

	const gap = "  "
	ew := &errWriter{w: ios.Out}
	ew.printf("%s%s%s%s%s%s%s%s%s%s%s%s%s\n",
		padRight(headerStyle.Render("ID"), len("ID"), idW), gap,
		padRight(headerStyle.Render("SLUG"), len("SLUG"), slugW), gap,
		padRight(headerStyle.Render("NAME"), len("NAME"), nameW), gap,
		padRight(headerStyle.Render("STATUS"), len("STATUS"), statusW), gap,
		padRight(headerStyle.Render("LEAD"), len("LEAD"), leadW), gap,
		padRight(headerStyle.Render("PROGRESS"), len("PROGRESS"), progW), gap,
		headerStyle.Render("DATES"),
	)
	for i, p := range projects {
		r := rows[i]
		styledStatus := r.status
		if p.Status.Type == "started" {
			styledStatus = projectStatusStartedStyle.Render(r.status)
		}
		ew.printf("%s%s%s%s%s%s%s%s%s%s%s%s%s\n",
			padRight(r.id, len(r.id), idW), gap,
			padRight(r.slug, len(r.slug), slugW), gap,
			padRight(projectNameStyle.Render(r.name), len(r.name), nameW), gap,
			padRight(styledStatus, len(r.status), statusW), gap,
			padRight(r.lead, len(r.lead), leadW), gap,
			padRight(r.progress, len(r.progress), progW), gap,
			r.dates,
		)
	}
	return ew.err
}

// PrintProjectDetail formats and prints a single project with full detail.
func PrintProjectDetail(ios *IOStreams, project *api.Project) error {
	if ios.IsPlain() {
		return printProjectDetailPlain(ios, project)
	}
	return printProjectDetailStyled(ios, project)
}

func printProjectDetailPlain(ios *IOStreams, p *api.Project) error {
	ew := &errWriter{w: ios.Out}

	ew.printf("%s\n", p.Name)
	ew.printf("Status:      %s (%s)\n", p.Status.Name, p.Status.Type)

	if p.Lead != nil {
		ew.printf("Lead:        %s\n", p.Lead.DisplayName)
	} else {
		ew.printf("Lead:        —\n")
	}

	if p.PriorityLabel != "" {
		ew.printf("Priority:    %s\n", p.PriorityLabel)
	}
	ew.printf("Progress:    %.0f%%\n", p.Progress*100)

	if p.StartDate != "" || p.TargetDate != "" {
		ew.printf("Dates:       %s — %s\n", p.StartDate, p.TargetDate)
	}
	if p.URL != "" {
		ew.printf("URL:         %s\n", p.URL)
	}

	if p.Description != "" {
		ew.printf("\n%s\n", p.Description)
	}

	if len(p.Milestones) > 0 {
		ew.printf("\nMilestones:\n")
		for _, m := range p.Milestones {
			datePart := ""
			if m.TargetDate != "" {
				datePart = "  " + m.TargetDate
			}
			ew.printf("  • %s%s\n", m.Name, datePart)
			if m.Description != "" {
				ew.printf("    %s\n", m.Description)
			}
		}
	}

	ew.printf("Created:     %s\n", p.CreatedAt.Format("2006-01-02 15:04"))
	ew.printf("Updated:     %s\n", p.UpdatedAt.Format("2006-01-02 15:04"))

	if p.CompletedAt != nil {
		ew.printf("Completed:   %s\n", p.CompletedAt.Format("2006-01-02 15:04"))
	}

	return ew.err
}

func printProjectDetailStyled(ios *IOStreams, p *api.Project) error {
	ew := &errWriter{w: ios.Out}

	ew.printf("%s\n", projectNameStyle.Render(p.Name))

	statusStr := fmt.Sprintf("%s (%s)", p.Status.Name, p.Status.Type)
	if p.Status.Type == "started" {
		statusStr = projectStatusStartedStyle.Render(statusStr)
	}
	ew.printf("%s  %s\n", headerStyle.Render("Status:"), statusStr)

	if p.Lead != nil {
		ew.printf("%s  %s\n", headerStyle.Render("Lead:"), p.Lead.DisplayName)
	} else {
		ew.printf("%s  —\n", headerStyle.Render("Lead:"))
	}

	if p.PriorityLabel != "" {
		ew.printf("%s  %s\n", headerStyle.Render("Priority:"), p.PriorityLabel)
	}
	ew.printf("%s  %.0f%%\n", headerStyle.Render("Progress:"), p.Progress*100)

	if p.StartDate != "" || p.TargetDate != "" {
		ew.printf("%s  %s — %s\n", headerStyle.Render("Dates:"), p.StartDate, p.TargetDate)
	}
	if p.URL != "" {
		ew.printf("%s  %s\n", headerStyle.Render("URL:"), p.URL)
	}

	if p.Description != "" {
		ew.printf("\n%s\n", p.Description)
	}

	if len(p.Milestones) > 0 {
		ew.printf("\n%s\n", headerStyle.Render("Milestones:"))
		for _, m := range p.Milestones {
			parts := []string{"  •", m.Name}
			if m.TargetDate != "" {
				parts = append(parts, " "+m.TargetDate)
			}
			ew.printf("%s\n", strings.Join(parts, " "))
			if m.Description != "" {
				ew.printf("    %s\n", m.Description)
			}
		}
	}

	ew.printf("%s  %s\n", headerStyle.Render("Created:"), p.CreatedAt.Format("2006-01-02 15:04"))
	ew.printf("%s  %s\n", headerStyle.Render("Updated:"), p.UpdatedAt.Format("2006-01-02 15:04"))

	if p.CompletedAt != nil {
		ew.printf("%s  %s\n", headerStyle.Render("Completed:"), p.CompletedAt.Format("2006-01-02 15:04"))
	}

	return ew.err
}
