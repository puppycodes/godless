package godless

import "sort"

// FIXME not really a stream, whole is kept in memory.
type NamespaceStreamEntry struct {
	Table  TableName
	Row    RowName
	Entry  EntryName
	Points []Point
}

func ReadNamespaceEntryMessage(message *NamespaceEntryMessage) NamespaceStreamEntry {
	entry := NamespaceStreamEntry{
		Table:  TableName(message.Table),
		Row:    RowName(message.Row),
		Entry:  EntryName(message.Entry),
		Points: make([]Point, len(message.Points)),
	}

	for i, p := range message.Points {
		entry.Points[i] = Point(p)
	}

	return entry
}

func MakeNamespaceEntryMessage(entry NamespaceStreamEntry) *NamespaceEntryMessage {
	pb := &NamespaceEntryMessage{
		Table:  string(entry.Table),
		Row:    string(entry.Row),
		Entry:  string(entry.Entry),
		Points: make([]string, len(entry.Points)),
	}

	for i, p := range entry.Points {
		pb.Points[i] = string(p)
	}

	return pb
}

type byNamespaceStreamOrder []NamespaceStreamEntry

func (stream byNamespaceStreamOrder) Len() int {
	return len(stream)
}

func (stream byNamespaceStreamOrder) Swap(i, j int) {
	stream[i], stream[j] = stream[j], stream[i]
}

func (stream byNamespaceStreamOrder) Less(i, j int) bool {
	a, b := stream[i], stream[j]

	if a.Table < b.Table {
		return true
	} else if a.Table > b.Table {
		return false
	}

	if a.Row < b.Row {
		return true
	} else if a.Row > b.Row {
		return false
	}

	if a.Entry < b.Entry {
		return true
	} else if a.Entry > b.Entry {
		return false
	}

	minSize := imin(len(a.Points), len(b.Points))

	for i := 0; i < minSize; i++ {
		ap := a.Points[i]
		bp := b.Points[j]
		if ap < bp {
			return true
		} else if ap > bp {
			return false
		}
	}

	return len(a.Points) < len(b.Points)
}

func MakeStreamEntry(tname TableName, rname RowName, ename EntryName, entry Entry) NamespaceStreamEntry {
	return NamespaceStreamEntry{
		Table:  tname,
		Row:    rname,
		Entry:  ename,
		Points: entry.GetValues(),
	}
}

func MakeNamespaceStream(ns Namespace) []NamespaceStreamEntry {
	stream := []NamespaceStreamEntry{}
	for tableName, table := range ns.Tables {
		for rowName, row := range table.Rows {
			for entryName, entry := range row.Entries {
				if len(entry.Set) > 0 {
					streamEntry := MakeStreamEntry(tableName, rowName, entryName, entry)
					stream = append(stream, streamEntry)
				}
			}
		}
	}

	sort.Sort(byNamespaceStreamOrder(stream))
	return stream
}

func MakeNamespaceStreamMessage(stream []NamespaceStreamEntry) *NamespaceMessage {
	message := &NamespaceMessage{Entries: make([]*NamespaceEntryMessage, len(stream))}

	for i, entry := range stream {
		message.Entries[i] = MakeNamespaceEntryMessage(entry)
	}

	return message
}

func ReadNamespaceStreamMessage(message *NamespaceMessage) []NamespaceStreamEntry {
	stream := make([]NamespaceStreamEntry, len(message.Entries))

	for i, emsg := range message.Entries {
		stream[i] = ReadNamespaceEntryMessage(emsg)
	}

	return stream
}

func ReadNamespaceStream(stream []NamespaceStreamEntry) Namespace {
	ns := EmptyNamespace()

	for _, streamEntry := range stream {
		ns.addStreamEntry(streamEntry)
	}

	return ns
}