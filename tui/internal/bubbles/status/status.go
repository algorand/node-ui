package status

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"

	"github.com/algorand/node-ui/messages"
	"github.com/algorand/node-ui/tui/internal/bubbles/explorer"
	"github.com/algorand/node-ui/tui/internal/style"
)

const roundTo = time.Second / 10

// consensus constants, in theory these could be modified by a consensus upgrade.
const (
	upgradeVoteRounds = 10000
	upgradeThreshold  = 9000
)

// Model representing the status.
type Model struct {
	Status  models.NodeStatus
	Header  types.BlockHeader
	Network messages.NetworkMsg
	Err     error

	style     *style.Styles
	requestor *messages.Requestor

	// fast catchup state
	progress          progress.Model
	processedAcctsPct float64
	verifiedAcctsPct  float64
	acquiredBlksPct   float64

	// round time calculation state
	startBlock  uint64
	startTime   time.Time
	latestBlock uint64
	latestTime  time.Time
}

// New creates a status Model.
func New(style *style.Styles, requestor *messages.Requestor) Model {
	return Model{
		style:     style,
		progress:  progress.New(progress.WithDefaultGradient()),
		requestor: requestor,
	}
}

// Init is part of the tea.Model interface.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.requestor.GetNetworkCmd(),
		m.requestor.GetStatusCmd(),
	)
}

func (m Model) averageBlockTime() time.Duration {
	numBlocks := int64(m.latestBlock - m.startBlock)

	// Default round time during first seen block
	if numBlocks == 0 {
		return 4400 * time.Millisecond
	}

	runtime := m.latestTime.Sub(m.startTime)
	dur := runtime.Nanoseconds() / numBlocks
	return time.Duration(dur)
}

// Update is part of the tea.Model interface.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case explorer.BlocksMsg:
		// Still initializing.
		if m.Status.LastRound == 0 {
			return m, nil
		}

		for _, blk := range msg.Blocks {
			if uint64(blk.Block.Block.Round) == m.Status.LastRound {
				m.Header = blk.Block.Block.BlockHeader
			}
		}
		return m, nil
	case messages.StatusMsg:
		if msg.Error != nil {
			m.Err = fmt.Errorf("error fetching status: %w", msg.Error)
			return m, tea.Quit
		}
		m.Status = msg.Status

		// Save the times for computing round time
		if m.latestBlock < m.Status.LastRound {
			m.latestBlock = m.Status.LastRound
			m.latestTime = time.Now().Add(-time.Duration(m.Status.TimeSinceLastRound))

			// Grab the start time
			if m.startBlock == 0 {
				m.startBlock = m.Status.LastRound
				since := time.Duration(m.Status.TimeSinceLastRound)
				m.startTime = time.Now().Add(-since)
			}
		}

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
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func formatVersion(v string) string {
	i := strings.LastIndex(v, "/")
	if i != 0 {
		i++
	}
	return v[i:]
}

func writeProgress(b *strings.Builder, prefix string, progress progress.Model, pct float64) {
	b.WriteString(prefix)
	b.WriteString(progress.ViewAs(pct))
	b.WriteString("\n")
}

func (m Model) calculateTimeToGo(start, end uint64, style lipgloss.Style) string {
	rounds := end - start
	timeRemaining := time.Duration(int64(rounds) * m.averageBlockTime().Nanoseconds()).Round(roundTo)
	return style.Render(fmt.Sprintf("%d to go, %s", rounds, timeRemaining))
}

// View is part of the tea.Model interface.
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
			builder.WriteString(fmt.Sprintf("Block wait time: %s\n", time.Duration(m.Status.TimeSinceLastRound).Round(roundTo)))
			builder.WriteString(fmt.Sprintf("Sync time:       %s\n", time.Duration(m.Status.CatchupTime).Round(roundTo)))
			height -= 3
			if m.Header.UpgradeState != (types.UpgradeState{}) {
				//remainingToUpgrade := m.calculateTimeToGo(
				//	m.Status.LastRound, uint64(m.Header.NextProtocolSwitchOn), m.style.AccountBlueText)
				remainingToVote := m.calculateTimeToGo(
					m.Status.LastRound, uint64(m.Header.NextProtocolVoteBefore), m.style.AccountBlueText)

				// calculate yes/no votes
				votesToGo := uint64(m.Header.NextProtocolVoteBefore) - m.Status.LastRound
				votes := upgradeVoteRounds - votesToGo
				voteYes := m.Header.NextProtocolApprovals
				voteNo := votes - voteYes
				voteString := fmt.Sprintf("%d / %d", voteYes, voteNo)
				yesPct := float64(voteYes) / float64(votes)
				windowPct := float64(votes) / float64(upgradeVoteRounds)
				builder.WriteString(fmt.Sprintf("%s\n", bold.Render("Consensus Upgrade Pending: Votes")))
				builder.WriteString(fmt.Sprintf("Next Protocol:     %s\n", formatVersion(m.Header.NextProtocol)))
				builder.WriteString(fmt.Sprintf("Yes/No votes:      %s (%.0f%%, 90%% required)\n", voteString, yesPct*100))
				//builder.WriteString(fmt.Sprintf("Vote window:      %s (%f%%)\n", voteString, *100))
				builder.WriteString(fmt.Sprintf("Vote window close: %d (%.0f%%, %s)\n",
					m.Header.UpgradeState.NextProtocolVoteBefore,
					windowPct*100,
					remainingToVote))

				height -= 5
			} else if m.Status.LastVersion == m.Status.NextVersion {
				// no upgrade in progress
				builder.WriteString(fmt.Sprintf("Protocol:        %s\n", formatVersion(m.Status.LastVersion)))
				builder.WriteString(fmt.Sprintf("                 %s\n", bold.Render("No upgrade in progress.")))
				height -= 2
			} else {
				// compute the time until the upgrade round and apply formatting to message
				togo := m.Status.NextVersionRound - m.Status.LastRound
				timeRemaining := time.Duration(int64(togo) * m.averageBlockTime().Nanoseconds()).Round(roundTo)
				remaining := m.style.AccountBlueText.Render(
					fmt.Sprintf("%d to go, %s", togo, timeRemaining))

				// upgrade in progress
				builder.WriteString(fmt.Sprintf("%s\n", bold.Render("Consensus Upgrade Scheduled")))
				builder.WriteString(fmt.Sprintf("Current Protocol: %s\n", formatVersion(m.Status.LastVersion)))
				builder.WriteString(fmt.Sprintf("Next Protocol:    %s\n", formatVersion(m.Status.NextVersion)))
				builder.WriteString(fmt.Sprintf("Upgrade round:    %d (%s)\n", m.Status.NextVersionRound, remaining))
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
