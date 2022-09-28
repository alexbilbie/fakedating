package model

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Location struct {
	Latitude  float64
	Longitude float64
}

// Scan decodes a custom type by the database driver
func (g *Location) Scan(src interface{}) error {
	// Found via https://antonibertel.medium.com/golang-scanning-mysql-mariadb-geographical-point-type-2c425c4f8bc
	switch b := src.(type) {
	case []byte:
		if len(b) != 25 {
			return fmt.Errorf("expected []bytes with length 25, got %d", len(b))
		}
		var longitude float64
		var latitude float64
		buf := bytes.NewReader(b[9:17])
		err := binary.Read(buf, binary.LittleEndian, &longitude)
		if err != nil {
			return err
		}
		buf = bytes.NewReader(b[17:25])
		err = binary.Read(buf, binary.LittleEndian, &latitude)
		if err != nil {
			return err
		}
		*g = Location{Latitude: latitude, Longitude: longitude}
	default:
		return fmt.Errorf("expected []byte for Location type, got  %T", src)
	}
	return nil
}
