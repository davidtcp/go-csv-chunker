package main

import (
  "bufio"
  "fmt"
  "github.com/c2h5oh/datasize"
  "github.com/tit/go-sysexits"
  "os"
)

const (
  newLine                    = "\n"
  envKeyChunkFilenamePrefix  = `CSV_CHUNKER_FILENAME_PREFIX`
  defaultChunkFilenamePrefix = `chunk_`
  chunkExtension             = `.csv`
  appName                    = `csv_chunker`
)

// I know, it is a magic, but...
const (
  _ = iota // binary
  argSourceCsvFilename
  argChunkSize
  argTotal
)

func main() {
  if !isArgumentsValid() {
    usage()
  }

  sourceCsvFilePath := os.Args[argSourceCsvFilename]

  chunkSize, err := chunkSize()
  if err != nil {
    usage()
  }

  chunkFilenamePrefix := chunkFilenamePrefix()

  sourceCsvFile, _ := os.Open(sourceCsvFilePath)
  defer fileCloser(sourceCsvFile)

  scanner := bufio.NewScanner(sourceCsvFile)

  scanner.Scan()
  header := scanner.Text() + newLine

  size := len(header)
  var line string
  var counter = 0
  chunk, _ := os.Create(chunkFilenamePrefix + fmt.Sprintf("%d", counter) + chunkExtension)
  _, _ = chunk.Write([]byte(header))
  for scanner.Scan() {
    line = scanner.Text() + newLine
    if size+len(line) > int(chunkSize.Bytes()) {
      counter++
      _ = chunk.Close()
      chunk, _ = os.Create(chunkFilenamePrefix + fmt.Sprintf("%d", counter) + chunkExtension)
      _, _ = chunk.Write([]byte(header))
      size = len(header)
    }
    _, _ = chunk.Write([]byte(line))
    size += len(line)
  }
  defer fileCloser(chunk)
}

func chunkSize() (chunkSize datasize.ByteSize, err error) {
  err = chunkSize.UnmarshalText([]byte(os.Args[argChunkSize]))
  return
}

func isArgumentsValid() (valid bool) {
  if len(os.Args) == argTotal {
    return true
  }
  return false
}

func chunkFilenamePrefix() (chunkFilenamePrefix string) {
  chunkFilenamePrefix, isEnvChunkFilenamePrefixExist := os.LookupEnv(envKeyChunkFilenamePrefix)
  if !isEnvChunkFilenamePrefixExist {
    chunkFilenamePrefix = defaultChunkFilenamePrefix
  }
  return
}

func usage() () {
  message := fmt.Sprintf("Usage: %s source_file_csv chunk_size", appName)
  message += newLine
  _, _ = os.Stderr.WriteString(message)
  os.Exit(sysexits.ExUsage)
}

func fileCloser(file *os.File) {
  err := file.Close()
  if err != nil {
    fmt.Printf(err.Error())
    panic(fmt.Errorf("can not close %s file", file.Name()))
  }
}
