package main

import (
	"os"
	"strings"
)

type liste struct {
	Static_dir     string
	Title          string
	Raw_body       string
	Processed_body map[string][]element
}

func (l *liste) processBody() {
	var carte = make(map[string][]element)
	var menu string
	for idx, line := range strings.Split(l.Raw_body, "\n") {
		line = strings.TrimSpace(line)
		if len(line) > 1 && line[0] != '#' && strings.TrimSpace(line) != "" {
			if len(line) > 2 && line[0] == '=' {
				menu = strings.TrimSpace(line[1:])
				carte[menu] = make([]element, 0)
			} else if len(line) > 2 && line[0] == '-' {
				carte[menu] = append(carte[menu], element{Index: idx, Valeur: strings.TrimSpace(line[1:])})
			}
		}
	}
	l.Processed_body = carte
}

func (l *liste) loadListe() error {
	filename := conf.resourceDir + "/" + l.Title + ".txt"
	stat, err := os.Stat(filename)
	if err != nil {
		return (err)
	}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	buf := make([]byte, stat.Size())
	if _, err := file.Read(buf); err != nil {
		return err
	}
	l.Raw_body = strings.TrimSpace(string(buf))
	return nil
}

func (l *liste) saveListe() error {
	filename := conf.resourceDir + "/" + l.Title + ".txt"
	file, err := os.OpenFile(filename, os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		return err
	}
	l.Raw_body = strings.TrimSpace(l.Raw_body)
	n, err := file.WriteString(l.Raw_body)
	if err != nil || n != len(l.Raw_body) {
		return err
	}
	if err := file.Truncate(int64(n)); err != nil { // truncate end of file if input is shorter than previous file
		return err
	}
	return nil
}


