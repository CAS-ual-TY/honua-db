package honuadb

import (
	"bufio"
	"os"
	"strings"
)

func read_and_parse_sql_file(filepath string) ([]string, error) {
	var statements []string

	// Open File
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Scanner zum Zeilenweisen Lesen der Datei erstellen
	scanner := bufio.NewScanner(file)

	// Variable zum Zwischenspeichern von mehrzeiligen Statements
	var statementBuilder strings.Builder

	// Zeilenweise Datei lesen
	for scanner.Scan() {
		line := scanner.Text()

		// Wenn die Zeile mit einem Kommentar beginnt, überspringen
		if strings.HasPrefix(strings.TrimSpace(line), "--") {
			continue
		}

		// Wenn die Zeile ein Teil eines mehrzeiligen Statements ist,
		// an den Builder anhängen
		if strings.HasSuffix(strings.TrimSpace(line), ";") && statementBuilder.Len() > 0 {
			statementBuilder.WriteString(" ")
			statementBuilder.WriteString(strings.TrimSpace(line))
			statement := statementBuilder.String()
			statements = append(statements, statement)
			statementBuilder.Reset()
		} else {
			// Ansonsten die Zeile an den Builder anhängen
			statementBuilder.WriteString(" ")
			statementBuilder.WriteString(strings.TrimSpace(line))
		}
	}

	// Letztes Statement hinzufügen, falls mehrzeiliges Statement am Ende
	if statementBuilder.Len() > 0 {
		statement := statementBuilder.String()
		statements = append(statements, strings.TrimSpace(statement))
	}

	// Fehler beim Scanner prüfen
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return statements, nil
}

// Determine if a string s is in the string array a
func string_array_contains_string(s string, a []string) bool {
	for _, k := range a {
		if s == k {
			return true
		}
	}
	return false
}
