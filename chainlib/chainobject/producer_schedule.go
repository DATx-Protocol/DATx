package chainobject

import "datx_chain/utils/common"

//ProducerKey struct
type ProducerKey struct {
	ProducerName string
	SigningKey   common.Hash //rlp code hash of public key
}

//ProducerSchedule struct
type ProducerSchedule struct {
	Version   uint32
	Producers []ProducerKey
}

//NewInitSchedule new
func NewInitSchedule(prok ProducerKey) ProducerSchedule {
	pros := make([]ProducerKey, 0)
	pros = append(pros, prok)

	return ProducerSchedule{
		Version:   0,
		Producers: pros,
	}
}

//Clear clear data
func (ps ProducerSchedule) Clear() {
	ps.Version = 0
	ps.Producers = make([]ProducerKey, 0)
}

//GetProducerKey get key
func (ps ProducerSchedule) GetProducerKey(p string) common.Hash {
	for _, v := range ps.Producers {
		if v.ProducerName == p {
			return v.SigningKey
		}
	}

	return common.Hash{}
}
