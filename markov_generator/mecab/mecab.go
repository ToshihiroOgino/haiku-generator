package mecab

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type MeCab struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	scanner *bufio.Scanner
}

func CreateInstance() *MeCab {
	home_dir, _ := os.UserHomeDir()
	dic_path := home_dir + "/.local/share/mecab/dic/unidic-waka"
	cmd := exec.Command("mecab", "-d", dic_path)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(stdout)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	instance := &MeCab{cmd, stdin, stdout, scanner}
	return instance
}

func (m *MeCab) Parse(text string) string {
	fmt.Fprintln(m.stdin, text)
	out := ""
	for m.scanner.Scan() {
		line := m.scanner.Text()
		if line == "EOS" {
			break
		}
		out += line + "\n"
	}
	return out
}

func (m *MeCab) Close() {
	m.stdin.Close()
	m.stdout.Close()
	m.cmd.Wait()
}
