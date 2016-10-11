<?php
// Copyright 1999-2016. Parallels IP Holdings GmbH.

class Modules_BackupAmazon_Task_ResumeDownload extends pm_LongTask_Task // Since Plesk 17.0
{
    const UID = 'resumeDownload';
    public $trackProgress = false;

    private $result = [];
    private $backupName;
    private $remotePath;


    public function run()
    {
        $this->backupName = $this->getParam('backupName', 'none');
        $this->remotePath = $this->getParam('remotePath', 'none');

        $this->result = Modules_BackupAmazon_Helper::resumeDownload($this->remotePath);
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
                return pm_Locale::lmsg('resumeDownloadTaskNotStarted', $localeArgs) . $linkToExtensionHome;
            case static::STATUS_RUNNING:
                return pm_Locale::lmsg('resumeDownloadTaskRunning', $localeArgs);
            case static::STATUS_DONE:
                return pm_Locale::lmsg('resumeDownloadTaskDone', $localeArgs) . $linkToExtensionHome;
            case static::STATUS_ERROR:
                $localeArgs['error'] = $this->getParam('onError', 'none');
                return pm_Locale::lmsg('resumeDownloadTaskError', $localeArgs) . $linkToExtensionHome;
        }
        return pm_Locale::lmsg('unknownTaskStatus', ['status' => $this->getStatus()]);
    }
}