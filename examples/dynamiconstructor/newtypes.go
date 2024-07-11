package main

import "time"

// Note is a simple structure to hold an annotation
type Note struct {
	Message string
	Time    time.Time
}

// NewAnnotation is a constructor for Note
func NewAnnotation(message string, currentTime time.Time) Note {
	return Note{
		Message: message,
		Time:    currentTime,
	}
}

// FolderChecker is an implementation of FileFolderChecker
type FolderChecker struct {
	Path string
}

// NewFolderChecker is a constructor for FolderChecker
func NewFolderChecker(note Note) FolderChecker {
	return FolderChecker{
		Path: "/absolute/path/" + note.Message,
	}
}

// RunningAbsolute returns the absolute path
func (fc FolderChecker) RunningAbsolute() string {
	return fc.Path
}
