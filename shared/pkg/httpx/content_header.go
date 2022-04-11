package httpx

import (
	"strconv"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

// returns required fields for an Attachment
func GetContentHeader(url string) (*ContentHeader, error) {
	var contentHeader = &ContentHeader{
		SizeInByte: 0,
		MIMEType:   "",
	}

	header, err := Head(url)

	if err != nil {
		logger.Errorf("cannot read content type of url: %s. error is : %v", url, err)
	}

	contentLength := header.Get("Content-Length")

	if contentLength != "" {
		if sizeInBytes, err := strconv.Atoi(contentLength); err != nil {
			logger.Errorf("cannot convert content length of url: %s to int: %s. error is : %v", url, contentLength, err)
		} else {
			contentHeader.SizeInByte = sizeInBytes
		}
	} else {
		contentHeader.SizeInByte = 0
	}

	contentHeader.MIMEType = header.Get("Content-Type")

	return contentHeader, nil
}
