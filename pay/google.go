/*
File Name:  google.go
Description:
Author:      Chenghu
Date:       2023/8/23 22:48
Change Activity:
*/
package pay

import (
	"context"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
	"log/slog"
)

var GooglePay *androidpublisher.Service

func initPay() (err error) {
	ctx := context.Background()
	GooglePay, err = androidpublisher.NewService(ctx, option.WithAPIKey("AIza..."))

	if err != nil {
		slog.Error("google pay client", "error", err.Error())
		return err
	}
	//退订
	//GooglePay.Orders.Refund()
	return nil
}
