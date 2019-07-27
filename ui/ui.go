//ui.go
package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadUI() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter API key: ")
	text, _ := reader.ReadString('\n')
	// convert CRLF to LF
	text = strings.Replace(text, "\n", "", -1)

	fmt.Println("API key is ", text)
}
