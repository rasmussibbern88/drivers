package region

type Settings interface {
	JoinRequestChannel() Channel
	JoinAcceptChannel() Channel
	UplinkChannel() Channel
	Rx2Channel() Channel
}

type settings struct {
	joinRequestChannel Channel
	joinAcceptChannel  Channel
	uplinkChannel      Channel
	rx2Channel         Channel
}

func (r *settings) JoinRequestChannel() Channel {
	return r.joinRequestChannel
}

func (r *settings) JoinAcceptChannel() Channel {
	return r.joinAcceptChannel
}

func (r *settings) UplinkChannel() Channel {
	return r.uplinkChannel
}

func (r *settings) Rx2Channel() Channel {
	return r.rx2Channel
}
