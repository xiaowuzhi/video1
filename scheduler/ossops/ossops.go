package ossops

import (
    "github.com/qiniu/api.v7/auth/qbox"
    "github.com/qiniu/api.v7/storage"
    "golang.org/x/net/context"
    "log"
    "time"
    "avenssi/config"
)

var AK string
var SK string
var DOMAIN string
var EP string

func init() {
    AK = "1XiP12dWrUkeTzJAKW-tdQ7iBbxHXtwcQyIssV34"
    SK = "zpxjHfl5GeOoZhCliMEEsia4mjUeXqImZlYfIghK"
    DOMAIN = "http://pewiokxkt.bkt.clouddn.com"
    EP = config.GetOssAddr()
}

func UploadToOss(filename string, path string, bn string) bool {

    putPolicy := storage.PutPolicy{
        Scope: bn,
    }
    mac := qbox.NewMac(AK, SK)
    upToken := putPolicy.UploadToken(mac)
    cfg := storage.Config{}
    // 空间对应的机房
    cfg.Zone = &storage.ZoneHuabei
    // 是否使用https域名
    cfg.UseHTTPS = false
    // 上传是否使用CDN上传加速
    cfg.UseCdnDomains = false
    // 构建表单上传的对象
    formUploader := storage.NewFormUploader(&cfg)
    ret := storage.PutRet{}
    // 可选配置
    putExtra := storage.PutExtra{
        Params: map[string]string{
            "x:name": filename,
        },
    }
    err := formUploader.PutFile(context.Background(), &ret, upToken, filename, path, &putExtra)
    if err != nil {
        log.Printf("Uploading object error: %s", err)
        return false
    }
    return true
}

func XiaoPrivateAccessURL(key string, deadline int64) (privateURL string) {
    mac := qbox.NewMac(AK, SK)
    if deadline == 0 {
        deadline = time.Now().Add(time.Second * 3600).Unix() //1小时有效期
    }
    privateURL = storage.MakePrivateURL(mac, DOMAIN, key, deadline)
    return
}

func DeleteObject(filename string, bn string) bool {

    mac := qbox.NewMac(AK, SK)
    cfg := storage.Config{
        // 是否使用https域名进行资源管理
        UseHTTPS: false,
    }
    // 指定空间所在的区域，如果不指定将自动探测
    // 如果没有特殊需求，默认不需要指定
    //cfg.Zone=&storage.ZoneHuabei
    bucketManager := storage.NewBucketManager(mac, &cfg)
    err := bucketManager.Delete(bn, filename)
    if err != nil {
        log.Printf("Deleting object error: %s", err)
        return false
    }

    return true

}
