//   File Name:  file.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/10/26 09:46
//    Change Activity:

package utils

import (
	"io"
	"log/slog"
	"os"
)

func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		slog.Error("打开google file 失败", "error", err.Error())
		return nil, err
	}
	defer file.Close()
	fd, err := io.ReadAll(file)
	if err != nil {
		slog.Error("read to fd fail", "error", err.Error())
		return nil, err
	}
	return fd, nil
}
