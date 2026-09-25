package text

import (
	"bytes"
	"errors"
	"testing"
)

func TestGetContentFromBalnced(t *testing.T) {
	tests := []struct {
		name                 string
		input                []byte
		start                int
		expectedRes          []byte
		expectedContentStart int
		expectedContentEnd   int
		err                  error
	}{
		{
			name:                 "correctly balanced link",
			input:                []byte("some text here [valid link](https://example.com) whatever"),
			start:                len("some text here "),
			expectedRes:          []byte("valid link"),
			expectedContentStart: len("some text here [valid link]("),
			expectedContentEnd:   len("some text here [valid link](https://example.com"),
			err:                  nil,
		},
		{
			name:                 "correctly balanced link with wrong start should be treated as not a markdown structure",
			input:                []byte("[valid link](https://example.com)"),
			start:                1,
			expectedRes:          nil,
			expectedContentStart: 0,
			expectedContentEnd:   0,
			err:                  nil,
		},
		{
			name:                 "uncorrectly balanced structure",
			input:                []byte("[not valid}(../source/image.txt)"),
			start:                0,
			expectedRes:          nil,
			expectedContentStart: 0,
			expectedContentEnd:   0,
			err:                  ErrNotValidMarkdownStructureName,
		},
		{
			name:                 "correctly with a sub structure balanced structure",
			input:                []byte("[^some-random-thing](This is something else ![image here](../../source/image.png))"),
			start:                0,
			expectedRes:          []byte("^some-random-thing"),
			expectedContentStart: len("[^some-random-thing]("),
			expectedContentEnd:   len("[^some-random-thing](This is something else ![image here](../../source/image.png)"),
			err:                  nil,
		},
		{
			name:                 "uncorrectly added",
			input:                []byte("[^some-random-thing])This is something else ![image here](../../source/image.png))"),
			start:                0,
			expectedRes:          nil,
			expectedContentStart: 0,
			expectedContentEnd:   0,
			err:                  ErrNotValidMarkdownStructureContent,
		},
		{
			name:                 "correctly with uncorrect sub structure",
			input:                []byte("[^some-random-thing](This is something else ![image here}(../../source/image.png))"),
			start:                0,
			expectedRes:          []byte("^some-random-thing"),
			expectedContentStart: len("[^some-random-thing]("),
			expectedContentEnd:   len("[^some-random-thing](This is something else ![image here](../../source/image.png)"),
			err:                  nil,
		},
		{
			name:                 "uncorrectly terminated",
			input:                []byte("[^some-random-thing](This is something \nelse)"),
			start:                0,
			expectedRes:          nil,
			expectedContentStart: 0,
			expectedContentEnd:   0,
			err:                  ErrNotValidMarkdownStructureContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, contentStart, contentEnd, err := getContentFromBalanced(tt.start, tt.input)
			if !bytes.Equal(tt.expectedRes, res) {
				t.Errorf("expeced %s but got %s", tt.expectedRes, res)
			}

			if tt.expectedContentStart != contentStart {
				t.Errorf("expected content to start in %d idx, but it starts in %d", tt.expectedContentStart, contentStart)
			}

			if tt.expectedContentEnd != contentEnd {
				t.Errorf("expected content to end i %d idx, but it ends in %d", tt.expectedContentEnd, contentEnd)
			}

			if err != nil && !errors.Is(tt.err, err) {
				t.Errorf("expected content balanced to be %s but got %s", tt.err, err)
			}
		})
	}
}
