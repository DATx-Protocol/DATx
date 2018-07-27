package chainobject

import (
	"datx_chain/utils/rlp"
	"errors"
	"fmt"
	"log"
)

//GetIDFromRaw get global id of this raw object
func GetIDFromRaw(objtype uint32, v interface{}) uint64 {
	var id uint64

	//return 0 if the input v is invalid
	defer func() {
		if err := recover(); err != nil {
			id = 0
		}
	}()

	switch objtype {
	case AccountType:
		id = v.(*AccountObject).ID
	case GlobalPropertyType:
		id = v.(*GlobalPropertyObject).ID
	case TranscationType:
		id = v.(*TransactionObject).ID
	default:
		log.Printf("getid on undefined value type={%v}", objtype)
		id = 0
	}

	return id
}

//GetIDFromRLPData get hash id of rlp object
func GetIDFromRLPData(objtyp uint32, v []byte) uint64 {
	var id uint64

	switch objtyp {
	case AccountType:
		var value AccountObject
		if err := rlp.DecodeBytes(v, &value); err != nil {
			log.Printf("decode rlp AccountObject err={%v}", err)
			return 0
		}
		id = value.ID
	case GlobalPropertyType:
		var value GlobalPropertyObject
		if err := rlp.DecodeBytes(v, &value); err != nil {
			log.Printf("decode rlp GlobalPropertyObject err={%v}", err)
			return 0
		}
		id = value.ID
	case TranscationType:
		var value TransactionObject
		if err := rlp.DecodeBytes(v, &value); err != nil {
			log.Printf("decode rlp TransactionObject err={%v}", err)
			return 0
		}
		id = value.ID
	default:
		log.Printf("getid on undefined value type={%v}", objtyp)
	}

	return id
}

//GetValue get value from RLP data by obj type
func GetValue(objtyp uint32, data []byte) (interface{}, error) {

	switch objtyp {
	case AccountType:
		var value AccountObject
		if err := rlp.DecodeBytes(data, &value); err != nil {
			return nil, err
		}
		return value, nil
	case GlobalPropertyType:
		var value GlobalPropertyObject
		if err := rlp.DecodeBytes(data, &value); err != nil {
			return nil, err
		}
		return value, nil
	case TranscationType:
		var value TransactionObject
		if err := rlp.DecodeBytes(data, &value); err != nil {
			return nil, err
		}
		return value, nil
	default:
		str := fmt.Sprintf("Find on unreg value type={%v}", objtyp)
		log.Printf(str)
		return nil, errors.New(str)
	}
}
