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
	"os/exec"
	"path"
	"path/filepath"
	"runtime/debug"
	"strings"
)

func GetCurrentPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	return path
}

//open file and return file data, support arbitrary type of file on config dir
func GetFileHelper(config string) (error, []byte) {
	//catch execption
	defer func() {
		if err := recover(); err != nil {
			es := fmt.Sprintf("GetHelperConfig panic. config={%s} err={%s}", config, err)
			log.Print(es)
		}
	}()

	//get exec file current path
	CurrentPath := GetCurrentPath()
	if len(CurrentPath) == 0 {
		return errors.New(fmt.Sprintf("Parse config={%s} error", config)), nil
	}
	//check suffix name of config file
	if suffix := path.Ext(config); strings.Compare(suffix, ".yaml") != 0 {
		return errors.New(fmt.Sprintf("config={%s} suffix is not yaml", config)), nil
	}

	ins := strings.Split(CurrentPath, string(os.PathSeparator))

	ps := append(ins[:len(ins)-2], "config", config)
	configpath := strings.Join(ps, string(os.PathSeparator))
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
func CatchException(errhandle func()) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("\n***********catch exception: \n%v\n\n", err)

			str, ok := err.(string)
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

	return nil
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
