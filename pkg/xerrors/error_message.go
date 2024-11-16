package xerrors

const (
	MsgSuccess = "Thành công!"
	MsgFail    = "Không thành công!"

	// General Message
	MsgGeneralError       = "Có lỗi xảy ra trong quá trình xử lý, vui lòng thực hiện lại sau!"
	MsgBadRequest         = "Yêu cầu không hợp lệ!"
	MsgAuthenticateFailed = "Không thể xác thực tài khoản!"

	MsgNoProfileID         string = "Profile không hợp lệ"
	MsgNoPortfolioAccepted string = "No Portfolio accepted"
	MsgNoRoleEditPortfolio string = "you need role edit that portfolio"
	MsgUserAlreadyThisRole string = "user already has this permission"
	MsgNoRoleEditYamsAdmin string = "you need role edit that YAMS admin"
)
