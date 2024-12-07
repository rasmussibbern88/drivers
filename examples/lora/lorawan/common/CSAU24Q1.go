//go:build CSAU24Q1

package common

import (
	"machine"

	"tinygo.org/x/drivers"
)

var (
	rstPin = machine.PA9
	//csPin   = machine.PB12 // NSS3
	csPin               = machine.SPI_NSS1
	dio0Pin             = machine.PA11
	dio1Pin             = machine.PA12
	spi     drivers.SPI = nil
)

func SetSPI(newSpi drivers.SPI) {
	spi = newSpi
}

//TODO func setup_rfm_io()
