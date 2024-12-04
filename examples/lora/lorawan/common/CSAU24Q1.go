//go:build CSAU24Q1

package common

import "machine"

var (
	rstPin = machine.PA9
	//csPin   = machine.PB12 // NSS3
	csPin   = machine.SPI_NSS3
	dio0Pin = machine.PA11
	dio1Pin = machine.PA12
	spi     = machine.SPI0
)

//TODO func setup_rfm_io()
