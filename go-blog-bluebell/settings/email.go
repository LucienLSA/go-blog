package settings

const (
	EmailOperationBinding = iota + 1
	EmailOperationNoBinding
)

var EmailOperationMap = map[int]string{
	EmailOperationBinding:   "您正在绑定邮箱, 请点击链接确定身份 %s",
	EmailOperationNoBinding: "您正在解邦邮箱, 请点击链接确定身份 %s",
}
