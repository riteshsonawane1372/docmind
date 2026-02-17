package chunker

import (
	"strings"
)

type Chunk struct {
	Content  string
	Source   string
	ChunkIdx int64
}

// ChunkMarkdown splits markdown content into heading-aware paragraph chunks.
// It splits on double newlines, accumulates up to maxSize characters with
// overlap characters of overlap, and prepends the nearest heading to each chunk.
func ChunkMarkdown(content, source string, maxSize, overlap int) []Chunk {
	paragraphs := strings.Split(content, "\n\n")

	var chunks []Chunk
	var currentHeading string
	var buf strings.Builder
	idx := int64(0)

	flush := func() {
		text := strings.TrimSpace(buf.String())
		if text == "" {
			return
		}
		if currentHeading != "" && !strings.HasPrefix(text, currentHeading) {
			text = currentHeading + "\n\n" + text
		}
		chunks = append(chunks, Chunk{
			Content:  text,
			Source:   source,
			ChunkIdx: idx,
		})
		idx++
	}

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		if strings.HasPrefix(para, "#") {
			currentHeading = para
		}

		// If adding this paragraph would exceed maxSize, flush and start new chunk with overlap.
		if buf.Len() > 0 && buf.Len()+len(para)+2 > maxSize {
			flush()

			// Create overlap from end of previous buffer.
			prev := buf.String()
			buf.Reset()
			if overlap > 0 && len(prev) > overlap {
				overlapText := prev[len(prev)-overlap:]
				// Start overlap at a word boundary.
				if spaceIdx := strings.Index(overlapText, " "); spaceIdx >= 0 {
					overlapText = overlapText[spaceIdx+1:]
				}
				buf.WriteString(overlapText)
			}
		}

		if buf.Len() > 0 {
			buf.WriteString("\n\n")
		}
		buf.WriteString(para)
	}

	flush()
	return chunks
}
