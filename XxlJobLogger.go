package xxl_job_logger

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xxl-job/xxl-job-executor-go"
	"io"
	"os"
	"strconv"
	"time"
)

const LOG_DIR = "/data/applogs/xxl-job/jobhandler/"

type Xxljob_logger_handler struct {
	log   *logrus.Logger
	logId int64
}

func GetLogHandler() *Xxljob_logger_handler {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(io.MultiWriter(os.Stdout))
	return &Xxljob_logger_handler{
		logger,
		0,
	}
}

func (s *Xxljob_logger_handler) getLogDir() string {
	date := time.Now().Format("2006-01-02")
	if _, err := os.Stat(LOG_DIR + date); os.IsNotExist(err) {
		err = os.MkdirAll(LOG_DIR+date, 0755)
		if err != nil {
			panic(fmt.Errorf("create log directory error:%s", err.Error()))
		}
	}
	return LOG_DIR + date
}
func (s *Xxljob_logger_handler) Info(format string, a ...interface{}) {
	s.log.Infof(format, a...)
}

func (s *Xxljob_logger_handler) Error(format string, a ...interface{}) {
	s.log.Errorf(format, a...)
}

func (s *Xxljob_logger_handler) SetLogId(log_id int) *Xxljob_logger_handler {
	s.logId = int64(log_id)
	if log_id != 0 {
		s.log.SetOutput(io.MultiWriter(s.createLogFile(), os.Stdout))
	} else {
		s.log.SetOutput(os.Stdout)
	}
	return s
}
func (s *Xxljob_logger_handler) createLogFile() *os.File {

	var filename string
	if s.logId == 0 {
		filename = "-" + strconv.Itoa(int(time.Now().Hour())) + ".log"
	} else {
		filename = strconv.Itoa(int(s.logId)) + ".log"
	}
	filepath := s.getLogDir() + "/" + filename
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	return file
}
func (s *Xxljob_logger_handler) makeLogFileName(log_id int64, log_date_time int64) string {

	date_str := time.UnixMilli(log_date_time).Format("2006-01-02")

	return LOG_DIR + date_str + "/" + strconv.Itoa(int(log_id)) + ".log"
}

func (s *Xxljob_logger_handler) ReadLog(req *xxl.LogReq) *xxl.LogRes {
	log_id := s.logId
	log_dateTime := req.LogDateTim
	from_line := req.FromLineNum
	filepath := s.makeLogFileName(log_id, log_dateTime)

	file, err := os.Open(filepath)
	if err != nil {
		return &xxl.LogRes{Code: 200, Msg: "error", Content: xxl.LogResContent{
			FromLineNum: req.FromLineNum,
			ToLineNum:   2,
			LogContent:  err.Error(),
			IsEnd:       true,
		}}
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	currentLine := 0
	var log_context []byte
	for {
		currentLine++
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		if currentLine >= from_line {
			var eof byte = '\n'
			log_context = append(log_context, line...)
			log_context = append(log_context, eof)
		}
	}

	return &xxl.LogRes{Code: 200, Msg: "ok", Content: xxl.LogResContent{
		FromLineNum: req.FromLineNum,
		ToLineNum:   currentLine,
		LogContent:  string(log_context),
		IsEnd:       false,
	}}

}
