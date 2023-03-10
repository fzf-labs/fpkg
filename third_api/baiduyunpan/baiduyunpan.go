package baiduyunpan

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"net/url"
	"path"
	"strconv"
)

type File struct {
	Path    string `json:"path"`
	Newname string `json:"newname"`
	Dest    string `json:"dest"`
	Ondup   string `json:"ondup"`
}

type FileMoveResp struct {
	Errno int `json:"errno"`
	Info  []struct {
		Errno int    `json:"errno"`
		Path  string `json:"path"`
	} `json:"info"`
	RequestId int64 `json:"request_id"`
}

func FileMove(accessToken string, fileList []File) error {
	fileMoveurl := "http://pan.baidu.com/rest/2.0/xpan/file"
	queryParams := map[string]string{
		"method":       "filemanager",
		"access_token": accessToken,
		"opera":        "move",
	}
	formParams := url.Values{}
	formParams.Add("async", "1")
	fileListJson, err := json.Marshal(fileList)
	if err != nil {
		return err
	}
	formParams.Add("filelist", string(fileListJson))
	var response FileMoveResp
	client := req.C()
	resp, err := client.R().SetQueryParams(queryParams).SetFormDataFromValues(formParams).SetSuccessResult(&response).Post(fileMoveurl)
	if err != nil {
		return err
	}
	if !resp.IsSuccessState() {
		return fmt.Errorf("bad response status: %s", resp.Status)
	}
	if response.Errno != 0 {
		return fmt.Errorf("response errno: %d", response.Errno)
	}
	return nil
}

type FileListAllReq struct {
}
type FileListAllResp struct {
	Cursor  int    `json:"cursor"`
	Errmsg  string `json:"errmsg"`
	Errno   int    `json:"errno"`
	HasMore int    `json:"has_more"`
	List    []struct {
		Category       int    `json:"category"`
		FsId           int64  `json:"fs_id"`
		Isdir          int    `json:"isdir"`
		LocalCtime     int    `json:"local_ctime"`
		LocalMtime     int    `json:"local_mtime"`
		Md5            string `json:"md5"`
		Path           string `json:"path"`
		ServerCtime    int    `json:"server_ctime"`
		ServerFilename string `json:"server_filename"`
		ServerMtime    int    `json:"server_mtime"`
		Size           int    `json:"size"`
		Thumbs         struct {
			Url1 string `json:"url1"`
			Url2 string `json:"url2"`
			Url3 string `json:"url3"`
			Icon string `json:"icon,omitempty"`
		} `json:"thumbs"`
	} `json:"list"`
	RequestId string `json:"request_id"`
}
type FileListResp struct {
	Errno    int    `json:"errno"`
	GuidInfo string `json:"guid_info"`
	List     []struct {
		ServerFilename string `json:"server_filename"`
		Privacy        int    `json:"privacy"`
		Category       int    `json:"category"`
		Unlist         int    `json:"unlist"`
		FsId           int64  `json:"fs_id"`
		DirEmpty       int    `json:"dir_empty"`
		ServerAtime    int    `json:"server_atime"`
		ServerCtime    int    `json:"server_ctime"`
		LocalMtime     int    `json:"local_mtime"`
		Size           int    `json:"size"`
		Isdir          int    `json:"isdir"`
		Share          int    `json:"share"`
		Path           string `json:"path"`
		LocalCtime     int    `json:"local_ctime"`
		ServerMtime    int    `json:"server_mtime"`
		Empty          int    `json:"empty"`
		OperId         int    `json:"oper_id"`
	} `json:"list"`
}

func FileList(accessToken string, path string) (*FileListResp, error) {
	url := "http://pan.baidu.com/rest/2.0/xpan/file"
	queryParams := map[string]string{
		"method":       "list",
		"access_token": accessToken,
		"dir":          path,
	}
	var response FileListResp
	client := req.C()
	resp, err := client.R().SetQueryParams(queryParams).SetSuccessResult(&response).Post(url)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("bad response status: %s", resp.Status)
	}
	return &response, nil
}

