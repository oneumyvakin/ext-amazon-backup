package plesk

import (
    "encoding/xml"
    "fmt"
    "strconv"
    "time"
    "regexp"
)

const (
    BackupReturnCodeOk int = 0
    BackupReturnCodeError int = 1
    BackupReturnCodeNotSupported int = 2

    BackupReturnCodeImportedExist  int = 111
    BackupReturnCodeImportedObjectNotMatch int  = 112
    BackupReturnCodeImportWrongPassword int  = 113
    BackupReturnCodeImportDeprecatedDumpVersion int  = 114
    BackupReturnCodeImportWinNativeMailContentSkipped int  = 115
    BackupReturnCodeImportErrorSign int  = 116
    BackupReturnCodeImportNotWellFormedXml int  = 117
    BackupReturnCodeImportDenied int  = 118

    BackupReturnCodeTransportPermissionDenied int  = 121
    BackupReturnCodeTransportWrongPassword int  = 122
    BackupReturnCodeTransportWrongLogin  int = 123
    BackupReturnCodeTransportResolveHost int  = 124
    BackupReturnCodeTransportUnableConnect int  = 125
    BackupReturnCodeTransportNetworkError  int = 126
    BackupReturnCodeTransportFileNotExist  int = 127

    BackupReturnCodeRepoDumpNotExist  int = 151
    BackupReturnCodeRepoBadDump  int = 152
    BackupReturnCodeRepoDumpExist  int = 153
    BackupReturnCodeRepoPathTooLong int  = 154
)

var PleskBackupReturnCode = map[int]string{
    BackupReturnCodeOk: "PleskBackupReturnCodeOk",
    BackupReturnCodeError: "PleskBackupReturnCodeError",
    BackupReturnCodeNotSupported: "PleskBackupReturnCodeNotSupported",
    
    BackupReturnCodeImportedExist: "PleskBackupReturnCodeImportedExist",
    BackupReturnCodeImportedObjectNotMatch: "PleskBackupReturnCodeImportedObjectNotMatch",
    BackupReturnCodeImportWrongPassword: "PleskBackupReturnCodeImportWrongPassword",
    BackupReturnCodeImportDeprecatedDumpVersion: "PleskBackupReturnCodeImportDeprecatedDumpVersion",
    BackupReturnCodeImportWinNativeMailContentSkipped: "PleskBackupReturnCodeImportWinNativeMailContentSkipped",
    BackupReturnCodeImportErrorSign: "PleskBackupReturnCodeImportErrorSign",
    BackupReturnCodeImportNotWellFormedXml: "PleskBackupReturnCodeImportNotWellFormedXml",
    BackupReturnCodeImportDenied: "PleskBackupReturnCodeImportDenied",
    
    BackupReturnCodeTransportPermissionDenied: "PleskBackupReturnCodeTransportPermissionDenied",
    BackupReturnCodeTransportWrongPassword: "PleskBackupReturnCodeTransportWrongPassword",
    BackupReturnCodeTransportWrongLogin: "PleskBackupReturnCodeTransportWrongLogin",
    BackupReturnCodeTransportResolveHost: "PleskBackupReturnCodeTransportResolveHost",
    BackupReturnCodeTransportUnableConnect: "PleskBackupReturnCodeTransportUnableConnect",
    BackupReturnCodeTransportNetworkError: "PleskBackupReturnCodeTransportNetworkError",
    BackupReturnCodeTransportFileNotExist: "PleskBackupReturnCodeTransportFileNotExist",
    
    BackupReturnCodeRepoDumpNotExist: "PleskBackupReturnCodeRepoDumpNotExist",
    BackupReturnCodeRepoBadDump: "PleskBackupReturnCodeRepoBadDump",
    BackupReturnCodeRepoDumpExist: "PleskBackupReturnCodeRepoDumpExist",
    BackupReturnCodeRepoPathTooLong: "PleskBackupReturnCodeRepoPathTooLong",
}

type BackupErr struct {
    isError    bool
	code       string
	message    string
	localeKey  string
	localeArgs map[string]string
}

