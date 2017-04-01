package main

import ()

type ConfigIO interface {
	Deserialize(reader io.Reader) (TaskList, error)
	Serialize(writer io.Writer, tasks TaskList) error
}
