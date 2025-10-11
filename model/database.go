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
	_ core.RecordProxy = (*Awards)(nil)
	_ core.RecordProxy = (*Histories)(nil)
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

const (
	DbNameAwards           = "awards"
	AwardsFieldLevel       = "level"
	AwardsFieldName        = "name"
	AwardsFieldAlias       = "alias"
	AwardsFieldPoint       = "point"
	AwardsFieldAmount      = "amount"
	AwardsFieldDescription = "description"
)

type Awards struct {
	core.BaseRecordProxy
}

func NewAwards(record *core.Record) *Awards {
	award := new(Awards)
	award.SetProxyRecord(record)
	return award
}

func NewAwardsFromCollection(collection *core.Collection) *Awards {
	record := core.NewRecord(collection)
	return NewAwards(record)
}

func (award *Awards) Level() int {
	return award.GetInt(AwardsFieldLevel)
}

func (award *Awards) SetLevel(value int) {
	award.Set(AwardsFieldLevel, value)
}

func (award *Awards) Name() string {
	return award.GetString(AwardsFieldName)
}

func (award *Awards) SetName(value string) {
	award.Set(AwardsFieldName, value)
}

func (award *Awards) Alias() string {
	return award.GetString(AwardsFieldAlias)
}

func (award *Awards) SetAlias(value string) {
	award.Set(AwardsFieldAlias, value)
}

func (award *Awards) Point() int {
	return award.GetInt(AwardsFieldPoint)
}

func (award *Awards) SetPoint(value int) {
	award.Set(AwardsFieldPoint, value)
}

func (award *Awards) Amount() int {
	return award.GetInt(AwardsFieldAmount)
}

func (award *Awards) SetAmount(value int) {
	award.Set(AwardsFieldAmount, value)
}

func (award *Awards) Description() string {
	return award.GetString(AwardsFieldDescription)
}

func (award *Awards) SetDescription(value string) {
	award.Set(AwardsFieldDescription, value)
}

const (
	DbNameHistories       = "histories"
	HistoriesFieldUserId  = "userId"
	HistoriesFieldTimes   = "times"
	HistoriesFieldAwardId = "awardId"
	HistoriesFieldIsTop   = "isTop"
	HistoriesFieldIsBest  = "isBest"
	HistoriesFieldDetails = "details"
	HistoriesFieldCreated = "created"
	HistoriesFieldUpdated = "updated"
)

type Histories struct {
	core.BaseRecordProxy
}

func NewHistories(record *core.Record) *Histories {
	history := new(Histories)
	history.SetProxyRecord(record)
	return history
}

func NewHistoriesFromCollection(collection *core.Collection) *Histories {
	record := core.NewRecord(collection)
	return NewHistories(record)
}
func (history *Histories) UserId() string {
	return history.GetString(HistoriesFieldUserId)
}

func (history *Histories) SetUserId(value string) {
	history.Set(HistoriesFieldUserId, value)
}

func (history *Histories) Times() int {
	return history.GetInt(HistoriesFieldTimes)
}

func (history *Histories) SetTimes(value int) {
	history.Set(HistoriesFieldTimes, value)
}

func (history *Histories) AwardId() string {
	return history.GetString(HistoriesFieldAwardId)
}

func (history *Histories) SetAwardId(value string) {
	history.Set(HistoriesFieldAwardId, value)
}

func (history *Histories) IsTop() bool {
	return history.GetBool(HistoriesFieldIsTop)
}

func (history *Histories) SetIsTop(value bool) {
	history.Set(HistoriesFieldIsTop, value)
}

func (history *Histories) IsBest() bool {
	return history.GetBool(HistoriesFieldIsBest)
}

func (history *Histories) SetIsBest(value bool) {
	history.Set(HistoriesFieldIsBest, value)
}

func (history *Histories) Details() [6]int {
	var details = types.JSONArray[int]{}
	_ = details.Scan(history.GetRaw(HistoriesFieldDetails))

	var arr [6]int
	for i := 0; i < 6 && i < len(details); i++ {
		arr[i] = details[i]
	}

	return arr
}

func (history *Histories) SetDetails(value [6]int) {
	history.Set(HistoriesFieldDetails, value)
}

func (history *Histories) Created() types.DateTime {
	return history.GetDateTime(HistoriesFieldCreated)
}

func (history *Histories) Updated() types.DateTime {
	return history.GetDateTime(HistoriesFieldUpdated)
}