func (e BackupErr) Error() string {
	return e.message
}
func (e BackupErr) IsError() bool {
	return e.isError
}
func (e BackupErr) Code() string {
	return e.code
}
func (e BackupErr) Message() string {
	return e.message
}
func (e BackupErr) LocaleKey() string {
	return e.localeKey
}
func (e BackupErr) LocaleArgs() map[string]string {
	return e.localeArgs
}
func (e BackupErr) Wrap(err error) error {
    e.message = err.Error() + ": " + e.message
    e.localeArgs["err"] = err.Error() + ": " + e.localeArgs["err"]
	return e
}


type DumpList struct {
    DumpList []Dump `xml:"dump"`
}

type Dump struct {
    //Inner tags
    DumpObject DumpObject `xml:"backup-object"`
    DumpStatus DumpStatus `xml:"dump-status"`

    //Attributes
    Name string `xml:"name,attr"`
    Fullname string `xml:"fullname,attr"`
    CreationDate string `xml:"creation-date,attr"`
    Size int64 `xml:"size,attr"`
    IsFull bool `xml:"isFull,attr"`
    Description string `xml:"description,attr"`
    OwnerGuid string `xml:"owner-guid,attr"`
    OwnerType string `xml:"owner-type,attr"`
    VerificationString string `xml:"verification-string,attr"`
    EncryptionType string `xml:"encryption-type,attr"`
    DumpOriginalVersion string `xml:"dump-original-version,attr"`
    DumpFormat string `xml:"dump-format,attr"`
    ContentIncluded bool `xml:"content-included,attr"`
    IncrementBase int64 `xml:"increment-base,attr"`
    IncrementBaseFullname string `xml:"increment-base-fullname,attr"`
}

type DumpObject struct {
     //Attributes
     Type string `xml:"type,attr"` // "domain" or "server"
     Name string `xml:"name,attr"` // subscription's name or "admin"
     Guid string `xml:"guid,attr"` // owner guid
}

type DumpStatus struct {
     //Attributes
     DumpStatus string `xml:"dump-status,attr"` // "OK"
     BackupProcessStatus string `xml:"backup-process-status,attr"` // "ERROR" or "WARNINGS" or "SUCCESS"
}

