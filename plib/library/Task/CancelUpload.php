<?php
// Copyright 1999-2016. Parallels IP Holdings GmbH.

class Modules_BackupAmazon_Task_CancelUpload extends pm_LongTask_Task // Since Plesk 17.0
{
    const UID = 'cancelUpload';
    public $trackProgress = false;

    private $result = [];
    private $backupName;
    private $remotePath;
    private $uploadId;

    public function run()
    {
        $this->backupName = $this->getParam('backupName', 'none');
        $this->remotePath = $this->getParam('remotePath', 'none');
        $this->uploadId = $this->getParam('uploadId', 'none');
        $this->result = Modules_BackupAmazon_Helper::cancelUpload($this->backupName);
        if (isset($this->result['is_error']) && $this->result['is_error']) {
            $msg = $this->result['message'];

            if (isset($this->result['locale_key']) && $this->result['locale_key'] <> '') {
                $msg = pm_Locale::lmsg($this->result['locale_key'], $this->result['locale_args']);
            }

            $this->setParam('onError', $msg);
            throw new pm_Exception();
        }
    }

    public function statusMessage()
    {
        $htdocs = pm_Context::getBaseUrl();
        $localeRefreshPage = pm_Locale::lmsg('refreshPage');
        $linkToExtensionHome = "<div><a href='${htdocs}'>${localeRefreshPage}</a></div>";
        $localeArgs = [
            'backupName' => $this->getParam('backupName', 'none'),
            'remotePath' => $this->getParam('remotePath', 'none'),
            'uploadId' => $this->getParam('uploadId', 'none'),
        ];

        switch ($this->getStatus()) {
            case static::STATUS_NOT_STARTED:
                return pm_Locale::lmsg('cancelUploadTaskNotStarted', $localeArgs) . $linkToExtensionHome;
            case static::STATUS_RUNNING:
                $localeArgs['taskId'] = $this->getId();
                return pm_Locale::lmsg('cancelUploadTaskRunning', $localeArgs);
            case static::STATUS_DONE:
                return pm_Locale::lmsg('cancelUploadTaskDone', $localeArgs) . $linkToExtensionHome;
            case static::STATUS_ERROR:
                $localeArgs['error'] = $this->getParam('onError', 'none');
                return pm_Locale::lmsg('cancelUploadTaskError', $localeArgs) . $linkToExtensionHome;
            case static::STATUS_CANCELED:
                return pm_Locale::lmsg('cancelUploadTaskCanceled', $localeArgs) . $linkToExtensionHome;
        }
        return pm_Locale::lmsg('unknownTaskStatus', ['status' => $this->getStatus()]);
    }

}