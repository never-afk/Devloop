package util

var DefaultTransMessage = map[string]map[string]string{
	"alphanum": {
		"en": "must be alphanum",
		"zh": "必须是字母或者数字",
	},
	"base64": {
		"en": "must be eq {0}",
		"zh": "必须等于{0}",
	},
	"contains": {
		"en": "must contains {0}",
		"zh": "必须等于{0}",
	},
	"email": {
		"en": "invalid email address",
		"zh": "无效的邮箱地址",
	},
	"gt": {
		"en": "must gt {0}",
		"zh": "必须大于{0}",
	},
	"gte": {
		"en": "must gte {0}",
		"zh": "必须大于等于{0}",
	},
	"ip": {
		"en": "invalid IP address",
		"zh": "无效的IP地址",
	},
	"len": {
		"en": "length must be {0}",
		"zh": "长度必须是{0}位",
	},
	"lt": {
		"en": "must lt {0}",
		"zh": "必须小于{0}",
	},
	"lte": {
		"en": "must lte {0}",
		"zh": "必须小于等于{0}",
	},
	"max": {
		"en": "max is {0}",
		"zh": "最大值是{0}",
	},
	"min": {
		"en": "min is {0}",
		"zh": "最小值是{0}",
	},
	"numeric": {
		"zh": "必须是数字",
		"en": "must numeric",
	},
	"oneof": {
		"en": "only allow in [{0}]",
		"zh": "只能是以下值: [{0}]",
	},
	"required": {
		"zh": "不能为空",
		"en": "can not be empty",
	},
	"url": {
		"en": "invalid url address",
		"zh": "无效的URL地址",
	},
}
