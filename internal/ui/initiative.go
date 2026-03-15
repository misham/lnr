package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/misham/linear-cli/internal/api"
)

var (
	initiativeActiveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	initiativeNameStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
)

// PrintInitiatives formats and prints a list of initiatives.
func PrintInitiatives(ios *IOStreams, initiatives []api.Initiative) error {
	if len(initiatives) == 0 {
		_, err := fmt.Fprintln(ios.Out, "No initiatives found")
		return err
	}
	if ios.IsPlain() {
		return printInitiativesPlain(ios, initiatives)
	}
	return printInitiativesStyled(ios, initiatives)
}

func initiativeOwnerName(ini *api.Initiative) string {
	if ini.Owner != nil {
		return ini.Owner.DisplayName
	}
	return "-"
}

func printInitiativesPlain(ios *IOStreams, initiatives []api.Initiative) error {
	for _, ini := range initiatives {
		target := ini.TargetDate
		if target == "" {
			target = "-"
		}
		if _, err := fmt.Fprintf(ios.Out, "%s\t%s\t%s\t%s\t%s\n",
			ini.ID,
			ini.Name,
			ini.Status,
			initiativeOwnerName(&ini),
			target,
		); err != nil {
			return err
		}
	}
	return nil
}

func printInitiativesStyled(ios *IOStreams, initiatives []api.Initiative) error {
	idW, nameW, statusW, ownerW, dateW := len("ID"), len("NAME"), len("STATUS"), len("OWNER"), len("TARGET DATE")
	type row struct {
		id, name, status, owner, targetDate string
	}
	rows := make([]row, len(initiatives))
	for i, ini := range initiatives {
		target := ini.TargetDate
		if target == "" {
			target = "-"
		}
		r := row{
			id:         ini.ID,
			name:       ini.Name,
			status:     ini.Status,
			owner:      initiativeOwnerName(&ini),
			targetDate: target,
		}
		if l := len(r.id); l > idW {
			idW = l
		}
		if l := len(r.name); l > nameW {
			nameW = l
		}
		if l := len(r.status); l > statusW {
			statusW = l
		}
		if l := len(r.owner); l > ownerW {
			ownerW = l
		}
		if l := len(r.targetDate); l > dateW {
			dateW = l
		}
		rows[i] = r
	}

	const gap = "  "
	ew := &errWriter{w: ios.Out}
	ew.printf("%s%s%s%s%s%s%s%s%s\n",
		padRight(headerStyle.Render("ID"), len("ID"), idW), gap,
		padRight(headerStyle.Render("NAME"), len("NAME"), nameW), gap,
		padRight(headerStyle.Render("STATUS"), len("STATUS"), statusW), gap,
		padRight(headerStyle.Render("OWNER"), len("OWNER"), ownerW), gap,
		headerStyle.Render("TARGET DATE"),
	)
	for i, ini := range initiatives {
		r := rows[i]
		styledStatus := r.status
		if ini.Status == "Active" {
			styledStatus = initiativeActiveStyle.Render(r.status)
		}
		ew.printf("%s%s%s%s%s%s%s%s%s\n",
			padRight(r.id, len(r.id), idW), gap,
			padRight(initiativeNameStyle.Render(r.name), len(r.name), nameW), gap,
			padRight(styledStatus, len(r.status), statusW), gap,
			padRight(r.owner, len(r.owner), ownerW), gap,
			r.targetDate,
		)
	}
	return ew.err
}

// PrintInitiativeDetail formats and prints a single initiative with full detail.
func PrintInitiativeDetail(ios *IOStreams, initiative *api.Initiative) error {
	if ios.IsPlain() {
		return printInitiativeDetailPlain(ios, initiative)
	}
	return printInitiativeDetailStyled(ios, initiative)
}

func printInitiativeDetailPlain(ios *IOStreams, ini *api.Initiative) error {
	ew := &errWriter{w: ios.Out}

	ew.printf("%s\n", ini.Name)
	ew.printf("Status:      %s\n", ini.Status)

	if ini.Owner != nil {
		ew.printf("Owner:       %s\n", ini.Owner.DisplayName)
	} else {
		ew.printf("Owner:       —\n")
	}

	if ini.Health != "" {
		ew.printf("Health:      %s\n", ini.Health)
	} else {
		ew.printf("Health:      —\n")
	}

	if ini.TargetDate != "" {
		ew.printf("Target:      %s\n", ini.TargetDate)
	}
	if ini.URL != "" {
		ew.printf("URL:         %s\n", ini.URL)
	}

	if ini.Description != "" {
		ew.printf("\n%s\n", ini.Description)
	}

	ew.printf("Created:     %s\n", ini.CreatedAt.Format("2006-01-02 15:04"))
	ew.printf("Updated:     %s\n", ini.UpdatedAt.Format("2006-01-02 15:04"))

	if ini.CompletedAt != nil {
		ew.printf("Completed:   %s\n", ini.CompletedAt.Format("2006-01-02 15:04"))
	}

	return ew.err
}

func printInitiativeDetailStyled(ios *IOStreams, ini *api.Initiative) error {
	ew := &errWriter{w: ios.Out}

	ew.printf("%s\n", initiativeNameStyle.Render(ini.Name))

	styledStatus := ini.Status
	if ini.Status == "Active" {
		styledStatus = initiativeActiveStyle.Render(ini.Status)
	}
	ew.printf("%s  %s\n", headerStyle.Render("Status:"), styledStatus)

	if ini.Owner != nil {
		ew.printf("%s  %s\n", headerStyle.Render("Owner:"), ini.Owner.DisplayName)
	} else {
		ew.printf("%s  —\n", headerStyle.Render("Owner:"))
	}

	if ini.Health != "" {
		ew.printf("%s  %s\n", headerStyle.Render("Health:"), ini.Health)
	} else {
		ew.printf("%s  —\n", headerStyle.Render("Health:"))
	}

	if ini.TargetDate != "" {
		ew.printf("%s  %s\n", headerStyle.Render("Target:"), ini.TargetDate)
	}
	if ini.URL != "" {
		ew.printf("%s  %s\n", headerStyle.Render("URL:"), ini.URL)
	}

	if ini.Description != "" {
		ew.printf("\n%s\n", ini.Description)
	}

	ew.printf("%s  %s\n", headerStyle.Render("Created:"), ini.CreatedAt.Format("2006-01-02 15:04"))
	ew.printf("%s  %s\n", headerStyle.Render("Updated:"), ini.UpdatedAt.Format("2006-01-02 15:04"))

	if ini.CompletedAt != nil {
		ew.printf("%s  %s\n", headerStyle.Render("Completed:"), ini.CompletedAt.Format("2006-01-02 15:04"))
	}

	return ew.err
}
