/**
 * @note
 * open_client.go
 *
 * @author	songtianming
 * @date 	2019-12-04
 */
package cas

import (
	"errors"
	"github.com/tmsong/hlog"
	"net/http"
	"net/url"
	"path"
)

var (
	ErrRespCode = errors.New("resp code is not 200")
)

type OpenClient struct {
	appId   int64
	appKey  string
	client  *http.Client
	logger  *hlog.Logger
	openUrl *url.URL
}

func NewOpenClient(appId int64, appKey string, client *http.Client, openURL *url.URL, l *hlog.Logger) *OpenClient {
	return &OpenClient{
		client:  client,
		openUrl: openURL,
		logger:  l,
		appId:   appId,
		appKey:  appKey,
	}
}

func (c *OpenClient) UserInfoDetailUrl(userId int64, employeeId string) (string, string, error) {
	u, err := c.openUrl.Parse(path.Join(c.openUrl.Path, "api/open/sso/get_user_info_detail"))
	if err != nil {
		return "", "", err
	}
	params := CreateBaseParams(c.appId, c.appKey)
	if len(employeeId) == 0 {
		params["uid"] = userId
	} else {
		params["employeeId"] = employeeId
	}
	return u.String(), JsonEncode(params), nil
}

func (c *OpenClient) DepartmentInfoUrl(departmentId int64) (string, string, error) {
	u, err := c.openUrl.Parse(path.Join(c.openUrl.Path, "api/open/sso/dept/get_dept"))
	if err != nil {
		return "", "", err
	}
	params := CreateBaseParams(c.appId, c.appKey)
	params["id"] = departmentId
	return u.String(), JsonEncode(params), nil
}

func (c *OpenClient) AllDepartmentInfoUrl() (string, string, error) {
	u, err := c.openUrl.Parse(path.Join(c.openUrl.Path, "api/open/sso/dept/get_all_dept"))
	if err != nil {
		return "", "", err
	}
	params := CreateBaseParams(c.appId, c.appKey)
	return u.String(), JsonEncode(params), nil
}

func (c *OpenClient) AllDepartmentUserUrl(departmentId int64, recursive, containsOutgoing bool) (string, string, error) {
	u, err := c.openUrl.Parse(path.Join(c.openUrl.Path, "api/open/sso/get_dept_user_list"))
	if err != nil {
		return "", "", err
	}
	params := CreateBaseParams(c.appId, c.appKey)
	params["deptId"] = departmentId
	params["recursive"] = recursive
	params["containsOutgoing"] = containsOutgoing
	return u.String(), JsonEncode(params), nil
}

func (c *OpenClient) UserPermissionListUrl(userId int64, filterMenu bool) (string, string, error) {
	u, err := c.openUrl.Parse(path.Join(c.openUrl.Path, "api/open/upm/user/permission_list"))
	if err != nil {
		return "", "", err
	}
	params := CreateBaseParams(c.appId, c.appKey)
	params["userId"] = userId
	params["filterMenu"] = filterMenu
	return u.String(), JsonEncode(params), nil
}

func (c *OpenClient) UserInfoVagueUrl(
	account string,
	accountVague bool,
	nameZh string,
	nameZhVague bool,
	nameEn string,
	nameEnVague bool,
	nameDisplay string,
	nameDisplayVague bool,
	email string,
	emailVague bool,
	phone string,
	phoneVague bool,
	employeeId string,
	employeeIdVague bool,
) (string, string, error) {
	u, err := c.openUrl.Parse(path.Join(c.openUrl.Path, "api/open/sso/get_user_vague"))
	if err != nil {
		return "", "", err
	}
	params := CreateBaseParams(c.appId, c.appKey)
	if len(account) > 0 {
		params["account"] = account
		if accountVague {
			params["accountVague"] = accountVague
		}
	}
	if len(nameZh) > 0 {
		params["nameZh"] = nameZh
		if nameZhVague {
			params["nameZhVague"] = nameZhVague
		}
	}
	if len(nameEn) > 0 {
		params["nameEn"] = nameEn
		if nameEnVague {
			params["nameEnVague"] = nameEnVague
		}
	}
	if len(nameDisplay) > 0 {
		params["nameDisplay"] = nameDisplay
		if nameDisplayVague {
			params["nameDisplayVague"] = nameDisplayVague
		}
	}
	if len(email) > 0 {
		params["email"] = email
		if emailVague {
			params["emailVague"] = emailVague
		}
	}
	if len(phone) > 0 {
		params["phone"] = phone
		if phoneVague {
			params["phoneVague"] = phoneVague
		}
	}
	if len(employeeId) > 0 {
		params["employeeId"] = employeeId
		if employeeIdVague {
			params["employeeIdVague"] = employeeIdVague
		}
	}
	return u.String(), JsonEncode(params), nil
}

