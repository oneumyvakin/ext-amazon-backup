<?php
// Copyright 1999-2016. Parallels IP Holdings GmbH.

$messages = array(
    'tabStorage' => 'Storage',
    'tabSettings' => 'Settings',
    'tabAbout' => 'About',
    'pageTitle' => 'Backup to Amazon S3',

    'settingsAwsAccessKeyIdHelp' => 'Important: <a target="_blank" rel="noopener noreferrer" href="https://console.aws.amazon.com/iam/home?#users">Amazon IAM</a> user should have "AmazonS3FullAccess" permission',
    'settingsTestFailed' => 'Failed to test Amazon S3 connection: %%message%%',
    'enabled' => 'Enabled',
    'awsRegion' => 'Region',
    'awsBucket' => 'Bucket',
    'awsAccessKeyId' => 'Access Key ID',
    'awsSecretAccessKey' => 'Secret Access Key',
    'awsSubFolder' => 'Sub-folder in bucket',
    'uploadNewBackups' => 'Upload new server backups to remote storage',
    'deleteSourceBackup' => 'Delete source Plesk backup after uploading (only incremental)',
    'emailNotificationUploadFailedTitle' => 'Upload backup to Amazon S3 has failed',
    'emailNotificationSubjectUploadNewBackupFailed' => 'Failed to upload new Plesk backup to Amazon S3',
    'emailNotificationBodyUploadNewBackupFailed' => 'Upload new Plesk backup to Amazon S3 failed with error: <error>',

    'us-east-1'      => 'us-east-1',
    'us-west-1'      => 'us-west-1',
    'us-west-2'      => 'us-west-2',
    'eu-west-1'      => 'eu-west-1',
    'eu-central-1'   => 'eu-central-1',
    'ap-southeast-1' => 'ap-southeast-1',
    'ap-southeast-2' => 'ap-southeast-2',
    'ap-northeast-1' => 'ap-northeast-1',
    'sa-east-1'      => 'sa-east-1',

    'SignatureDoesNotMatch' => 'Access Key or Secret Access Key are wrong',
    'InvalidAccessKeyId' => 'Amazon user Access Key is invalid',
    'AccessDenied' => 'Amazon S3 Access denied, check Amazon user\'s permissions and security policy',
    'failedQueryNtpServer' => 'Failed to check local time: %%error%%',
    'gapInTime' => 'Gap in local and real time detected, local time is %%local%%, real time is %%real%%, gap in time is %%diff%%',
    'settingsWasSuccessfullySaved' => 'Settings were successfully saved',

    'date' => 'Date',
    'backupName' => 'Backup Name',
    'isRemote' => 'Remote Storage',
    'isPartialDownload' => 'Partial Download',
    'isLocal' => 'Local Storage',
    'isPartialUpload' => 'Partial Upload',
    'size' => 'Size',
    'mb' => 'Mb',
    'bytes' => 'bytes',
    'actions' => '',
    'backupInvalid' => 'Backup was created with errors',

    'refreshPage' => 'Refresh page',
    'uploadAll' => 'Upload All',
    'uploadAllButtonDescription' => 'Upload all local backups to remote storage',
    'upload' => 'Upload',
    'download' => 'Download',
    'resumeUpload' => 'Resume Upload',
    'resumeDownload' => 'Resume Download',
    'cancelUpload' => 'Cancel Upload',
    'cancelDownload' => 'Cancel Download',

    'uploadAllTaskDone' => 'All new backups were uploaded.',
    'uploadAllTaskError' => 'All new backups were not uploaded due to the following error: %%error%%',
    'uploadAllTaskRunning' => 'Uploading all new backups',
    'uploadAllTaskNotStarted' => 'Upload of all new backups has not started',

    'uploadTaskDone' => 'Backup %%backupName%% was uploaded.',
    'uploadTaskError' => 'Backup %%backupName%% was not uploaded due to the following error: %%error%%',
    'uploadTaskRunning' => 'Backup %%backupName%% is being uploaded',
    'uploadTaskNotStarted' => 'Upload of backup %%backupName%% has not started',

    'resumeUploadTaskDone' => 'Backup %%backupName%% upload has been resumed.',
    'resumeUploadTaskError' => 'Backup %%backupName%% upload has not been resumed due to the following error: %%error%%',
    'resumeUploadTaskRunning' => 'Backup %%backupName%% upload is resumed',
    'resumeUploadTaskNotStarted' => 'Backup %%backupName%% upload has not been resumed',

    'cancelUploadTaskDone' => 'Backup %%backupName%% upload canceled.',
    'cancelUploadTaskError' => 'Backup %%backupName%% upload canceled due to the following error: %%error%%',
    'cancelUploadTaskRunning' => 'Backup %%backupName%% upload canceled',
    'cancelUploadTaskNotStarted' => 'Backup %%backupName%% upload canceled',

    'downloadTaskDone' => 'Backup %%remotePath%% was downloaded.',
    'downloadTaskError' => 'Backup %%remotePath%% was not downloaded due to the following error: %%error%%',
    'downloadTaskRunning' => 'Backup %%remotePath%% is being downloaded',
    'downloadTaskNotStarted' => 'Download of backup %%remotePath%% has not started',

    'resumeDownloadTaskDone' => 'Backup %%remotePath%% download has been resumed.',
    'resumeDownloadTaskError' => 'Backup %%remotePath%% download has not been resumed due to the following error: %%error%%',
    'resumeDownloadTaskRunning' => 'Backup %%remotePath%% download is resumed',
    'resumeDownloadTaskNotStarted' => 'Backup %%remotePath%% download has not been resumed',

    'cancelDownloadTaskDone' => 'Backup %%remotePath%% download canceled.',
    'cancelDownloadTaskError' => 'Backup %%remotePath%% download canceled due to the following error: %%error%%',
    'cancelDownloadTaskRunning' => 'Backup %%remotePath%% download canceled',
    'cancelDownloadTaskNotStarted' => 'Backup %%remotePath%% download canceled',

    'BackupUploadFailedNoDiskSpaceForExport' => 'Free disk space %%freeSpace%% bytes in upload dir %%uploadDir%% is not enough to export backup %%backupName%% of size %%backupSize%% bytes',

    'PleskBackupReturnCodeError' => 'Backup import/export error',

    'PleskBackupReturnCodeImportedExist' => 'Imported backup %%filePath%% already exists in Plesk local storage %%dumpD%%',
    'PleskBackupReturnCodeImportedObjectNotMatch' => 'Imported backup %%filePath%% does not match',
    'PleskBackupReturnCodeImportWrongPassword' => 'Wrong backup %%filePath%% password',
    'PleskBackupReturnCodeImportDeprecatedDumpVersion' => 'Unsupported backup %%filePath%% version',
    'PleskBackupReturnCodeImportWinNativeMailContentSkipped' => 'Native mail content was skipped on import %%filePath%%',
    'PleskBackupReturnCodeImportErrorSign' => 'Backup has a bad signature %%filePath%%',
    'PleskBackupReturnCodeImportNotWellFormedXml' => 'Non-valid XML in backup import %%filePath%%',
    'PleskBackupReturnCodeImportDenied' => 'Backup import denied %%filePath%%',

    'PleskBackupReturnCodeTransportPermissionDenied' => 'Transport Permission Denied',
    'PleskBackupReturnCodeTransportWrongPassword' => 'Transport Password Incorrect',
    'PleskBackupReturnCodeTransportWrongLogin' => 'Transport Username Incorrect',
    'PleskBackupReturnCodeTransportResolveHost' => 'Failed to resolve hostname',
    'PleskBackupReturnCodeTransportUnableConnect' => 'Unable to connect',
    'PleskBackupReturnCodeTransportNetworkError' => 'Transport Network Error',
    'PleskBackupReturnCodeTransportFileNotExist' => 'Transport file does not exist',

    'PleskBackupReturnCodeRepoDumpNotExist' => 'Backup %%backupName%% does not exists',
    'PleskBackupReturnCodeRepoBadDump' => 'Backup imported from %%dstFilePath%% is bad: <a href=\'/admin/backup/restore/?type=local&dumpId=%%backupName%%\'>%%backupName%%</a>',
    'PleskBackupReturnCodeRepoDumpExist' => 'Backup %%backupName%% already exists',
    'PleskBackupReturnCodeRepoPathTooLong' => 'Path too long %%dumpD%%',

    'about' => 'This extension stores your server backups on Amazon S3',
    'feedback' => 'If you have any questions or concerns about this extension, please report your issue on <a rel="noopener noreferrer" target="_blank" href="https://github.com/plesk/ext-backup-amazon">GitHub</a>',
    'faq' => 'FAQ',
    'faqFullNotDeletingQuestion' => 'Q: Why full backups are not removed after upload?',
    'faqFullNotDeletingAnswer' => 'A: Because Plesk can create incremental backups only based on full backup. Incremental backups provide better storage space utilization and reduce storage cost.',

    'faqPartialDownloadPartialUploadQuestion' => 'Q: What do "Partial Download" and "Partial Upload" fields mean?',
    'faqPartialDownloadPartialUploadAnswer' => 'A: They indicate that your backup was not uploaded or downloaded correctly.',

    'faqPartialUploadStoringQuestion' => 'Q: If backup upload to Amazon S3 has failed, where is the uploaded data stored?',
    'faqPartialUploadStoringAnswer' => 'A: Exported backup file is located in DUMP_TMP_D/upload folder, by default it\'s /tmp/upload. Also note that Amazon S3 keeps all uploaded parts (and continutes to charge you) until backup upload is finished or canceled.',

    'faqPartialDownloadStoringQuestion' => 'Q: In backup download from Amazon S3 has failed, where is the downloaded data stored?',
    'faqPartialDownloadStoringAnswer' => 'A: In DUMP_TMP_D/download folder, by default it\'s /tmp/download. You can change DUMP_TMP_D in /etc/psa/psa.conf or in the Windows registry HKLM\SOFTWARE\Wow6432Node\PLESK\PSA Config\Config',

    'faqResumeUploadDownloadQuestion' => 'Q: What can I do in case of a failed upload or download?',
    'faqResumeUploadDownloadAnswer' => 'A: Extension provides the ability to Cancel or Resume failed upload and download operations. The operation continues from the last successfully processed bytes.',

    'faqResumeBackupPlanQuestion' => 'Q: How should I back up my server?',
    'faqResumeBackupPlanAnswer' => 'A: 
        <ol>
          <li>Enable this extension on "Settings" tab</li>
          <li>Create full server backup</li>  
          <li>Upload this full backup to Amazon S3 using this extension from "Storage" tab</li>  
          <li><a href="/admin/backup/schedule" target="_blank">Schedule backup task</a> for creating incremental backups every day.</li>  
          <li>Since this moment all new incremental backup files will be uploaded to Amazon S3. Pay attention that you always have to keep at least one full backup in Plesk backup repository.</li>  
        </ol>',

    'faqResumeRestorePlanQuestion' => 'Q: How should I restore my data from backup files?',
    'faqResumeRestorePlanAnswer' => 'A: Download the latest full backup and latest incremental backup from Amazon S3 using this extension. After downloading these backup files become available in <a target="_blank" href="/admin/backup/list">Plesk Backup Manager</a> and you can restore your data from the latest incremental backup.',

    'faqLogFileLocationQuestion' => 'Q: Where can I find the log file of Amazon S3 operations?',
    'faqLogFileLocationAnswer' => 'A: All actions specific to Amazon S3 are logged to <b>/usr/local/psa/admin/sbin/modules/backup-amazon/backup-on-amazon.log</b>',

);