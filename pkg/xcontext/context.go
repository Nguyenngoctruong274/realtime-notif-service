package xcontext

import (
	"context"
	"fmt"
	"strings"
	"time"
	"yes4all/ads-noti-api/pkg/middleware/auth/auth_model"
	"yes4all/ads-noti-api/pkg/utils/constants"

	"github.com/gin-gonic/gin"
)

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

type clone struct {
	ctx context.Context
}

func (clone) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (clone) Done() <-chan struct{} {
	return nil
}

func (clone) Err() error {
	return nil
}

func (d clone) Value(key interface{}) interface{} {
	return d.ctx.Value(key)
}

func Clone(ctx context.Context) context.Context {
	return clone{ctx: ctx}
}

const (
	KeyContextID ContextKey = "context_id"
	KeyUserID    ContextKey = "user_id"
)

var AllKeys = []ContextKey{
	KeyContextID,
	KeyUserID,
}

func AttachContext(c context.Context, key ContextKey, value string) context.Context {
	return context.WithValue(c, key, value)
}

func AttachGinContext(c *gin.Context, key ContextKey, value string) *gin.Context {
	c.Set(key.String(), value)
	return c
}

func GetUser(ctx context.Context) (res auth_model.InsideAuthClaim, err error) {
	ctxUser := ctx.Value(auth_model.InsideUser)
	res, ok := ctxUser.(auth_model.InsideAuthClaim)
	if !ok {
		err = fmt.Errorf("user không hợp lệ")
		return
	}
	return
}

func GetUserName(ctx context.Context) (res string, err error) {
	ctxUser := ctx.Value(auth_model.InsideUser)
	user, ok := ctxUser.(auth_model.InsideAuthClaim)
	if !ok {
		err = fmt.Errorf("user không hợp lệ")
		return
	}
	res = user.Email
	return
}

func IsAdmin(ctx context.Context) (res bool) {
	ctxPortfolios := ctx.Value(auth_model.AwsPortfolioIDsEdit)
	portfoliosIDs, _ := ctxPortfolios.([]string)
	if len(portfoliosIDs) == 0 {
		return true
	}

	ctxAdmin := ctx.Value(auth_model.Admin)
	admin, _ := ctxAdmin.(string)
	if admin == auth_model.EmailAdmin {
		res = true
		return
	}
	return
}

func IsAdminYams(ctx context.Context) (res bool) {
	ctxAdmin := ctx.Value(auth_model.AdminYAMS)
	admin, _ := ctxAdmin.(string)
	if len(admin) > 0 {
		res = true
		return
	}
	return
}

func GetProfileID(ctx context.Context) (res string, err error) {
	ctxProfileID := ctx.Value(auth_model.AwsProfileID)
	res, ok := ctxProfileID.(string)
	if !ok {
		err = fmt.Errorf("profile id không tồn tại")
		return
	}
	return
}

func GetParamTypeAPI(ctx context.Context) (res string) {
	typeParam := ctx.Value(constants.TypeAPIParam)
	res, ok := typeParam.(string)
	if !ok {
		return ""
	}
	return
}

func SetProfileID(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, auth_model.AwsProfileID, value)
}

func GetUserEmail(ctx context.Context) (res string, err error) {
	ctxUser := ctx.Value(auth_model.InsideUser)
	user, ok := ctxUser.(auth_model.InsideAuthClaim)
	if !ok {
		err = fmt.Errorf("user không hợp lệ")
		return
	}
	res = strings.ToLower(user.Email)
	return
}

func GetPortfolioIDs(ctx context.Context) (res []string, err error) {
	ctxPortfolioIDs := ctx.Value(auth_model.AwsPortfolioIDs)
	res, _ = ctxPortfolioIDs.([]string)
	return
}

func GetPortfolioIDsEdit(ctx context.Context) (res []string, err error) {
	ctxPortfolioIDs := ctx.Value(auth_model.AwsPortfolioIDsEdit)
	res, _ = ctxPortfolioIDs.([]string)
	return
}

func SetUserName(ctx context.Context, value string) context.Context {
	user := auth_model.InsideAuthClaim{
		Email: value,
	}
	return context.WithValue(ctx, auth_model.InsideUser, user)
}

func SetFromSource(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, auth_model.FromSource, value)
}

func GetFromSource(ctx context.Context) (res string, err error) {
	ctxFromSource := ctx.Value(auth_model.FromSource)
	res, ok := ctxFromSource.(string)
	if !ok {
		err = fmt.Errorf("FromSource không tồn tại")
		return
	}
	return
}

func CloneByGinContext(c *gin.Context) (ctx context.Context) {
	ctx = context.Background()
	for key, value := range c.Keys {
		ctx = context.WithValue(ctx, key, value)
	}
	return
}
