package migration

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const scanBufSize = 4 * 1024 * 1024

// implemet only simple variant migration.
func ParseSQLMigration(f io.Reader) ([]string, error) {
	scanBuf := make([]byte, scanBufSize)

	scanner := bufio.NewScanner(f)
	scanner.Buffer(scanBuf, scanBufSize)

	// res[0] is Up transactions, res[1] is Down transactions.
	res := make([]string, 2)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Up") {
			for scanner.Scan() {
				line = scanner.Text()
				if strings.Contains(line, ";") {
					res[0] += " " + line
					break
				}
				res[0] += " " + line
			}
		} else if strings.Contains(line, "Down") {
			for scanner.Scan() {
				line = scanner.Text()
				if strings.Contains(line, ";") {
					res[1] += " " + line
					break
				}
				res[1] += " " + line
			}
		}
	}

	// check strings.
	if (len(res[0]) > 0 && !strings.Contains(res[0], ";")) ||
		(len(res[1]) > 0 && !strings.Contains(res[1], ";")) {
		return res, fmt.Errorf("rows for database transactions do not contain the ;")
	}
	return res, nil
}
