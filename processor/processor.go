package processor

import (
	"encoding/json"
	"log"
)

type ParserJSON struct {
	ResultChannel chan bool
	ErrorChannel  chan error 
}

func NewParserJSON(resultChannel chan bool, errorChannel chan error) *ParserJSON {
	return &ParserJSON{
		ResultChannel: resultChannel,
		ErrorChannel:  errorChannel, 
	}
}

func (p *ParserJSON) Run(value []byte, attempts []ValidationFunc[string]) bool {
	var data map[string]interface{}
	if err := json.Unmarshal(value, &data); err != nil {
		log.Println("Error parsing JSON:", err)
		p.ErrorChannel <- err 
		return false
	}
	log.Println("Parsed data:", data)

	if len(attempts) == 0 {
		p.ResultChannel <- true
	}

	return true
}
