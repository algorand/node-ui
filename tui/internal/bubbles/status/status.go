package status

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"

	"github.com/algorand/node-ui/messages"
	"github.com/algorand/node-ui/tui/internal/style"
)

type Model struct {
	Status  models.NodeStatus
	Network messages.NetworkMsg
	Err     error

	style             *style.Styles
	requestor         *messages.Requestor
	progress          progress.Model
	processedAcctsPct float64
	verifiedAcctsPct  float64
	acquiredBlksPct   float64
}

func New(style *style.Styles, requestor *messages.Requestor) Model {
	return Model{
		style:     style,
		progress:  progress.New(progress.WithDefaultGradient()),
		requestor: requestor,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.requestor.GetNetworkCmd(),
		m.requestor.GetStatusCmd(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.StatusMsg:
		if msg.Error != nil {
			m.Err = fmt.Errorf("error fetching status: %w", msg.Error)
			return m, tea.Quit
		}
		m.Status = msg.Status
		if m.Status.CatchpointTotalAccounts > 0 {
			m.processedAcctsPct = float64(m.Status.CatchpointProcessedAccounts) / float64(m.Status.CatchpointTotalAccounts)
			m.verifiedAcctsPct = float64(m.Status.CatchpointVerifiedAccounts) / float64(m.Status.CatchpointTotalAccounts)
		}
		if m.Status.CatchpointTotalBlocks > 0 {
			m.processedAcctsPct = 1
			m.verifiedAcctsPct = 1
			m.acquiredBlksPct = float64(m.Status.CatchpointAcquiredBlocks) / float64(m.Status.CatchpointTotalBlocks)
		}
		return m, tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
			return m.requestor.GetStatusCmd()()
		})

	case messages.NetworkMsg:
		m.Network = msg

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	return m, nil
}

func formatVersion(v string) string {
	i := strings.LastIndex(v, "/")
	if i != 0 {
		i++
	}
	return v[i:]
}

func formatNextVersion(last, next string, round uint64) string {
	if last == next {
		return "N/A"
	}
	return strconv.FormatUint(round, 10)
}

func writeProgress(b *strings.Builder, prefix string, progress progress.Model, pct float64) {
	b.WriteString(prefix)
	b.WriteString(progress.ViewAs(pct))
	b.WriteString("\n")
}

func (m Model) View() string {
	bold := m.style.StatusBoldText
	key := m.style.BottomListItemKey.Copy().MarginLeft(0)
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("%s %s\n", bold.Render("Network:"), m.Network.GenesisID))
	builder.WriteString(fmt.Sprintf("%s %s\n", bold.Render("Genesis:"), base64.StdEncoding.EncodeToString(m.Network.GenesisHash[:])))
	// TODO: get rid of magic number
	height := style.TopHeight - 2 - 3 // 3 is the padding/margin/border
	// status
	if (m.Status != models.NodeStatus{}) {
		nextVersion := formatNextVersion(
			m.Status.LastVersion,
			m.Status.NextVersion,
			m.Status.NextVersionRound)

		switch {
		case m.Status.Catchpoint != "":
			// Catchpoint view
			builder.WriteString(fmt.Sprintf("\n    Catchpoint: %s\n", key.Render(strings.Split(m.Status.Catchpoint, "#")[0])))
			var catchupStatus string
			switch {
			case m.Status.CatchpointAcquiredBlocks > 0:
				catchupStatus = fmt.Sprintf("    Downloading blocks:   %5d / %d\n", m.Status.CatchpointAcquiredBlocks, m.Status.CatchpointTotalBlocks)
			case m.Status.CatchpointVerifiedAccounts > 0:
				catchupStatus = fmt.Sprintf("    Processing accounts:   %d / %d\n", m.Status.CatchpointVerifiedAccounts, m.Status.CatchpointTotalAccounts)
			case m.Status.CatchpointProcessedAccounts > 0:
				catchupStatus = fmt.Sprintf("    Downloading accounts: %d / %d\n", m.Status.CatchpointProcessedAccounts, m.Status.CatchpointTotalAccounts)
			default:
				catchupStatus = "\n"
			}
			builder.WriteString(bold.Render(catchupStatus))
			builder.WriteString("\n")
			writeProgress(&builder, "Downloading accounts: ", m.progress, m.processedAcctsPct)
			writeProgress(&builder, "Processing accounts:  ", m.progress, m.verifiedAcctsPct)
			writeProgress(&builder, "Downloading blocks:   ", m.progress, m.acquiredBlksPct)
			height -= 7
		default:
			builder.WriteString(fmt.Sprintf("Current round:   %s\n", key.Render(strconv.FormatUint(m.Status.LastRound, 10))))
			builder.WriteString(fmt.Sprintf("Block wait time: %s\n", time.Nanosecond*time.Duration(m.Status.TimeSinceLastRound)))
			builder.WriteString(fmt.Sprintf("Sync time:       %s\n", time.Second*time.Duration(m.Status.CatchupTime)))
			height -= 3
			// TODO: Display consensus upgrade progress
			if m.Status.LastVersion == m.Status.NextVersion {
				// no upgrade in progress
				builder.WriteString(fmt.Sprintf("Protocol:        %s\n", formatVersion(string(m.Status.LastVersion))))
				builder.WriteString(fmt.Sprintf("                 %s\n", bold.Render("No upgrade in progress.")))
				height -= 2
			} else {
				// upgrade in progress
				builder.WriteString(fmt.Sprintf("%s\n", bold.Render("Consensus Upgrade Pending")))
				builder.WriteString(fmt.Sprintf("Current Protocol: %s\n", formatVersion(string(m.Status.LastVersion))))
				builder.WriteString(fmt.Sprintf("Next Protocol:    %s\n", formatVersion(string(m.Status.NextVersion))))
				builder.WriteString(fmt.Sprintf("Upgrade round:    %s\n", nextVersion))
				height -= 4
			}
		}
	}

	// pad the box
	for height > 0 {
		builder.WriteString("\n")
		height--
	}

	return m.style.Status.Render(builder.String())
}
