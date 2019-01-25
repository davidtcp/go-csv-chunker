package main

import (
  "os"
  "testing"

  "github.com/stretchr/testify/assert"
)

func Test_fileCloser(t *testing.T) {
  type args struct {
    file *os.File
  }
  existFile, _ := os.Open(`/dev/null`)
  notExistFile, _ := os.Open(`/dev/fake`)
  tests := []struct {
    name      string
    args      args
    wantPanic bool
  }{
    {
      name:      `exist`,
      args:      args{existFile},
      wantPanic: false,
    },
    {
      name:      `not exist`,
      args:      args{notExistFile},
      wantPanic: true,
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      if tt.wantPanic {
        assert.Panics(t, func() { fileCloser(tt.args.file) }, "The code did not panic")
      }
    })
  }
}

func Test_chunkFilenamePrefix(t *testing.T) {
  tests := []struct {
    name                    string
    wantChunkFilenamePrefix string
  }{
    {
      name:                    `from env`,
      wantChunkFilenamePrefix: `part_`,
    },
    {
      name:                    `from default`,
      wantChunkFilenamePrefix: `chunk_`,
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      _ = os.Unsetenv(`CSV_CHUNKER_FILENAME_PREFIX`)
      if tt.name == `from env` {
        _ = os.Setenv(`CSV_CHUNKER_FILENAME_PREFIX`, `part_`)
      }

      if gotChunkFilenamePrefix := chunkFilenamePrefix(); gotChunkFilenamePrefix != tt.wantChunkFilenamePrefix {
        t.Errorf("chunkFilenamePrefix() = %v, want %v", gotChunkFilenamePrefix, tt.wantChunkFilenamePrefix)
      }
    })
  }
}

func Test_isArgumentsValid(t *testing.T) {
  tests := []struct {
    name      string
    wantValid bool
    args      []string
  }{
    {
      name:      `valid`,
      wantValid: true,
      args:      []string{`binary`, `foo`, `bar`},
    },
    {
      name:      `less`,
      wantValid: false,
      args:      []string{`binary`, `foo`},
    },
    {
      name:      `more`,
      wantValid: false,
      args:      []string{`binary`, `foo`, `bar`, `quz`},
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      os.Args = tt.args
      if gotValid := isArgumentsValid(); gotValid != tt.wantValid {
        t.Errorf("isArgumentsValid() = %v, want %v", gotValid, tt.wantValid)
      }
    })
  }
}
