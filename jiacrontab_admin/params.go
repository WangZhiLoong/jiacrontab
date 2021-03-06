package admin

import (
	"errors"
	"fmt"
	"jiacrontab/models"
	"jiacrontab/pkg/proto"
	"jiacrontab/pkg/util"
	"strings"
)

var (
	paramsError = errors.New("参数错误")
)

type Parameter interface {
	Verify(*myctx) error
}

type JobReqParams struct {
	JobID uint   `json:"jobID" rule:"required,请填写jobID"`
	Addr  string `json:"addr"  rule:"required,请填写addr"`
}

func (p *JobReqParams) Verify(*myctx) error {
	if p.JobID == 0 || p.Addr == "" {
		return paramsError
	}
	return nil
}

type JobsReqParams struct {
	JobIDs []uint `json:"jobIDs" `
	Addr   string `json:"addr"`
}

func (p *JobsReqParams) Verify(ctx *myctx) error {

	if len(p.JobIDs) == 0 || p.Addr == "" {
		return paramsError
	}

	return nil
}

type EditJobReqParams struct {
	JobID            uint              `json:"jobID"`
	Addr             string            `json:"addr" rule:"required,请填写addr"`
	IsSync           bool              `json:"isSync"`
	Name             string            `json:"name" rule:"required,请填写name"`
	Command          []string          `json:"command" rule:"required,请填写name"`
	Code             string            `json:"code"`
	Timeout          int               `json:"timeout"`
	MaxConcurrent    uint              `json:"maxConcurrent"`
	ErrorMailNotify  bool              `json:"errorMailNotify"`
	ErrorAPINotify   bool              `json:"errorAPINotify"`
	ErrorDingdingNotify   bool         `json:"errorDingdingNotify"`
	MailTo           []string          `json:"mailTo"`
	APITo            []string          `json:"APITo"`
	DingdingTo       []string          `json:"DingdingTo"`
	RetryNum         int               `json:"retryNum"`
	WorkDir          string            `json:"workDir"`
	WorkUser         string            `json:"workUser"`
	WorkEnv          []string          `json:"workEnv"`
	WorkIp           []string          `json:"workIp"`
	KillChildProcess bool              `json:"killChildProcess"`
	DependJobs       models.DependJobs `json:"dependJobs"`
	Month            string            `json:"month"`
	Weekday          string            `json:"weekday"`
	Day              string            `json:"day"`
	Hour             string            `json:"hour"`
	Minute           string            `json:"minute"`
	Second           string            `json:"second"`
	TimeoutTrigger   []string          `json:"timeoutTrigger"`
}

func (p *EditJobReqParams) Verify(ctx *myctx) error {
	ts := map[string]bool{
		proto.TimeoutTrigger_CallApi:   true,
		proto.TimeoutTrigger_SendEmail: true,
		proto.TimeoutTrigger_Kill:      true,
		proto.TimeoutTrigger_DingdingWebhook: true,
	}

	for _, v := range p.TimeoutTrigger {
		if !ts[v] {
			return fmt.Errorf("%s:%v", v, paramsError)
		}
	}

	p.Command = util.FilterEmptyEle(p.Command)
	p.MailTo = util.FilterEmptyEle(p.MailTo)
	p.APITo = util.FilterEmptyEle(p.APITo)
	p.DingdingTo = util.FilterEmptyEle(p.DingdingTo)
	p.WorkEnv = util.FilterEmptyEle(p.WorkEnv)
	p.WorkIp = util.FilterEmptyEle(p.WorkIp)

	if p.Month == "" {
		p.Month = "*"
	}

	if p.Weekday == "" {
		p.Weekday = "*"
	}

	if p.Day == "" {
		p.Day = "*"
	}

	if p.Hour == "" {
		p.Hour = "*"
	}

	if p.Minute == "" {
		p.Minute = "*"
	}

	if p.Second == "" {
		p.Second = "*"
	}

	return nil
}

type GetLogReqParams struct {
	Addr     string `json:"addr"`
	JobID    uint   `json:"jobID"`
	Date     string `json:"date"`
	Pattern  string `json:"pattern"`
	IsTail   bool   `json:"isTail"`
	Offset   int64  `json:"offset"`
	Pagesize int    `json:"pagesize"`
}

func (p *GetLogReqParams) Verify(ctx *myctx) error {

	if p.Pagesize <= 0 {
		p.Pagesize = 50
	}

	return nil
}

type DeleteNodeReqParams struct {
	Addr    string `json:"addr" rule:"required,请填写addr"`
	GroupID uint   `json:"groupID"`
}

