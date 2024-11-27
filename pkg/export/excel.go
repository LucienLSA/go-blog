package export

import "blog-service/global"

func GetExcelFullUrl(name string) string {
	return global.AppSetting.PrefixUrl + "/" + GetExcelPath() + name
}

func GetExcelPath() string {
	return global.AppSetting.ExportSavePath
}

func GetExcelFullPath() string {
	return global.AppSetting.RuntimeRootPath + GetExcelPath()
}
