package utils

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// ExitOnError : Terminate Program with Error
func ExitOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

// CreateFile : Create New File
func CreateFile(path string) string {
	f, err := os.Create(path)
	ExitOnError(err)
	defer f.Close()

	return path
}

// PromptConfirm : Prompt for Confirmation
func PromptConfirm() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	if strings.ToLower(string(rune(response[0]))) == "y" {
		return true
	}
	if strings.ToLower(string(rune(response[0]))) == "n" {
		return false
	}
	fmt.Printf("Please type y (for yes) or n (for no) and then press enter: ")
	return PromptConfirm()
}

// CleanText ...
func CleanText(str string) string {
	// Return if empty
	if len(str) < 0 {
		return ""
	}

	cleaned := ""
	// Regexp for leading space (\r), new line (\n) and tab (\t)
	re := regexp.MustCompile(`\r|\t|\n`)
	cleaned = re.ReplaceAllString(str, " ")
	// Convert all chars to Uppercase
	// Uppercase is prefered against Lowercase because on some langeuages (ie. Greek) lowercase ending tokens (ie. σ, ς) are different from the uppercase
	// But yet again languages like Greek have accents
	// Need more testing here
	cleaned = strings.TrimSpace(cleaned)
	// Remove Regexp
	// cleaned = re.ReplaceAllString(cleaned, " ")
	// Remove single quote, this is important
	// Greek single quote will not split
	// cleaned = strings.Replace(cleaned, "'", "", -1)
	cleaned = strings.Replace(cleaned, "\"", "'", -1)
	// TODO
	// Remove links
	// Remove quoted text and pointer if exists (ie. βόλφγκανγκ σόιμπλε: «ο αριθμός των υποστηρικτών μιας ελάφρυνσης χρέους για την ελλάδα αυξάνεται, αφότου η αθήνα επικύρωσε ένα νέο μεταρρυθμιστικό πακέτο. ωστόσο, κάποιος δεν θέλει ακόμη να πάει μαζί τους»)
	// Remove HTML
	htmlre := regexp.MustCompile(`<[^>]*>`)

	// Space
	space := regexp.MustCompile(`\s+`)
	cleaned = space.ReplaceAllString(cleaned, " ")
	cleaned = htmlre.ReplaceAllString(cleaned, "")

	mm := regexp.MustCompile(`(?m)(^"|"\r?$|";")|"`)
	cleaned = mm.ReplaceAllString(cleaned, "${1}")

	// commre := regexp.MustCompile(`/\/\*[\s\S]*?\*\/|([^\\:]|^)\/\/.*|<!--[\s\S]*?-->$/`)
	// cleaned = commre.ReplaceAllString(cleaned, "")
	// Retrun cleaned text
	return cleaned
}