func (self Plesk) GetBackupListFromLocalStorage() ([]Dump, error) {
    /*
    <?xml version="1.0" encoding="UTF-8"?>
    <dump-list>
     <dump name="backup_info_1609042052_1609101705.xml" fullname="backup_info_1609042052_1609101705.xml" creation-date="1609101705" size="222300" isFull="true" description="All configuration and content" owner-guid="64110f72-f73f-4c2a-bd4c-c9ffaf40b75f" owner-type="server" verification-string="" encryption-type="" dump-version="17.0.15" dump-original-version="17.0.15" dump-format="panel" content-included="true" increment-base="1609042052" increment-base-fullname="backup_info_1609042052.xml">
      <backup-object type="server" name="admin" guid="64110f72-f73f-4c2a-bd4c-c9ffaf40b75f"/>
      <dump-status dump-status="OK" backup-process-status="ERROR">
       <details>
        <message>backup_info_1609042052_1609101705.xml: </message>
       </details>
      </dump-status>
      <related-dumps>
       <related-dump>1609042052</related-dump>
       <related-dump>1609050158</related-dump>
      </related-dumps>
      <dump-result></dump-result>
     </dump>
      <dump name="backup_info_1609042052_1609050158.xml" fullname="backup_info_1609042052_1609050158.xml" creation-date="1609050158" size="18101729" isFull="true" description="All configuration and content" owner-guid="64110f72-f73f-4c2a-bd4c-c9ffaf40b75f" owner-type="server" verification-string="$AES-128-CBC$ZywWvHjz1MZpIaaHwYM2Aw==$C6gbYwAtRp0Otl1JJEmMbN2k1dAzYa8Ts+ulqFIpy3Q=" encryption-type="panel-key" dump-version="17.0.15" dump-original-version="17.0.15" dump-format="panel" content-included="true" increment-base="1609042052" increment-base-fullname="backup_info_1609042052.xml">
      <backup-object type="server" name="admin" guid="64110f72-f73f-4c2a-bd4c-c9ffaf40b75f"/>
      <dump-status dump-status="OK" backup-process-status="WARNINGS">
       <details>
        <message>backup_info_1609042052_1609050158.xml: </message>
       </details>
      </dump-status>
      <related-dumps>
       <related-dump>1609042052</related-dump>
      </related-dumps>
      <dump-result></dump-result>
     </dump>
     <dump name="backup_info_1605120057.xml" fullname="backup_info_1605120057.xml" creation-date="1605120057" size="50997425" isFull="true" description="п²п╟я│я┌я─п╬п╧п╨п╦ п╦ п╨п╬пҐя┌п╣пҐя┌ я│п╣я─п╡п╣я─п╟" owner-guid="df70004e-c33a-4b9d-8317-0748c2256303" owner-type="server" verification-string="$AES-128-CBC$r2S+Eyk0B6nhBPOfkGdRPw==$pFXWFU1UC+cOqjy788Fyu8ndOS1iNlrZEDjSTp2nl5U=" encryption-type="panel-key" dump-original-version="12.5.30" dump-format="panel" content-included="true">
      <backup-object type="server" name="admin" guid="df70004e-c33a-4b9d-8317-0748c2256303"/>
      <dump-status dump-status="OK" backup-process-status="SUCCESS">
       <details>
        <message>backup_info_1605120057.xml: </message>
       </details>
      </dump-status>
      <dump-result></dump-result>
     </dump>
     <dump name="backup_info_1605120057_1607160620.xml" fullname="backup_info_1605120057_1607160620.xml" creation-date="1607160620" size="51398931" isFull="true" description="Server configuration and content" owner-guid="df70004e-c33a-4b9d-8317-0748c2256303" owner-type="server" verification-string="$AES-128-CBC$F0sNnj9+C2pOPvOyZX+Qzg==$GbKOR31DRgP5vgvfOL/w/uTI6iEscNrq4LTsFUnR1xw=" encryption-type="panel-key" dump-original-version="12.5.30" dump-format="panel" content-included="true" increment-base="1605120057" increment-base-fullname="backup_info_1605120057.xml">
      <backup-object type="server" name="admin" guid="df70004e-c33a-4b9d-8317-0748c2256303"/>
      <dump-status dump-status="OK" backup-process-status="SUCCESS">
       <details>
        <message>backup_info_1605120057_1607160620.xml: </message>
       </details>
      </dump-status>
      <dump-result></dump-result>
     </dump>
     <dump name="backup_info_1607160616.xml" fullname="clients/cl1/domains/cl1-110.tld/backup_info_1607160616.xml" creation-date="1607160616" size="115613" isFull="true" description="All configuration and content" owner-guid="eeca455f-31e0-4351-a720-edebf3ee7f6e" owner-type="client" verification-string="$AES-128-CBC$ezivrE8dkAJqQ6VyXbpwfQ==$Eni/IqkAuUj+JsLkDCi9IST9Y7t0PMLbWqLZ9pK/G1w=" encryption-type="panel-key" dump-original-version="12.5.30" dump-format="panel" content-included="true">
      <backup-object type="domain" name="cl1-110.tld" guid="085c6a65-5f05-4b05-a339-ea89b5eb3419"/>
      <dump-status dump-status="OK" backup-process-status="SUCCESS">
       <details>
        <message>clients/cl1/domains/cl1-110.tld/backup_info_1607160616.xml: </message>
       </details>
      </dump-status>
      <dump-result></dump-result>
     </dump>
     <dump name="backup_info_1607160616_1607160618.xml" fullname="clients/cl1/domains/cl1-110.tld/backup_info_1607160616_1607160618.xml" creation-date="1607160618" size="21449" isFull="true" description="All configuration and content" owner-guid="eeca455f-31e0-4351-a720-edebf3ee7f6e" owner-type="client" verification-string="$AES-128-CBC$/5ctMiJiGPEnWb5wnz4KWQ==$3HUTNjInqCQZQy10PGXGLmNjhwHl3CAHFwKxzh7MEFk=" encryption-type="panel-key" dump-original-version="12.5.30" dump-format="panel" content-included="true" increment-base="1607160616" increment-base-fullname="clients/cl1/domains/cl1-110.tld/backup_info_1607160616.xml">
      <backup-object type="domain" name="cl1-110.tld" guid="085c6a65-5f05-4b05-a339-ea89b5eb3419"/>
      <dump-status dump-status="OK" backup-process-status="SUCCESS">
       <details>
        <message>clients/cl1/domains/cl1-110.tld/backup_info_1607160616_1607160618.xml: </message>
       </details>
      </dump-status>
      <related-dumps>
       <related-dump>1607160616</related-dump>
      </related-dumps>
      <dump-result></dump-result>
     </dump>
     <dump name="backup_info_1601190210.xml" fullname="domains/deb8x64-plesk12-5.local/backup_info_1601190210.xml" creation-date="1601190210" size="16999254" isFull="true" description="п▓я│п╣ пҐп╟я│я┌я─п╬п╧п╨п╦ п╦ п╨п╬пҐя┌п╣пҐя┌" owner-guid="df70004e-c33a-4b9d-8317-0748c2256303" owner-type="server" verification-string="$AES-128-CBC$tLOwG2a9KGIZM1ezMQjfZg==$nqFes8kbJopZZioeMpX1Clxhx2rm6pX1Ry9s3H0huOE=" encryption-type="panel-key" dump-original-version="12.5.30" dump-format="panel" content-included="true">
      <backup-object type="domain" name="deb8x64-plesk12-5.local" guid="1dff611d-8a49-4c67-bfff-ae4e0d490fd5"/>
      <dump-status dump-status="OK" backup-process-status="SUCCESS">
       <details>
        <message>domains/deb8x64-plesk12-5.local/backup_info_1601190210.xml: </message>
       </details>
      </dump-status>
      <dump-result></dump-result>
     </dump>
     <dump name="backup_info_1607160958.xml" fullname="domains/lo.xsstest.ru/backup_info_1607160958.xml" creation-date="1607160958" size="134826" isFull="true" description="backup lo.xsstest.ru" owner-guid="df70004e-c33a-4b9d-8317-0748c2256303" owner-type="server" verification-string="$AES-128-CBC$UqJzK3v0WdHmhWkD0IuqtA==$wwn8SE80JIVtkYNdYV2E46tWCCpMnhx/P6fY1KPmZWY=" encryption-type="panel-key" dump-original-version="12.5.30" dump-format="panel" content-included="true">
      <backup-object type="domain" name="lo.xsstest.ru" guid="485e32b6-b31e-49aa-abab-c501e20611a6"/>
      <dump-status dump-status="OK" backup-process-status="SUCCESS">
       <details>
        <message>domains/lo.xsstest.ru/backup_info_1607160958.xml: </message>
       </details>
      </dump-status>
      <dump-result></dump-result>
     </dump>
    </dump-list>
    */

    _, output, _, err := execute(self.Log, self.Config["pmm-ras"], "--get-dump-list", "--dump-storage=" + self.Config["DUMP_D"])
    if err != nil {
        return nil, fmt.Errorf("Failed to get backup list: Failed to execute pmm-ras: %s\n", err)
    }
    var dumpList DumpList
    err = xml.Unmarshal(output, &dumpList)
    if err != nil {
        return nil, fmt.Errorf("Failed to get backup list: Failed to parse xml: %s\n", err)
    }
    
    for i, _ := range dumpList.DumpList {
        rfc3339, err := self.BackupDateToRFC3339(dumpList.DumpList[i].CreationDate)
        if err != nil {
            return nil, fmt.Errorf("Failed to parse date: %s", err)
        }
        dumpList.DumpList[i].CreationDate = rfc3339
    }

    self.Log.Printf("Successfully get backup list %#v\n", dumpList.DumpList)
    return dumpList.DumpList, err
}

