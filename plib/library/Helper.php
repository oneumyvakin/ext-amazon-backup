<?php
// Copyright 1999-2016. Parallels IP Holdings GmbH.

class Modules_BackupAmazon_Helper
{

    /**
     * @param array $args
     * @param array $env
     * @return array|mixed
     */
    public static function execute($args = [], $env = [])
    {
        if (pm_Settings::get('awsSubFolder')) {
            $args[] = '-sub-folder';
            $args[] = pm_Settings::get('awsSubFolder');
        }

        if (pm_Settings::get('deleteSourceBackup')) {
            $args[] = '-delete-after-upload';
        }

        if (!isset($env['AWS_BUCKET'])){
            $env = array_merge($env,[
                'AWS_BUCKET' => pm_Settings::get('awsBucket'),
                'AWS_ACCESS_KEY_ID' => pm_Settings::get('awsAccessKeyId'),
                'AWS_SECRET_ACCESS_KEY' => pm_Settings::get('awsSecretAccessKey'),
            ]);
        }

        $binary = 'backup-amazon.386';

        if (PHP_OS == 'WINNT') {
            $binary = 'backup-amazon.exe';
        }

        pm_Log::debug(print_r($env, true));
        $result = pm_ApiCli::callSbin($binary, $args, pm_ApiCli::RESULT_FULL, $env);
        pm_Log::debug('Execution result : ' . print_r($result, true));

        $response = $result['stdout'];
        if ($result['code'] <> 0) {
            $response = $result['stderr'];
        }

        $jsonResult = json_decode($response, true);
        if (!empty($response) && $jsonResult === null) {
            return [
                'is_error' => true,
                'message' => $result['stdout'] . "\n" . $result['stderr'] . "\n" . 'Exit code: ' . $result['code'],
            ];
        }

        return $jsonResult;
    }

    /** Test settings and access to Amazon S3
     * @param $awsBucket string
     * @param $awsAccessKeyId string
     * @param $awsSecretAccessKey string
     * @param $awsSubFolder string
     * @return array|mixed
     */
    public static function test($awsBucket, $awsAccessKeyId, $awsSecretAccessKey, $awsSubFolder = null)
    {
        $args = [
            '-test'
        ];

        if ($awsSubFolder) {
            $args[] = '-sub-folder';
            $args[] = $awsSubFolder;
        }

        return self::execute(
            $args,
            [
                'AWS_BUCKET' => $awsBucket,
                'AWS_ACCESS_KEY_ID' => $awsAccessKeyId,
                'AWS_SECRET_ACCESS_KEY' => $awsSecretAccessKey
            ]
        );
    }

    /** List Local and Amazon S3 storage
     * @return array {
    "backup_info_1508141631.xml": {
        "IsRemote": false,
        "RemotePath": "",
        "IsLocal": true,
        "Backup": {
            "DumpObject": {
                "Type": "server",
                "Name": "admin",
                "Guid": "013be43d-47af-417f-96e1-ffee2fd52f54"
            },
            "Name": "backup_info_1508141631.xml",
            "Fullname": "backup_info_1508141631.xml",
            "CreationDate": "1508141631",
            "Size": 744789,
            "IsFull": true,
            "Description": "Server configuration and content",
            "OwnerGuid": "013be43d-47af-417f-96e1-ffee2fd52f54",
            "OwnerType": "server",
            "VerificationString": "$AES-128-CBC$GzYSjMy1XHrE3zg/f5MbKw==$mnuGAU7Ny1eIp6bLIa9R1nHUva+Vw+pq5NU0vvPgyeo=",
            "EncryptionType": "panel-key",
            "DumpOriginalVersion": "12.5.30",
            "DumpFormat": "panel",
            "ContentIncluded": true,
            "IncrementBase": 0,
            "IncrementBaseFullname": ""
        }
    },
    "backup_info_1607311943.xml": {
        "IsRemote": true,
        "RemotePath": "backup_info_1607311943.xml.tar",
        "IsLocal": false,
        "Backup": {
            "DumpObject": {
                "Type": "",
                "Name": "",
                "Guid": ""
            },
            "Name": "",
            "Fullname": "",
            "CreationDate": "",
            "Size": 0,
            "IsFull": false,
            "Description": "",
            "OwnerGuid": "",
            "OwnerType": "",
            "VerificationString": "",
            "EncryptionType": "",
            "DumpOriginalVersion": "",
            "DumpFormat": "",
            "ContentIncluded": false,
            "IncrementBase": 0,
            "IncrementBaseFullname": ""
            }
        }
    }
     */
    public static function storage()
    {
        $args = [
            '-list'
        ];

        return self::execute(
            $args
        );
    }

    /** Upload all local backups. Example: -upload-all-backups
     * @return array|mixed
     */
    public static function uploadAllBackups()
    {
        $args = [
            '-upload-all-backups',
        ];

        return self::execute(
            $args
        );
    }

    /** Download and import backup. Example: -download-backup subscription.com/backup_info_1607161246_1607161349.xml.tar.gz
     * @param $remotePath string
     * @return array|mixed
     */
    public static function downloadBackup($remotePath)
    {
        $args = [
            '-download-backup',
            $remotePath
        ];

        return self::execute(
            $args
        );
    }

    /** Exports and uploads backup. Example: -upload-backup backup_info_1606061802.xml
     * @param $backupName string
     * @return array|mixed
     */
    public static function uploadBackup($backupName)
    {
        $args = [
            '-upload-backup',
            $backupName
        ];

        return self::execute(
            $args
        );
    }

    /** Cancel upload. Example: -cancel-upload backup_info_1606061802.xml
     * @param $backupName string
     * @return array|mixed
     */
    public static function cancelUpload($backupName)
    {
        $args = [
            '-cancel-upload',
            $backupName
        ];

        return self::execute(
            $args
        );
    }

    /** Cancel download. Example: -cancel-download backup_info_1606061802.xml
     * @param $remotePath string
     * @return array|mixed
     */
    public static function cancelDownload($remotePath)
    {
        $args = [
            '-cancel-download',
            $remotePath
        ];

        return self::execute(
            $args
        );
    }

    /** Resume backup upload. Example: -resume-upload-backup backup_info_1606061802.xml
     * @param $backupName string
     * @return array|mixed
     */
    public static function resumeUpload($backupName)
    {
        $args = [
            '-resume-upload-backup',
            $backupName
        ];

        return self::execute(
            $args
        );
    }

    /** Resume download and import backup. Example: -resume-download-backup path/file.tar
     * @param $remotePath string
     * @return array|mixed
     */
    public static function resumeDownload($remotePath)
    {
        $args = [
            '-resume-download-backup',
            $remotePath
        ];

        return self::execute(
            $args
        );
    }
}