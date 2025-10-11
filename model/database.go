package model

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	CommonFieldId = "id"
)

var (
	_ core.RecordProxy = (*User)(nil)
	_ core.RecordProxy = (*Config)(nil)
	_ core.RecordProxy = (*Article)(nil)
)

const (
	DbNameUsers               = "users"
	UsersFieldEmail           = "email"
	UsersFieldEmailVisibility = "emailVisibility"
	UsersFieldVerified        = "verified"
	UsersFieldName            = "name"
	UsersFieldNickname        = "nickname"
	UsersFieldAvatar          = "avatar"
	UsersFieldOId             = "oId"
	UsersFieldCreated         = "created"
	UsersFieldUpdated         = "updated"
)

type User struct {
	core.BaseRecordProxy
}

func NewUser(record *core.Record) *User {
	user := new(User)
	user.SetProxyRecord(record)
	return user
}

func NewUserFromCollection(collection *core.Collection) *User {
	record := core.NewRecord(collection)
	return NewUser(record)
}

func (user *User) Name() string {
	return user.GetString(UsersFieldName)
}

func (user *User) SetName(value string) {
	user.Set(UsersFieldName, value)
}

func (user *User) Nickname() string {
	return user.GetString(UsersFieldNickname)
}

func (user *User) SetNickname(value string) {
	user.Set(UsersFieldNickname, value)
}

func (user *User) Avatar() string {
	return user.GetString(UsersFieldAvatar)
}

func (user *User) SetAvatar(value string) {
	user.Set(UsersFieldAvatar, value)
}

func (user *User) OId() string {
	return user.GetString(UsersFieldOId)
}

func (user *User) SetOId(value string) {
	user.Set(UsersFieldOId, value)
}

func (user *User) Created() types.DateTime {
	return user.GetDateTime(UsersFieldCreated)
}

func (user *User) Updated() types.DateTime {
	return user.GetDateTime(UsersFieldUpdated)
}

const (
	DbNameConfigs     = "configs"
	ConfigsFieldKey   = "key"
	ConfigsFieldValue = "value"
)

type Config struct {
	core.BaseRecordProxy
}

func NewConfig(record *core.Record) *Config {
	config := new(Config)
	config.SetProxyRecord(record)
	return config
}

func NewConfigFromCollection(collection *core.Collection) *Config {
	record := core.NewRecord(collection)
	return NewConfig(record)
}

func (config *Config) Key() ConfigKey {
	return MustParseConfigKey(config.GetString(ConfigsFieldKey))
}

func (config *Config) SetKey(value ConfigKey) {
	config.Set(ConfigsFieldKey, value)
}

func (config *Config) Value() string {
	return config.GetString(ConfigsFieldValue)
}

func (config *Config) SetValue(value string) {
	config.Set(ConfigsFieldValue, value)
}

const (
	DbNameArticles              = "articles"
	ArticlesFieldUserId         = "userId"
	ArticlesFieldOId            = "oId"
	ArticlesFieldTitle          = "title"
	ArticlesFieldPreviewContent = "previewContent"
	ArticlesFieldViewCount      = "viewCount"
	ArticlesFieldGoodCnt        = "goodCnt"
	ArticlesFieldCommentCount   = "commentCount"
	ArticlesFieldCollectCnt     = "collectCnt"
	ArticlesFieldThankCnt       = "thankCnt"
	ArticlesFieldCreatedAt      = "createdAt"
	ArticlesFieldUpdatedAt      = "updatedAt"
	ArticlesFieldCreated        = "created"
	ArticlesFieldUpdated        = "updated"
)

type Article struct {
	core.BaseRecordProxy
}

func NewArticle(record *core.Record) *Article {
	article := new(Article)
	article.SetProxyRecord(record)
	return article
}

func NewArticleFromCollection(collection *core.Collection) *Article {
	record := core.NewRecord(collection)
	return NewArticle(record)
}

func (article *Article) UserId() string {
	return article.GetString(ArticlesFieldUserId)
}

func (article *Article) SetUserId(value string) {
	article.Set(ArticlesFieldUserId, value)
}

func (article *Article) OId() string {
	return article.GetString(ArticlesFieldOId)
}

func (article *Article) SetOId(value string) {
	article.Set(ArticlesFieldOId, value)
}

func (article *Article) Title() string {
	return article.GetString(ArticlesFieldTitle)
}

func (article *Article) SetTitle(value string) {
	article.Set(ArticlesFieldTitle, value)
}

func (article *Article) PreviewContent() string {
	return article.GetString(ArticlesFieldPreviewContent)
}

func (article *Article) SetPreviewContent(value string) {
	article.Set(ArticlesFieldPreviewContent, value)
}

func (article *Article) ViewCount() int {
	return article.GetInt(ArticlesFieldViewCount)
}

func (article *Article) SetViewCount(value int) {
	article.Set(ArticlesFieldViewCount, value)
}

func (article *Article) GoodCnt() int {
	return article.GetInt(ArticlesFieldGoodCnt)
}

func (article *Article) SetGoodCnt(value int) {
	article.Set(ArticlesFieldGoodCnt, value)
}

func (article *Article) CommentCount() int {
	return article.GetInt(ArticlesFieldCommentCount)
}

func (article *Article) SetCommentCount(value int) {
	article.Set(ArticlesFieldCommentCount, value)
}

func (article *Article) CollectCnt() int {
	return article.GetInt(ArticlesFieldCollectCnt)
}

func (article *Article) SetCollectCnt(value int) {
	article.Set(ArticlesFieldCollectCnt, value)
}

func (article *Article) ThankCnt() int {
	return article.GetInt(ArticlesFieldThankCnt)
}

func (article *Article) SetThankCnt(value int) {
	article.Set(ArticlesFieldThankCnt, value)
}

func (article *Article) CreatedAt() types.DateTime {
	return article.GetDateTime(ArticlesFieldCreatedAt)
}

func (article *Article) SetCreatedAt(value types.DateTime) {
	article.Set(ArticlesFieldCreatedAt, value)
}

func (article *Article) UpdatedAt() types.DateTime {
	return article.GetDateTime(ArticlesFieldUpdatedAt)
}

func (article *Article) SetUpdatedAt(value types.DateTime) {
	article.Set(ArticlesFieldUpdatedAt, value)
}

func (article *Article) Created() types.DateTime {
	return article.GetDateTime(ArticlesFieldCreated)
}

func (article *Article) Updated() types.DateTime {
	return article.GetDateTime(ArticlesFieldUpdated)
}
