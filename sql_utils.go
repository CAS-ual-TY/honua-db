package honuadb

import (
	"bufio"
	"os"
	"strings"
)

func readAndParseSqlFile(filepath string) ([]string, error) {
	var statements []string

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var statementBuilder strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(strings.TrimSpace(line), "--") {
			continue
		}

		if strings.HasSuffix(strings.TrimSpace(line), ";") && statementBuilder.Len() > 0 {
			statementBuilder.WriteString(" ")
			statementBuilder.WriteString(strings.TrimSpace(line))
			statement := statementBuilder.String()
			statements = append(statements, statement)
			statementBuilder.Reset()
		} else {
			statementBuilder.WriteString(" ")
			statementBuilder.WriteString(strings.TrimSpace(line))
		}
	}

	if statementBuilder.Len() > 0 {
		statement := statementBuilder.String()
		statements = append(statements, strings.TrimSpace(statement))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return statements, nil
}

func string_array_contains_string(s string, a []string) bool {
	for _, k := range a {
		if s == k {
			return true
		}
	}
	return false
}
