package chainobject

type GlobalPropertyObject struct {
	ID                       uint64 //global obj id
	ProposedScheduleBlockNum uint32
	ProposedSchedule         ProducerSchedule
}

func NewGlobalPropertyObject() *GlobalPropertyObject {
	return &GlobalPropertyObject{
		ID: GetOID(GlobalPropertyType),
	}
}