// Export backup as tar file
func (self Plesk) ExportBackupFromLocalStorage(name, dstPath string, includeIncrements bool) (error) {
    args := []string{
        "--export-dump-as-file",
        "--dump-specification=" + name,
        "--dump-file-specification=" + dstPath,
    }
    if includeIncrements {
        args = append(args, "--include-increments")
    }
    
    _, _, exitCode, err := execute(
        self.Log,
        self.Config["pmm-ras"],
        args...
    )
    
    if err != nil {
        return BackupErr{
            isError:    true,
            code:       self.GetBackupErrorCode(exitCode),
            message:    fmt.Sprintf("Failed to export backup %s to %s: %s\n", name, dstPath, err),
            localeKey:  PleskBackupReturnCode[exitCode],
            localeArgs: map[string]string {
                "backupName": name,
                "dstPath": dstPath,
            },
        }
    }

    self.Log.Println("Successfully export backup %s to %s", name, dstPath)
    return nil
}

// Import backup from tar file
func (self Plesk) ImportBackupToLocalStorage(filePath, backupPassword string, checkSign bool) (error) {
    args := []string{
        "--import-file-as-dump",
        "--dump-file-specification=" + filePath,
        "--dump-storage=" + self.Config["DUMP_D"],
        "--force",
        // --type=server --session-path=/var/log/plesk/PMM
    }
    if checkSign {
        args = append(args, "--check-sign")
    }

    _, _, exitCode, err := execute(
        self.Log,
        self.Config["pmm-ras"],
        args...
    )
    
    if err != nil {
        returnErr := BackupErr{
            isError:    true,
            code:       self.GetBackupErrorCode(exitCode),
            message:    fmt.Sprintf("Failed to import backup file %s to dump storage %s with error: %s\n", filePath, self.Config["DUMP_D"], err),
            localeKey:  PleskBackupReturnCode[exitCode],
            localeArgs: map[string]string {
                "filePath": filePath,
                "dumpD": self.Config["DUMP_D"],
                "checkSign": strconv.FormatBool(checkSign),
            },
        }
        switch exitCode {
        case BackupReturnCodeImportErrorSign:
            if !checkSign {
                break
            }
        default:
            return returnErr
        }
    }

    self.Log.Println("Successfully import backup file %s", filePath)
    return nil
}

