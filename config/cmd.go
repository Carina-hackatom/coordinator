package config

import (
	"fmt"
	"github.com/Carina-hackatom/coordinator/utils"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
)

type IBCPort struct {
	Transfer string
}

type IBCChan struct {
	Nova IBCPort
	Host IBCPort
}

type FlagOpts struct {
	Test   bool
	New    bool
	Disp   bool
	ExtIP  string
	Kn     string
	Host   string
	Period int
	IBCChan
}

func SetInitialDir(krDir string, logDir string) (string, string) {
	cwd, err := os.Getwd()
	utils.CheckErr(err, "cannot get working directory", 0)

	krDir = path.Join(cwd, "/bot", krDir)
	err = os.MkdirAll(krDir, 0740)
	if os.IsExist(err) {
		log.Println("** bot directory already exist **")
	} else if err != nil {
		log.Fatal(err)
	}

	logDir = path.Join(cwd, logDir)
	err = os.MkdirAll(logDir, 0740)
	if os.IsExist(err) {
		log.Println("** log directory already exist **")
	} else if err != nil {
		log.Fatal(err)
	}

	return krDir, logDir
}

func GetPassphrase(vp *viper.Viper) string {
	pw := vp.GetString("pw")
	pp := fmt.Sprintf("%s\n%s\n", pw, pw)
	return pp
}

func SetAllLogger(logDir, stdLogName, errLogName, errRedirectLogName string, isDisp bool) (*os.File, *os.File, *os.File) {
	var fdLog, fdErr, fdErrExt *os.File
	var err error
	if !isDisp {
		fdLog, err = os.OpenFile(path.Join(logDir, stdLogName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		utils.CheckErr(err, "cannot open logfp", 0)

		// 반환되서 처리할 수 있는 에러 핸들링
		fdErr, err = os.OpenFile(path.Join(logDir, errLogName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		utils.CheckErr(err, "cannot open novaerr", 0)

		// 외부 라이브러리에서 fmt.Fprintf(os.stderr)로 처리하는 애들 핸들링
		fdErrExt, err = os.OpenFile(path.Join(logDir, errRedirectLogName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		utils.CheckErr(err, "cannot open otherErr", 0)

		os.Stderr = fdErrExt
	} else {
		fdLog = os.Stdout
		fdErr = os.Stderr
		fdErrExt = os.Stderr
	}

	return fdLog, fdErr, fdErrExt
}