func FileListMultimedia(accessToken string, path string, start int, limit int) (*FileListAllResp, error) {
	url := "http://pan.baidu.com/rest/2.0/xpan/multimedia"
	queryParams := map[string]string{
		"method":       "listall",
		"access_token": accessToken,
		"path":         path,
		"recursion":    "1",
		"start":        strconv.Itoa(start),
		"limit":        strconv.Itoa(limit),
	}
	var response FileListAllResp
	client := req.C()
	resp, err := client.R().SetQueryParams(queryParams).SetSuccessResult(&response).Post(url)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("bad response status: %s", resp.Status)
	}
	if response.Errno != 0 {
		return nil, fmt.Errorf("response errno: %d", response.Errno)
	}
	return &response, nil
}

func FileListAll(accessToken string, path string) ([]string, error) {
	start := 0
	limit := 1000
	files := make([]string, 0)
	for {
		listAll, err := FileListMultimedia(accessToken, path, start, limit)
		if err != nil {
			return nil, err
		}
		for _, v := range listAll.List {
			if v.Isdir == 0 {
				files = append(files, v.Path)
			}
		}
		if listAll.HasMore == 0 {
			break
		}
		start = listAll.Cursor
	}
	return files, nil
}

// srcPath := "/_pcs_.appdata/xpan/zhuanshu/276828528969728/编程教程"
// dstPath := "/专属空间/编程教程"
func Move(accessToken string, srcPath string, dstPath string) error {
	fileMoveUrl := "http://pan.baidu.com/rest/2.0/xpan/file"
	dst := path.Dir(dstPath)
	newName := path.Base(dstPath)
	fileList := make([]File, 0)
	fileList = append(fileList, File{
		Path:    srcPath,
		Newname: newName,
		Dest:    dst,
		Ondup:   "overwrite",
	})
	queryParams := map[string]string{
		"method":       "filemanager",
		"access_token": accessToken,
		"opera":        "move",
	}
	formParams := url.Values{}
	formParams.Add("async", "1")
	fileListJson, err := json.Marshal(fileList)
	if err != nil {
		return err
	}
	formParams.Add("filelist", string(fileListJson))
	var response FileMoveResp
	client := req.C()
	resp, err := client.R().SetQueryParams(queryParams).SetFormDataFromValues(formParams).SetSuccessResult(&response).Post(fileMoveUrl)
	if err != nil {
		return err
	}
	if !resp.IsSuccessState() {
		return fmt.Errorf("bad response status: %s", resp.Status)
	}
	if response.Errno != 0 {
		fmt.Println(response)
		return fmt.Errorf("response errno: %d", response.Errno)
	}
	return nil
}

type FileDirResp struct {
	Path       string `json:"path"`
	Uploadid   string `json:"uploadid"`
	ReturnType int    `json:"return_type"`
	BlockList  []int  `json:"block_list"`
	Errno      int    `json:"errno"`
	RequestId  int64  `json:"request_id"`
}

func Dir(accessToken string, path string) error {
	dirUrl := "http://pan.baidu.com/rest/2.0/xpan/file"
	queryParams := map[string]string{
		"method":       "create",
		"access_token": accessToken,
	}
	formParams := url.Values{}
	formParams.Add("path", path)
	formParams.Add("size", "0")
	formParams.Add("isdir", "1")
	formParams.Add("block_list", "[]")
	formParams.Add("uploadid", "")
	var response FileDirResp
	client := req.C()
	resp, err := client.R().SetQueryParams(queryParams).SetFormDataFromValues(formParams).SetSuccessResult(&response).Post(dirUrl)
	if err != nil {
		return err
	}
	if !resp.IsSuccessState() {
		return fmt.Errorf("bad response status: %s", resp.Status)
	}
	if response.Errno != 0 {
		return fmt.Errorf("response errno: %d", response.Errno)
	}
	return nil
}
