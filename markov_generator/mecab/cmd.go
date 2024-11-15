package mecab

import (
	"bufio"
	"fmt"
	"io"
	. "markov_generator/domain"
	"os"
	"os/exec"
	"strings"

	"github.com/gookit/slog"
)

type MeCab struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	scanner *bufio.Scanner
}

func CreateInstance() *MeCab {
	// TODO: 変えられるようにする
	homeDir, _ := os.UserHomeDir()
	dicPath := homeDir + "/.local/share/mecab/dic/unidic-waka"
	cmd := exec.Command("mecab", "-d", dicPath)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(stdout)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	instance := &MeCab{cmd, stdin, stdout, scanner}

	slog.Info(fmt.Sprintf("MeCab process created (PID: %d)", cmd.Process.Pid))
	return instance
}

func (m *MeCab) Close() {
	m.stdin.Close()
	m.stdout.Close()
	m.cmd.Wait()
	slog.Info(fmt.Sprintf("MeCab process closed (PID: %d)", m.cmd.Process.Pid))
}

func (m *MeCab) Exec(text string) []*Morpheme {
	fmt.Fprintln(m.stdin, text)
	arr := []*Morpheme{}
	for m.scanner.Scan() {
		line := m.scanner.Text()
		if line == "EOS" {
			break
		}
		morpheme, err := parseResult(line)
		if err != nil {
			slog.Error("failed to parse MeCab result", "line", line, "error", err)
			return arr
		}
		arr = append(arr, morpheme)
	}
	return arr
}

func parseResult(line string) (*Morpheme, error) {
	arr := strings.Fields(line)
	if len(arr) != 2 {
		err := fmt.Errorf("failed to parse MeCab result (%s)", line)
		return nil, err
	}
	rawWord := arr[0]
	detail := strings.Split(arr[1], ",")

	if len(detail) < 5 {
		err := fmt.Errorf("failed to parse MeCab result (%s)", line)
		return nil, err
	} else if len(detail) == 6 {
		return &Morpheme{
			Surface: Surface(rawWord),
			Reading: Reading(rawWord),
		}, nil
	}

	return &Morpheme{
		Surface: Surface(rawWord),
		Reading: Reading(detail[9]),
	}, nil
}
