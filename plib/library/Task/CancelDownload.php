<?php
// Copyright 1999-2016. Parallels IP Holdings GmbH.

class Modules_BackupAmazon_Task_CancelDownload extends pm_LongTask_Task // Since Plesk 17.0
{
    const UID = 'cancelDownload';
    public $trackProgress = false;

    private $result = [];
    private $backupName;
    private $remotePath;

    public function run()
    {
        $this->backupName = $this->getParam('backupName', 'none');
        $this->remotePath = $this->getParam('remotePath', 'none');

        $this->result = Modules_BackupAmazon_Helper::cancelDownload($this->backupName);
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
        ];

        switch ($this->getStatus()) {
            case static::STATUS_NOT_STARTED:
                return pm_Locale::lmsg('cancelDownloadTaskNotStarted', $localeArgs) . $linkToExtensionHome;
            case static::STATUS_RUNNING:
                $localeArgs['taskId'] = $this->getId();
                return pm_Locale::lmsg('cancelDownloadTaskRunning', $localeArgs);
            case static::STATUS_DONE:
                return pm_Locale::lmsg('cancelDownloadTaskDone', $localeArgs) . $linkToExtensionHome;
            case static::STATUS_ERROR:
                $localeArgs['error'] = $this->getParam('onError', 'none') . $linkToExtensionHome;
                return pm_Locale::lmsg('cancelDownloadTaskError', $localeArgs);
            case static::STATUS_CANCELED:
                return pm_Locale::lmsg('cancelDownloadTaskCanceled', $localeArgs) . $linkToExtensionHome;
        }
        return pm_Locale::lmsg('unknownTaskStatus', ['status' => $this->getStatus()]);
    }

}