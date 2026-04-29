package sandbox

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// isInteractive 判断 stdin 是否连接到终端，用于在 CI、AI agent、管道、后台等
// 非交互场景下避免阻塞在交互式提示上。
func isInteractive() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

// Confirm 打印 prompt 并从 stdin 读取用户输入，仅在输入为 y/Y 时返回 true。
// 非交互场景会直接打印错误提示并返回 false，调用方应据此提前返回。
func Confirm(format string, args ...any) bool {
	if !isInteractive() {
		PrintError("confirmation required but stdin is not a terminal; pass --yes to confirm in non-interactive mode")
		return false
	}
	fmt.Printf(format+" [y/N] ", args...)
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" && confirm != "Y" {
		fmt.Println("Aborted")
		return false
	}
	return true
}