func (p *DeleteNodeReqParams) Verify(ctx *myctx) error {
	return nil
}

type CleanNodeLogReqParams struct {
	Unit   string `json:"unit" rule:"required,请填写时间单位"`
	Offset int    `json:"offset"`
	Addr   string `json:"addr" rule:"required,请填写addr"`
}

func (p *CleanNodeLogReqParams) Verify(ctx *myctx) error {
	if p.Unit != "day" && p.Unit != "month" {
		return errors.New("不支持的时间单位")
	}
	return nil
}

type SendTestMailReqParams struct {
	MailTo string `json:"mailTo" rule:"required,请填写mailTo"`
}

func (p *SendTestMailReqParams) Verify(ctx *myctx) error {
	return nil
}

type SystemInfoReqParams struct {
	Addr string `json:"addr" rule:"required,请填写addr"`
}

func (p *SystemInfoReqParams) Verify(ctx *myctx) error {
	return nil
}

type GetJobListReqParams struct {
	Addr      string `json:"addr" rule:"required,请填写addr"`
	SearchTxt string `json:"searchTxt"`
	PageReqParams
}

func (p *GetJobListReqParams) Verify(ctx *myctx) error {

	if p.Page <= 1 {
		p.Page = 1
	}

	if p.Pagesize <= 0 {
		p.Pagesize = 50
	}
	return nil
}

type GetGroupListReqParams struct {
	SearchTxt string `json:"searchTxt"`
	PageReqParams
}

func (p *GetGroupListReqParams) Verify(ctx *myctx) error {

	if p.Page <= 1 {
		p.Page = 1
	}

	if p.Pagesize <= 0 {
		p.Pagesize = 50
	}
	return nil
}

type ActionTaskReqParams struct {
	Action string `json:"action" rule:"required,请填写action"`
	Addr   string `json:"addr" rule:"required,请填写addr"`
	JobIDs []uint `json:"jobIDs" rule:"required,请填写jobIDs"`
}

func (p *ActionTaskReqParams) Verify(ctx *myctx) error {
	if len(p.JobIDs) == 0 {
		return paramsError
	}
	return nil
}

type EditDaemonJobReqParams struct {
	Addr            string   `json:"addr" rule:"required,请填写addr"`
	JobID           uint     `json:"jobID"`
	Name            string   `json:"name" rule:"required,请填写name"`
	MailTo          []string `json:"mailTo"`
	APITo           []string `json:"APITo"`
	DingdingTo      []string `json:"DingdingTo"`
	Command         []string `json:"command"  rule:"required,请填写command"`
	Code            string   `json:"code"`
	WorkUser        string   `json:"workUser"`
	WorkIp          []string `json:"workIp"`
	WorkEnv         []string `json:"workEnv"`
	WorkDir         string   `json:"workDir"`
	FailRestart     bool     `json:"failRestart"`
	RetryNum        int      `json:"retryNum"`
	ErrorMailNotify bool     `json:"errorMailNotify"`
	ErrorAPINotify  bool     `json:"errorAPINotify"`
	ErrorDingdingNotify  bool     `json:"errorDingdingNotify"`
}

func (p *EditDaemonJobReqParams) Verify(ctx *myctx) error {
	p.MailTo = util.FilterEmptyEle(p.MailTo)
	p.APITo = util.FilterEmptyEle(p.APITo)
	p.Command = util.FilterEmptyEle(p.Command)
	p.WorkEnv = util.FilterEmptyEle(p.WorkEnv)
	p.WorkIp = util.FilterEmptyEle(p.WorkIp)
	return nil
}

type GetJobReqParams struct {
	JobID uint   `json:"jobID" rule:"required,请填写jobID"`
	Addr  string `json:"addr" rule:"required,请填写addr"`
}

func (p *GetJobReqParams) Verify(ctx *myctx) error {
	return nil
}

type UserReqParams struct {
	Username  string `json:"username" rule:"required,请输入用户名"`
	Passwd    string `json:"passwd,omitempty" rule:"required,请输入密码"`
	GroupID   uint   `json:"groupID"`
	GroupName string `json:"groupName"`
	Avatar    string `json:"avatar"`
	Root      bool   `json:"root"`
	Mail      string `json:"mail"`
}

func (p *UserReqParams) Verify(ctx *myctx) error {
	return nil
}

type InitAppReqParams struct {
	Username string `json:"username" rule:"required,请输入用户名"`
	Passwd   string `json:"passwd" rule:"required,请输入密码"`
	Avatar   string `json:"avatar"`
	Mail     string `json:"mail"`
}

