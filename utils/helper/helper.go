package helper

import (
	"datx_chain/utils/common"
	"datx_chain/utils/crypto/sha3"
	"datx_chain/utils/rlp"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strings"
)

func GetCurrentPath() string {
	file, _ := os.Getwd()

	path, _ := filepath.Abs(file)
	// path := filepath.Dir(file)
	log.Printf("get current path:%v", path)
	if len(path) == 0 {
		return path
	}

	ins := strings.Split(path, string(os.PathSeparator))
	return strings.Join(ins[:], string(os.PathSeparator))
}

func MakePath(path ...string) string {
	var first string
	if len(path) < 1 {
		return first
	}
	first = path[0]
	for i := 1; i < len(path); i++ {
		first = first + string(os.PathSeparator) + path[i]
	}

	return first
}

//open file and return file data, support arbitrary type of file on config dir
func GetFileHelper(config string, configFolder string) (error, []byte) {
	//catch execption
	defer func() {
		if err := recover(); err != nil {
			es := fmt.Sprintf("GetHelperConfig panic. config={%s} err={%s}", config, err)
			log.Print(es)
		}
	}()

	//check suffix name of config file
	if suffix := path.Ext(config); strings.Compare(suffix, ".yaml") != 0 {
		return errors.New(fmt.Sprintf("config={%s} suffix is not yaml", config)), nil
	}

	//get exec file current path
	CurrentPath := GetCurrentPath()
	if len(CurrentPath) == 0 {
		return errors.New(fmt.Sprintf("Parse config={%s} error", config)), nil
	}
	log.Printf("current path:%v", CurrentPath)
	in := strings.Index(CurrentPath, "datx_chain")
	if in < 1 {
		return errors.New(fmt.Sprintf("Parse config={%s} index error, not include datx_chain.", config)), nil
	}

	configpath := MakePath(CurrentPath[:in-1], "datx_chain", configFolder, config)
	log.Printf("Config={%s} path={%s}", config, configpath)

	//check file exist or not
	_, err := os.Stat(configpath)
	if err != nil && os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("config={%s} %s ", config, err)), nil
	}

	//open file
	data, err := ioutil.ReadFile(configpath)
	if err != nil {
		return errors.New(fmt.Sprintf("config={%s} %s ", config, err)), nil
	}

	return nil, data
}

//catch exception, call errhandle func to release resource if there are exist panic
func CatchException(err error, errhandle func()) {
	defer func() {
		if cerr := recover(); cerr != nil {
			log.Printf("\n***********catch exception: \n%v\n\n", cerr)

			str, ok := cerr.(string)
			if ok {
				err = errors.New(str)
			} else {
				err = errors.New("panic")
			}

			//err handle,release resource
			errhandle()

			//print stack info
			debug.PrintStack()
		}
	}()
}

func RLPHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func ToBytes(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))

	return b
}

func EncodeBit(k1, k2 uint32) uint64 {
	pair := uint64(k1)<<32 | uint64(k2)
	return pair
}

func DecodeBit(pair uint64) (uint32, uint32) {
	k1 := uint32(pair >> 32)
	k2 := uint32(pair) & 0xFFFFFFFF
	return k1, k2
}
