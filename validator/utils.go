package validator

import (
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

// CamelToSnake 驼峰命名转蛇形命名
// UsernameUnique -> username_unique
// XMLParser -> xml_parser
func CamelToSnake(s string) string {
	// 处理连续大写字母：XMLParser -> XmlParser
	re1 := regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
	s = re1.ReplaceAllString(s, "${1}_${2}")

	// 处理普通驼峰：userName -> user_name
	re2 := regexp.MustCompile(`([a-z\d])([A-Z])`)
	s = re2.ReplaceAllString(s, "${1}_${2}")

	return strings.ToLower(s)
}

// SnakeToCamel 蛇形命名转驼峰命名
// username_unique -> UsernameUnique
func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(parts, "")
}

// IsValidatorMethod 检查方法是否符合验证器签名
func IsValidatorMethod(method reflect.Method) bool {
	methodType := method.Type

	// 必须有2个参数（receiver + FieldLevel）和1个返回值
	if methodType.NumIn() != 2 || methodType.NumOut() != 1 {
		return false
	}

	// 第二个参数必须是 validator.FieldLevel
	if !methodType.In(1).Implements(reflect.TypeOf((*FieldLevel)(nil)).Elem()) {
		return false
	}

	// 返回值必须是 bool
	if methodType.Out(0).Kind() != reflect.Bool {
		return false
	}

	// 方法名必须是导出的（首字母大写）
	return unicode.IsUpper(rune(method.Name[0]))
}

// ExtractValidatorMethods 从结构体中提取验证器方法
func ExtractValidatorMethods(pkg any) map[string]ValidatorFunc {
	validators := make(map[string]ValidatorFunc)

	pkgValue := reflect.ValueOf(pkg)
	pkgType := reflect.TypeOf(pkg)

	// 注意：不要转换指针类型为值类型，因为这会丢失指针接收者的方法
	// 直接使用原始的pkgValue和pkgType来遍历方法

	// 遍历所有方法
	for i := 0; i < pkgValue.NumMethod(); i++ {
		method := pkgType.Method(i)
		methodValue := pkgValue.Method(i)

		if IsValidatorMethod(method) {
			// 转换为 ValidatorFunc 类型
			validatorFunc := func(mv reflect.Value) ValidatorFunc {
				return func(fl FieldLevel) bool {
					results := mv.Call([]reflect.Value{reflect.ValueOf(fl)})
					return results[0].Bool()
				}
			}(methodValue)

			validators[method.Name] = validatorFunc
		}
	}

	return validators
}

// BuildTagName 构建验证器标签名
func BuildTagName(methodName string, config *PackageConfig) string {
	// 检查自定义标签映射
	if config.Tags != nil {
		if customTag, exists := config.Tags[methodName]; exists {
			return customTag
		}
	}

	// 使用自定义转换器
	if config.Transform != nil {
		tagName := config.Transform(methodName)
		if config.Prefix != "" {
			return config.Prefix + "_" + tagName
		}
		return tagName
	}

	// 默认转换
	tagName := CamelToSnake(methodName)
	if config.Prefix != "" {
		return config.Prefix + "_" + tagName
	}

	return tagName
}