func (p *InitAppReqParams) Verify(ctx *myctx) error {
	return nil
}

type EditUserReqParams struct {
	UserID   uint   `json:"userID" rule:"required,缺少userID"`
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
	OldPwd   string `json:"oldpwd"`
	Avatar   string `json:"avatar"`
	Mail     string `json:"mail"`
}

func (p *EditUserReqParams) Verify(ctx *myctx) error {
	return nil
}

type DeleteUserReqParams struct {
	UserID uint `json:"userID" rule:"required,缺少userID"`
}

func (p *DeleteUserReqParams) Verify(ctx *myctx) error {
	return nil
}

type LoginReqParams struct {
	Username string `json:"username" rule:"required,请输入用户名"`
	Passwd   string `json:"passwd" rule:"required,请输入密码"`
	Remember bool   `json:"remember"`
}

func (p *LoginReqParams) Verify(ctx *myctx) error {
	return nil
}

type PageReqParams struct {
	Page     int `json:"page"`
	Pagesize int `json:"pagesize"`
}

type GetNodeListReqParams struct {
	PageReqParams
	SearchTxt    string `json:"searchTxt"`
	QueryGroupID uint   `json:"queryGroupID"`
	QueryStatus  uint   `json:"queryStatus"`
}

func (p *GetNodeListReqParams) Verify(ctx *myctx) error {

	if p.Page == 0 {
		p.Page = 1
	}
	if p.Pagesize <= 0 {
		p.Pagesize = 50
	}
	return nil
}

type EditGroupReqParams struct {
	GroupID   uint   `json:"groupID" rule:"required,请填写groupID"`
	GroupName string `json:"groupName"  rule:"required,请填写groupName"`
}

func (p *EditGroupReqParams) Verify(ctx *myctx) error {
	return nil
}

type SetGroupReqParams struct {
	TargetGroupID   uint   `json:"targetGroupID"`
	TargetGroupName string `json:"targetGroupName"`
	UserID          uint   `json:"userID" rule:"required,请填写用户ID"`
	Root            bool   `json:"root"`
}

func (p *SetGroupReqParams) Verify(ctx *myctx) error {
	return nil
}

type ReadMoreReqParams struct {
	LastID   int    `json:"lastID"`
	Pagesize int    `json:"pagesize"`
	Keywords string `json:"keywords"`
	Orderby  string `json:"orderby"`
}

func (p *ReadMoreReqParams) Verify(ctx *myctx) error {
	if p.Pagesize == 0 {
		p.Pagesize = 50
	}

	if p.Orderby == "" {
		p.Orderby = "desc"
	}

	p.Keywords = strings.TrimSpace(p.Keywords)
	return nil
}

type GroupNodeReqParams struct {
	Addr            string `json:"addr" rule:"required,请填写addr"`
	TargetNodeName  string `json:"targetNodeName"`
	TargetGroupName string `json:"targetGroupName"`
	TargetGroupID   uint   `json:"targetGroupID"`
}

func (p *GroupNodeReqParams) Verify(ctx *myctx) error {
	return nil
}

type AuditJobReqParams struct {
	JobsReqParams
	JobType string `json:"jobType"`
}

func (p *AuditJobReqParams) Verify(ctx *myctx) error {

	if p.Addr == "" {
		return paramsError
	}

	jobTypeMap := map[string]bool{
		"crontab": true,
		"daemon":  true,
	}

	if err := p.JobsReqParams.Verify(nil); err != nil {
		return err
	}

	if jobTypeMap[p.JobType] == false {
		return paramsError
	}

	return nil
}

type GetUsersParams struct {
	PageReqParams
	SearchTxt    string `json:"searchTxt"`
	IsAll        bool   `json:"isAll"`
	QueryGroupID uint   `json:"queryGroupID"`
}

func (p *GetUsersParams) Verify(ctx *myctx) error {

	if p.Page <= 1 {
		p.Page = 1
	}

	if p.Pagesize <= 0 {
		p.Pagesize = 50
	}
	return nil
}

type CleanLogParams struct {
	IsEvent bool   `json:"isEvent"`
	Unit    string `json:"unit" rule:"required,请填写时间单位"`
	Offset  int    `json:"offset"`
}

func (c *CleanLogParams) Verify(ctx *myctx) error {
	if c.Unit != "day" && c.Unit != "month" {
		return errors.New("不支持的时间单位")
	}
	return nil
}
