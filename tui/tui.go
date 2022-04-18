package tui

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/algorand/node-ui/tui/internal/model"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
)

const host = "0.0.0.0"

func teaHandler(_ ssh.Session) (tea.Model, []tea.ProgramOption) {
	return model.New(), []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseCellMotion()}
}

// Start ...
func Start(port uint64) {
	if port == 0 {
		// Run directly
		p := tea.NewProgram(model.New(), tea.WithAltScreen(), tea.WithMouseCellMotion())
		if err := p.Start(); err != nil {
			fmt.Printf("Error in UI: %v", err)
			os.Exit(1)
		}

		fmt.Printf("\nUI Terminated, shutting down node.\n")
		os.Exit(0)
	}

	// Run on ssh server.
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	sshServer, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(path.Join(dirname, ".ssh/term_info_ed25519")),
		wish.WithMiddleware(
			bm.Middleware(teaHandler),
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting SSH server on %s:%d", host, port)
	go func() {
		if err = sshServer.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := sshServer.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
