/*
 * @Author: Casso-Wong
 * @Date: 2021-09-09 14:18:34
 * @Last Modified by: Casso-Wong
 * @Last Modified time: 2021-09-10 15:24:04
 * 数据复制，专门用于请求体数据映射到DTO模型。
 * 原因：ShouldBind时需要限制字段，如果直接使用DTO模型接受请求数据则不好限制字段；
 * 又不希望每次都每个字段都显示赋值，所以封装一个基于json包的数据复制功能
 */
package datacopy

import (
	"encoding/json"
)

// DataCopy 本来想使用reflect包来完成这个功能，json.Marshal()其实也是借助reflect包来实现的.
// TIP: 少于五个字段不建议使用
func DataCopy(data, res interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, res)
}
