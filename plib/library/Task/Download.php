<?php
// Copyright 1999-2016. Parallels IP Holdings GmbH.

class Modules_BackupAmazon_Task_Download extends pm_LongTask_Task // Since Plesk 17.0
{
    const UID = 'download';
    public $trackProgress = false;

    private $result = [];
    private $remotePath;

    public function run()
    {
        $this->remotePath = $this->getParam('remotePath', 'none');
        $this->result = Modules_BackupAmazon_Helper::downloadBackup($this->remotePath);
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
        $localeArgs = ['remotePath' => $this->getParam('remotePath', 'none')];

        switch ($this->getStatus()) {
            case static::STATUS_NOT_STARTED:
                return pm_Locale::lmsg('downloadTaskNotStarted', $localeArgs) . $linkToExtensionHome;
            case static::STATUS_RUNNING:
                return pm_Locale::lmsg('downloadTaskRunning', $localeArgs);
            case static::STATUS_DONE:
                return pm_Locale::lmsg('downloadTaskDone', $localeArgs) . $linkToExtensionHome;
            case static::STATUS_ERROR:
                $localeArgs['error'] = $this->getParam('onError', 'none') . $linkToExtensionHome;
                return pm_Locale::lmsg('downloadTaskError', $localeArgs);
        }
        return pm_Locale::lmsg('unknownTaskStatus', ['status' => $this->getStatus()]);
    }
}