// Delete backup in local storage
func (self Plesk) DeleteBackupFromLocalStorage(name string) (error) {
    _, _, exitCode, err := execute(
        self.Log,
        self.Config["pmm-ras"],
        "--delete-dump",
        "--dump-specification=" + name,
        "--dump-storage=" + self.Config["DUMP_D"],
    )
    if err != nil {
        return BackupErr{
            isError:    true,
            code:       self.GetBackupErrorCode(exitCode),
            message:    fmt.Sprintf("Failed to delete backup %s: %s\n", name, err),
            localeKey:  PleskBackupReturnCode[exitCode],
            localeArgs: map[string]string {
                "backupName": name,
                "dumpD": self.Config["DUMP_D"],
                "err": fmt.Sprintf("Failed to delete backup %s: %s\n", name, err),
            },
        }
    }

    self.Log.Println("Successfully delete backup %s", name)
    return nil
}

// Format date from 1607160958 to 2016-07-16 09:58:00 +0000 UTC
func (self Plesk) BackupDateToRFC3339(date string) (string, error) {
    datePartsStrings := regexp.MustCompile("\\d{2}").FindAllString (date, -1)
    if len(datePartsStrings) != 5 {
        return date, fmt.Errorf("Failed to convert date %s to RFC3339. Format must be like 1607160958", date)
    }
    dateParts := make([]int, 5)
    for i, part := range datePartsStrings {
        partInt, err := strconv.Atoi(part)
        if err != nil {
            return date, err
        }
         
        dateParts[i] = partInt
    }

    t := time.Date(2000 + dateParts[0], time.Month(dateParts[1]), dateParts[2], dateParts[3], dateParts[4], 0, 0, time.Local)
    return t.String(), nil
}

func (self Plesk) GetBackupErrorCode(exitCode int) (errCode string) {
    if code, ok := PleskBackupReturnCode[exitCode]; ok {
        errCode = code
    } else {
        errCode = "PleskBackupReturnCodeUnknown"
    }
    return
}