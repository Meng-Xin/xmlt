package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xmlt/global"
	"xmlt/internal/domain"
	"xmlt/internal/expand/e"
	"xmlt/internal/expand/public"
	"xmlt/internal/service"
)

type UserController struct {
	service service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{service: service}
}

func (u *UserController) Register(ctx *gin.Context) {
	code := e.SUCCESS
	var vo userRegister
	err := ctx.Bind(&vo)
	if err != nil {
		global.Log.Info("userRegister Json解析失败")
		return
	}
	// 执行注册
	res, err := u.service.Register(ctx, domain.User{
		NickName: vo.NickName,
		UserName: vo.UserName,
		Password: vo.Password,
	})
	if err != nil {
		if err == e.UserNameExistedError {
			// 账户已被注册
			code = e.Registered
			ctx.JSON(http.StatusOK, public.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			})
			return
		}
		if err == e.NickNameExistedError {
			// 昵称已存在-注册失败
			code = e.UserExisted
			ctx.JSON(http.StatusOK, public.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			})
			return
		}
		// 出现问题的注册失败
		global.Log.Warn("注册失败%s", err.Error())
		code = e.ERROR
		ctx.JSON(http.StatusOK, public.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		})
		return
	}
	// 注册成功
	vo.init(res)
	ctx.JSON(code, public.Response{
		Status: code,
		Data:   vo,
	})
}

func (u *UserController) Login(ctx *gin.Context) {
	code := e.SUCCESS
	var vo userLoginReq
	err := ctx.Bind(&vo)
	if err != nil {
		code = e.ERROR
		global.Log.Info("userLogin Json解析失败")
		ctx.JSON(code, public.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		})
		return
	}
	// 登录验证
	res, err := u.service.Login(ctx, domain.User{
		UserName: vo.UserName,
		Password: vo.Password,
	})
	if err != nil {
		if err == e.UserOrPwdError {
			code = e.UserOrPasswordError
			ctx.JSON(http.StatusOK, public.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			})
			return
		}
		global.Log.Warn("登录失败：", err.Error())
		code = e.ERROR
		ctx.JSON(http.StatusOK, public.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		})
		return
	}
	data := userLoginRes{}
	data.init(res)
	ctx.JSON(code, public.Response{
		Status: code,
		Data:   data,
	})
}

type UserV0 struct {
	ID       uint64 `form:"id" json:"id"`              // UID
	UserName string `from:"user_name" json:"userName"` // 账号
	Password string `from:"password" json:"password"`  // 密码
	NickName string `from:"nick_name" json:"nickName"` // 昵称
	Email    string `from:"email" json:"email"`        // 邮箱
	Phone    string `from:"phone" json:"phone"`        // 手机号
	Avatar   string `from:"avatar" json:"avatar"`      // 头像

	Ctime string
	Utime string
}

type userRegister struct {
	ID       uint64 `form:"id" json:"id"`              // UID
	NickName string `from:"nick_name" json:"nickName"` // 昵称
	UserName string `from:"user_name" json:"userName"` // 账号
	Password string `from:"password" json:"password"`  // 密码
}

type userLoginReq struct {
	UserName string `from:"user_name" json:"userName"` // 账号
	Password string `from:"password" json:"password"`  // 密码
}

type userLoginRes struct {
	ID       uint64 `form:"id" json:"id"`              // UID
	UserName string `from:"user_name" json:"userName"` // 账号
	NickName string `from:"nick_name" json:"nickName"` // 昵称
	Avatar   string `from:"avatar" json:"avatar"`      // 头像
	Token    string `from:"token" json:"token"`        // 登陆凭证
}

func (u *UserV0) init(user domain.User) {
	u.ID = user.ID
	u.UserName = user.UserName
	u.Password = user.Password
	u.NickName = user.NickName
	u.Email = user.Email
	u.Phone = user.Phone
	u.Avatar = user.Avatar
	u.Ctime = user.Ctime.String()
	u.Utime = user.Utime.String()
}

func (u *userRegister) init(user domain.User) {
	u.ID = user.ID
	u.UserName = user.UserName
	u.Password = user.Password
	u.NickName = user.NickName
}

func (u *userLoginRes) init(user domain.User) {
	u.ID = user.ID
	u.UserName = user.UserName
	u.NickName = user.NickName
	u.Avatar = user.Avatar
	u.Token = user.Token
}