///////////////////////////////////////////////////////

func (c *OpenClient) UserInfoDetail(userId int64, employeeId string) (*UserInfoDetailResponse, error) {
	u, body, err := c.UserInfoDetailUrl(userId, employeeId)
	if err != nil {
		return nil, err
	}
	ret := PostByJson(u, body, c.logger)
	r := PermissionResponse{}
	err = JsonDecode(ret, &r)
	if err != nil {
		return nil, err
	}
	if r.Code != 200 {
		return nil, ErrRespCode
	}
	re := &UserInfoDetailResponse{}
	err = InterfaceToStruct(r.Data, re)
	if err != nil {
		return nil, err
	}
	return re, nil
}

func (c *OpenClient) DepartmentInfo(departmentId int64) (*DepartmentInfoRespose, error) {
	u, body, err := c.DepartmentInfoUrl(departmentId)
	if err != nil {
		return nil, err
	}
	ret := PostByJson(u, body, c.logger)
	r := PermissionResponse{}
	err = JsonDecode(ret, &r)
	if err != nil {
		return nil, err
	}
	if r.Code != 200 {
		return nil, ErrRespCode
	}
	re := &DepartmentInfoRespose{}
	err = InterfaceToStruct(r.Data, re)
	if err != nil {
		return nil, err
	}
	return re, nil
}

func (c *OpenClient) AllDepartmentInfo() ([]*DepartmentInfoRespose, error) {
	u, body, err := c.AllDepartmentInfoUrl()
	if err != nil {
		return nil, err
	}
	ret := PostByJson(u, body, c.logger)
	r := PermissionResponse{}
	err = JsonDecode(ret, &r)
	if err != nil {
		return nil, err
	}
	if r.Code != 200 {
		return nil, ErrRespCode
	}
	re := []*DepartmentInfoRespose{}
	err = InterfaceToStruct(r.Data, &re)
	if err != nil {
		return nil, err
	}
	return re, nil
}

func (c *OpenClient) AllDepartmentUserInfo(departmentId int64, recursive, containsOutgoing bool) ([]*UserInfoDetailResponse, error) {
	u, body, err := c.AllDepartmentUserUrl(departmentId, recursive, containsOutgoing)
	if err != nil {
		return nil, err
	}
	ret := PostByJson(u, body, c.logger)
	r := PermissionResponse{}
	err = JsonDecode(ret, &r)
	if err != nil {
		return nil, err
	}
	if r.Code != 200 {
		return nil, ErrRespCode
	}
	re := []*UserInfoDetailResponse{}
	err = InterfaceToStruct(r.Data, &re)
	if err != nil {
		return nil, err
	}
	return re, nil
}

func (c *OpenClient) UserPermissionList(userId int64, filterMenu bool) ([]*UserPermissionListResponse, error) {
	u, body, err := c.UserPermissionListUrl(userId, filterMenu)
	if err != nil {
		return nil, err
	}
	ret := PostByJson(u, body, c.logger)
	r := PermissionResponse{}
	err = JsonDecode(ret, &r)
	if r.Code != 200 {
		return nil, err
	}
	re := []*UserPermissionListResponse{}
	InterfaceToStruct(r.Data, &re)
	return re, nil
}

func (c *OpenClient) UserInfoVague(
	account string,
	accountVague bool,
	nameZh string,
	nameZhVague bool,
	nameEn string,
	nameEnVague bool,
	nameDisplay string,
	nameDisplayVague bool,
	email string,
	emailVague bool,
	phone string,
	phoneVague bool,
	employeeId string,
	employeeIdVague bool) ([]*UserInfoVagueResponse, error) {
	u, body, err := c.UserInfoVagueUrl(account,
		accountVague,
		nameZh,
		nameZhVague,
		nameEn,
		nameEnVague,
		nameDisplay,
		nameDisplayVague,
		email,
		emailVague,
		phone,
		phoneVague,
		employeeId,
		employeeIdVague)
	if err != nil {
		return nil, err
	}
	ret := PostByJson(u, body, c.logger)
	r := PermissionResponse{}
	err = JsonDecode(ret, &r)
	if r.Code != 200 {
		return nil, err
	}
	re := []*UserInfoVagueResponse{}
	InterfaceToStruct(r.Data, &re)
	return re, nil
}
