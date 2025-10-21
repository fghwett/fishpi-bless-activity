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
	_ core.RecordProxy = (*Reward)(nil)
	_ core.RecordProxy = (*Awards)(nil)
	_ core.RecordProxy = (*Histories)(nil)
	_ core.RecordProxy = (*Vote)(nil)
	_ core.RecordProxy = (*Points)(nil)
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
	DbNameRewards      = "rewards"
	RewardsFieldLevel  = "level"
	RewardsFieldName   = "name"
	RewardsFieldPoint  = "point"
	RewardsFieldAmount = "amount"
	RewardsFiledMore   = "more"
)

type Reward struct {
	core.BaseRecordProxy
}

func NewReward(record *core.Record) *Reward {
	reward := new(Reward)
	reward.SetProxyRecord(record)
	return reward
}

func NewRewardFromCollection(collection *core.Collection) *Reward {
	record := core.NewRecord(collection)
	return NewReward(record)
}

func (reward *Reward) Level() int {
	return reward.GetInt(RewardsFieldLevel)
}

func (reward *Reward) SetLevel(value int) {
	reward.Set(RewardsFieldLevel, value)
}

func (reward *Reward) Name() string {
	return reward.GetString(RewardsFieldName)
}

func (reward *Reward) SetName(value string) {
	reward.Set(RewardsFieldName, value)
}

func (reward *Reward) Point() int {
	return reward.GetInt(RewardsFieldPoint)
}

func (reward *Reward) SetPoint(value int) {
	reward.Set(RewardsFieldPoint, value)
}

func (reward *Reward) Amount() int {
	return reward.GetInt(RewardsFieldAmount)
}

func (reward *Reward) SetAmount(value int) {
	reward.Set(RewardsFieldAmount, value)
}

func (reward *Reward) More() string {
	return reward.GetString(RewardsFiledMore)
}

func (reward *Reward) SetMore(value string) {
	reward.Set(RewardsFiledMore, value)
}

const (
	DbNameAwards           = "awards"
	AwardsFieldLevel       = "level"
	AwardsFieldName        = "name"
	AwardsFieldRewardId    = "rewardId"
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

func (award *Awards) RewardId() string {
	return award.GetString(AwardsFieldRewardId)
}

func (award *Awards) SetRewardId(value string) {
	award.Set(AwardsFieldRewardId, value)
}

func (award *Awards) Description() string {
	return award.GetString(AwardsFieldDescription)
}

func (award *Awards) SetDescription(value string) {
	award.Set(AwardsFieldDescription, value)
}

const (
	DbNameHistories         = "histories"
	HistoriesFieldUserId    = "userId"
	HistoriesFieldTimes     = "times"
	HistoriesFieldAwardId   = "awardId"
	HistoriesFieldRewardId  = "rewardId"
	HistoriesFieldIsTop     = "isTop"
	HistoriesFieldIsBest    = "isBest"
	HistoriesFieldGotReward = "gotReward"
	HistoriesFieldDetails   = "details"
	HistoriesFieldCreated   = "created"
	HistoriesFieldUpdated   = "updated"
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

func (history *Histories) RewardId() string {
	return history.GetString(HistoriesFieldRewardId)
}

func (history *Histories) SetRewardId(value string) {
	history.Set(HistoriesFieldRewardId, value)
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

func (history *Histories) GotReward() bool {
	return history.GetBool(HistoriesFieldGotReward)
}

func (history *Histories) SetGotReward(value bool) {
	history.Set(HistoriesFieldGotReward, value)
}

func (history *Histories) Details() [6]int {
	var details = types.JSONArray[int]{}
	_ = details.Scan(history.GetString(HistoriesFieldDetails))

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

const (
	DbNameVotes          = "votes"
	VotesFieldFromUserId = "fromUserId"
	VotesFieldToUserId   = "toUserId"
	VotesFieldArticleId  = "articleId"
	VotesFieldVoteType   = "voteType"
	VotesFieldCreated    = "created"
	VotesFieldUpdated    = "updated"
)

type Vote struct {
	core.BaseRecordProxy
}

func NewVote(record *core.Record) *Vote {
	vote := new(Vote)
	vote.SetProxyRecord(record)
	return vote
}

func NewVoteFromCollection(collection *core.Collection) *Vote {
	record := core.NewRecord(collection)
	return NewVote(record)
}

func (vote *Vote) FromUserId() string {
	return vote.GetString(VotesFieldFromUserId)
}

func (vote *Vote) SetFromUserId(value string) {
	vote.Set(VotesFieldFromUserId, value)
}

func (vote *Vote) ToUserId() string {
	return vote.GetString(VotesFieldToUserId)
}

func (vote *Vote) SetToUserId(value string) {
	vote.Set(VotesFieldToUserId, value)
}

func (vote *Vote) ArticleId() string {
	return vote.GetString(VotesFieldArticleId)
}

func (vote *Vote) SetArticleId(value string) {
	vote.Set(VotesFieldArticleId, value)
}

func (vote *Vote) VoteType() string {
	return vote.GetString(VotesFieldVoteType)
}

func (vote *Vote) SetVoteType(value string) {
	vote.Set(VotesFieldVoteType, value)
}

func (vote *Vote) Created() types.DateTime {
	return vote.GetDateTime(VotesFieldCreated)
}

func (vote *Vote) Updated() types.DateTime {
	return vote.GetDateTime(VotesFieldUpdated)
}

const (
	DbNamePoints         = "points"
	PointsFieldUserId    = "userId"
	PointsFieldHistoryId = "historyId"
	PointsFieldPoint     = "point"
	PointsFieldStatus    = "status"
	PointsFieldMemo      = "memo"
	PointsFieldError     = "error"
	PointsFieldCreated   = "created"
	PointsFieldUpdated   = "updated"
)

type Points struct {
	core.BaseRecordProxy
}

func NewPoints(record *core.Record) *Points {
	points := new(Points)
	points.SetProxyRecord(record)
	return points
}

func NewPointsFromCollection(collection *core.Collection) *Points {
	record := core.NewRecord(collection)
	return NewPoints(record)
}

func (points *Points) UserId() string {
	return points.GetString(PointsFieldUserId)
}

func (points *Points) SetUserId(value string) {
	points.Set(PointsFieldUserId, value)
}

func (points *Points) HistoryId() string {
	return points.GetString(PointsFieldHistoryId)
}

func (points *Points) SetHistoryId(value string) {
	points.Set(PointsFieldHistoryId, value)
}

func (points *Points) Point() int {
	return points.GetInt(PointsFieldPoint)
}

func (points *Points) SetPoint(value int) {
	points.Set(PointsFieldPoint, value)
}

func (points *Points) Status() PointStatus {
	return MustParsePointStatus(points.GetString(PointsFieldStatus))
}

func (points *Points) SetStatus(value PointStatus) {
	points.Set(PointsFieldStatus, value.String())
}

func (points *Points) Memo() string {
	return points.GetString(PointsFieldMemo)
}

func (points *Points) SetMemo(value string) {
	points.Set(PointsFieldMemo, value)
}

func (points *Points) Error() string {
	return points.GetString(PointsFieldError)
}

func (points *Points) SetError(value string) {
	points.Set(PointsFieldError, value)
}

func (points *Points) Created() types.DateTime {
	return points.GetDateTime(PointsFieldCreated)
}

func (points *Points) Updated() types.DateTime {
	return points.GetDateTime(PointsFieldUpdated)
}
