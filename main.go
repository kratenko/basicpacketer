package basicpacketer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type Writer struct {
	writer         io.Writer
	bytesWritten   uint64
	packetsWritten uint64
}

func newWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer}
}

func (writer *Writer) Write(data []byte) (err error) {
	var size uint32 = uint32(len(data))
	binary.Write(writer.writer, binary.BigEndian, &size)
	writer.bytesWritten += 4
	written, err := writer.writer.Write(data)
	if err != nil {
		return err
	}
	if int64(written) != int64(size) {
		return errors.New(fmt.Sprintf("wrote wrong number of bytes. did: %d, should: %d", written, size))
	}
	writer.bytesWritten += uint64(written)
	writer.packetsWritten += 1
	return nil
}

// --

type Reader struct {
	reader      io.Reader
	bytesRead   uint64
	packetsRead uint64
}

func newReader(reader io.Reader) *Reader {
	return &Reader{reader: reader}
}

func (reader *Reader) Read() (data []byte, err error) {
	var size uint32
	binary.Read(reader.reader, binary.BigEndian, &size)
	reader.bytesRead += 4
	buffer := make([]byte, size)
	read, err := reader.reader.Read(buffer)
	reader.bytesRead += uint64(read)
	if err != nil {
		return nil, err
	}
	if int64(read) < int64(size) {
		return nil, errors.New(fmt.Sprintf("read wrong number of bytes. did: %d, should: %d", read, size))
	}
	reader.packetsRead += 1
	return buffer, nil
}
