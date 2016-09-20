package log

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type Log struct {
	filename   string
	fileSystem FileSystemInterface
	messages   chan Communication
	screen     bool // Print message on screen
}

type Communication struct {
	message string
	done    chan error
}

////////////////////////////////////////////////////////////////////////////////
// Public functions

// Gets a new functionning log struct
func New() *Log {
	var log Log
	log.fileSystem = FileSystem{}
	log.screen = true
	return &log
}

// Start the log agent
func (log *Log) Start(filename string) error {
	if log == nil { // Check if function is used on empty pointer
		return errors.New("No log struct initialized")
	}

	log.filename = filename

	if err := log.createFile(); err != nil {
		return err
	}

	log.messages = make(chan Communication, 10)
	go log.routine()

	return nil
}

// Print logs in log system
func (log *Log) Print(line string) error {
	line = log.formatLine(line)

	if log == nil { // Check if function is used on empty pointer
		return errors.New("No log struct initialized")
	}

	if log.messages == nil { // Check if the log has been started
		return errors.New("Log system not started")
	}

	done := make(chan error, 1)
	log.messages <- Communication{line, done}
	select {
	case err := <-done:
		return err
	}
}

// Clear file containing logs
// Check if file exists and erase it
func (log *Log) Clear() error {
	if _, err := log.fileSystem.Stat(log.filename); err != nil {
		goto create
	}

	if err := log.fileSystem.Remove(log.filename); err != nil {
		return err
	}

create:
	return log.createFile()
}

////////////////////////////////////////////////////////////////////////////////
// Private functions

// Goroutine created by Start to handle new lines
func (log *Log) routine() {
	for {
		select {
		case communication := <-log.messages:
			// Print message at screen if wanted
			if log.screen {
				fmt.Print(communication.message)
			}
			// Write file
			if err := log.writeFile(communication.message); err != nil {
				communication.done <- err
			}
			communication.done <- nil
		}
	}
}

//  Format line with time
func (log *Log) formatLine(line string) string {
	return "[" + time.Now().String() + "] " + line + "\n"
}

// Create the file if it doesn't exist
func (log *Log) createFile() error {
	if _, err := log.fileSystem.Stat(log.filename); err != nil {
		goto noError
	}
	if _, err := log.fileSystem.Create(log.filename); err != nil {
		return err
	}
noError:
	return nil
}

func (log *Log) writeFile(message string) error {
	file, err := log.fileSystem.OpenFile(log.filename, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	// Write file
	_, err = log.fileSystem.WriteString(file, message)
	if err != nil {
		return err
	}
	return nil
}
