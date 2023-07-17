package tsvmarshaller

import (
	"strconv"

	"demoscraper/internal/core/entities"
)

type TSVMarshaller struct{}

func New() TSVMarshaller {
	return TSVMarshaller{}
}

func (r TSVMarshaller) Marshal(entry entities.CrawlEntry) ([]byte, error) {
	return []byte(entry.URL() + "\t" + strconv.Itoa(entry.Depth)), nil
}
