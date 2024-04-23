package lorawan

import (
	"errors"

	"tinygo.org/x/drivers/lora"
	"tinygo.org/x/drivers/lora/lorawan/region"
)

var (
	ErrNoJoinAcceptReceived    = errors.New("no JoinAccept packet received")
	ErrNoRadioAttached         = errors.New("no LoRa radio attached")
	ErrInvalidEuiLength        = errors.New("invalid EUI length")
	ErrInvalidAppKeyLength     = errors.New("invalid AppKey length")
	ErrInvalidPacketLength     = errors.New("invalid packet length")
	ErrInvalidDevAddrLength    = errors.New("invalid DevAddr length")
	ErrInvalidMic              = errors.New("invalid Mic")
	ErrFrmPayloadTooLarge      = errors.New("FRM payload too large")
	ErrInvalidNetIDLength      = errors.New("invalid NetID length")
	ErrInvalidNwkSKeyLength    = errors.New("invalid NwkSKey length")
	ErrInvalidAppSKeyLength    = errors.New("invalid AppSKey length")
	ErrUndefinedRegionSettings = errors.New("undefined Regionnal Settings ")
)

const (
	LORA_TX_TIMEOUT = 2000
	LORA_RX_TIMEOUT = 10000
)

var (
	ActiveRadio    lora.Radio
	Retries        = 15
	regionSettings region.Settings
)

// UseRegionSettings sets current Lorawan Regional parameters
func UseRegionSettings(rs region.Settings) {
	regionSettings = rs
}

// UseRadio attaches Lora radio driver to Lorawan
func UseRadio(r lora.Radio) {
	if ActiveRadio != nil {
		panic("lorawan.ActiveRadio is already set")
	}
	ActiveRadio = r
}

// SetPublicNetwork defines Lora Sync Word according to network type (public/private)
func SetPublicNetwork(enabled bool) {
	ActiveRadio.SetPublicNetwork(enabled)
}

// ApplyChannelConfig sets current Lora modulation according to current regional settings
func applyChannelConfig(ch region.Channel) {
	ActiveRadio.SetFrequency(ch.Frequency())
	ActiveRadio.SetBandwidth(ch.Bandwidth())
	ActiveRadio.SetCodingRate(ch.CodingRate())
	ActiveRadio.SetSpreadingFactor(ch.SpreadingFactor())
	ActiveRadio.SetPreambleLength(ch.PreambleLength())
	ActiveRadio.SetTxPower(ch.TxPowerDBm())
	// Lorawan defaults to explicit headers
	ActiveRadio.SetHeaderType(lora.HeaderExplicit)
	ActiveRadio.SetCrc(true)
}

// Join tries to connect Lorawan Gateway
func Join(otaa *Otaa, session *Session) error {
	var resp []uint8

	if ActiveRadio == nil {
		return ErrNoRadioAttached
	}

	if regionSettings == nil {
		return ErrUndefinedRegionSettings
	}

	otaa.Init()

	// Send join packet
	payload, err := otaa.GenerateJoinRequest()
	if err != nil {
		return err
	}

	for {
		joinRequestChannel := regionSettings.JoinRequestChannel()
		joinAcceptChannel := regionSettings.JoinAcceptChannel()

		// Prepare radio for Join Tx
		applyChannelConfig(joinRequestChannel)
		ActiveRadio.SetIqMode(lora.IQStandard)
		ActiveRadio.Tx(payload, LORA_TX_TIMEOUT)
		if err != nil {
			return err
		}

		// Wait for JoinAccept
		if joinAcceptChannel.Frequency() != 0 {
			applyChannelConfig(joinAcceptChannel)
		}
		ActiveRadio.SetIqMode(lora.IQInverted)
		resp, err = ActiveRadio.Rx(LORA_RX_TIMEOUT)
		if err == nil && resp != nil {
			break
		}
		if !joinAcceptChannel.Next() {
			return ErrNoJoinAcceptReceived
		}
	}

	err = otaa.DecodeJoinAccept(resp, session)
	if err != nil {
		return err
	}

	applyChannelConfig(regionSettings.UplinkChannel()) // JUST STUPID...
	return nil
}

// SendUplink sends Lorawan Uplink message
func SendUplink(data []uint8, session *Session) error {
	if regionSettings == nil {
		return ErrUndefinedRegionSettings
	}

	payload, err := session.GenMessage(0, 0x01, []byte(data))
	if err != nil {
		return err
	}

	// applyChannelConfig(regionSettings.UplinkChannel())
	ActiveRadio.SetIqMode(lora.IQStandard)
	err = ActiveRadio.Tx(payload, LORA_TX_TIMEOUT)
	if err != nil {
		return err
	}
	return nil
}

// SendUplink sends Lorawan Uplink message //Small change in GenMassage -> might be better to fix elsewhere.
func SendUplinkMac(data []uint8, session *Session) error {
	if regionSettings == nil {
		return ErrUndefinedRegionSettings
	}

	// Need to build the correct payload and calculate the mic...

	payload, err := session.GenMessage(0, 0x00, []byte(data))
	if err != nil {
		return err
	}

	// applyChannelConfig(regionSettings.UplinkChannel())
	ActiveRadio.SetIqMode(lora.IQStandard)
	ActiveRadio.Tx(payload, LORA_TX_TIMEOUT)
	if err != nil {
		return err
	}
	return nil
}

// Listen downlink, window 1 for RX1 window config (uplink), everything else is RX2 window config (worstcase)
func ListenDownlink(window int) ([]uint8, error) {
	if regionSettings == nil {
		return nil, ErrUndefinedRegionSettings
	}

	if window == 1 {
		// applyChannelConfig(regionSettings.UplinkChannel())
		ActiveRadio.SetIqMode(lora.IQInverted)
		resp, err := ActiveRadio.Rx(LORA_RX_TIMEOUT)
		if err == nil && resp != nil {
			return resp, nil
		}
	} else {
		// applyChannelConfig(regionSettings.Rx2Channel()) //TODO Fix this channel and reset the normal one after...
		ActiveRadio.SetIqMode(lora.IQInverted)
		resp, err := ActiveRadio.Rx(LORA_RX_TIMEOUT)
		if err == nil && resp != nil {
			return resp, nil
		}
	}

	return nil, nil
}
