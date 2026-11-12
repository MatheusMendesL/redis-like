package protocol

import (
	"fmt"
	"redis-project/internal/domain/command"
	"strings"
)

type Command struct {
	Type command.Type
	Args []string
}

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseCommand(line string) (*Command, error) {
	line = strings.TrimSpace(line)

	if line == "" {
		return nil, fmt.Errorf("empty command")
	}

	parts := strings.Fields(line)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	cmdType := command.Type(strings.ToUpper(parts[0]))

	if !cmdType.IsValid() {
		return nil, fmt.Errorf("unknown command: %s", parts[0])
	}

	cmd := &Command{
		Type: cmdType,
		Args: []string{},
	}

	cmd.Args = parts[1:]

	return cmd, nil
}

func (p *Parser) FormatResponse(result any) string {
	switch v := result.(type) {
	case string:
		return v
	case int, int64:
		return fmt.Sprintf("%d", v)
	case bool:
		if v {
			return p.FormatOK()
		}
		return "ERR operation failed"
	case error:
		return p.FormatError(v.Error())
	default:
		return fmt.Sprintf("%v", result)
	}
}

func (p *Parser) FormatOK() string { return "OK" }

func (p *Parser) FormatError(msg string) string {
	return fmt.Sprintf("ERR: %s", msg)
}

func (p *Parser) FormatNil() string {
	return "nil"
}
