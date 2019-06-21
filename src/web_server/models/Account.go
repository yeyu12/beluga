package models

import (
	"github.com/astaxie/beego/orm"
	"time"
	"encoding/json"
	"beluga/src/web_server/helpers"
	"strconv"
	"math"
)

type Account struct {
	Id               int       `json:"id"`                // 主键id
	Token            string    `json:"token"`             // token
	Username         string    `json:"username"`          // 用户名
	Passwd           string    `json:"passwd"`            // 密码
	Nickname         string    `json:"nickname"`          // 昵称
	Avatar           string    `json:"avatar"`            // 头像
	Gender           byte      `json:"gender"`            // 性别，0未知，1男，2女
	Age              byte      `json:"age"`               // 年龄
	Status           byte      `json:"status"`            // 账号状态（1正常，0禁止
	CreateTime       time.Time `json:"create_time"`       // 创建时间
	TokenExpiryTime  time.Time `json:"token_expiry_time"` // token过期时间
	ConfigurationNum int       `json:"configuration_num"` // 配置项目数
}

// 当前对象
func NewAccount() *Account {
	return &Account{}
}

// 加入表前缀的表名
func (m *Account) TableNamePrefix() string {
	return helpers.GetTablePrefix() + "account"
}

// 用户user_id 查找用户信息
func (m *Account) UseridFind(o orm.Ormer, id int) (Account, error) {
	var account Account
	err := o.QueryTable(m.TableNamePrefix()).Filter("id", id).One(&account)

	return account, err
}

// 用户昵称获取用户信息
func (m *Account) UserNicknameToInfo(o orm.Ormer, nickname string) (Account, error) {
	var account Account
	err := o.QueryTable(m.TableNamePrefix()).Filter("nickname", nickname).One(&account)

	return account, err
}

// 用户名与密码查找用户信息
func (m *Account) UsernameFind(o orm.Ormer, username string, passwd string) (*Account, error) {
	account := &Account{}
	err := o.QueryTable(m.TableNamePrefix()).Filter("username", username).Filter("passwd", passwd).One(account)

	return account, err
}

// 用户名获取用户数组
func (m *Account) UsernameToData(o orm.Ormer, username string) (*Account, error) {
	account := &Account{}
	err := o.QueryTable(m.TableNamePrefix()).Filter("username", username).One(account)

	return account, err
}

// token 获取用户信息
func (m *Account) TokenToUser(o orm.Ormer, token string) (*Account, error) {
	account := &Account{}
	err := o.QueryTable(m.TableNamePrefix()).Filter("token", token).One(account)

	return account, err
}

// 更新用户数据
func (m *Account) Update(o orm.Ormer, data *Account) bool {
	obj := o.QueryTable(m.TableNamePrefix())
	var account_data map[string]interface{}
	account_json, _ := json.Marshal(data)
	json.Unmarshal(account_json, &account_data)

	_, err := obj.Filter("id", data.Id).Update(account_data)

	if err == nil {
		return true
	}

	return false
}

// 账号修改用户数据
func (m *Account) IdSetUserInfo(o orm.Ormer, username string, data map[string]interface{}) bool {
	_, err := o.QueryTable(m.TableNamePrefix()).Filter("username", username).Update(data)

	if err == nil {
		return true
	}

	return false
}

// 添加用户
func (m *Account) Add(o orm.Ormer) (int64, error) {
	id, err := o.Insert(m)

	return id, err
}

// 删除用户

// 获取用户列表
func (m *Account) List(o orm.Ormer, page int, page_size int, search string) helpers.Page {
	page_size_str := strconv.Itoa(page_size)

	total_obj := o.QueryTable(m.TableNamePrefix())

	if search != "" {
		total_obj = total_obj.Filter("nickname__icontains", search)
	}

	total, _ := total_obj.Count()

	offset := strconv.Itoa(((page - 1) * page_size))
	sql := "select id,username,nickname,status,configuration_num,create_time from `" + m.TableNamePrefix() + "`"

	if search != "" {
		sql += " where `nickname` like '%" + search + "%' "
	}
	sql += " order by id desc limit " + offset + "," + page_size_str

	var lists []orm.Params
	o.Raw(sql).Values(&lists)

	account_page_data := helpers.Page{
		Total:     total,
		TotalPage: math.Ceil(float64(total) / float64(page_size)),
		PageSize:  page_size,
		Page:      page,
		List:      lists,
	}

	return account_page_data
}
