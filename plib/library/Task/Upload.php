<?php
// Copyright 1999-2016. Parallels IP Holdings GmbH.

class Modules_BackupAmazon_Task_Upload extends pm_LongTask_Task // Since Plesk 17.0
{
    const UID = 'upload';
    public $trackProgress = false;
    public $hasDangerousMessage = true;

    private $result = [];
    private $backupName;

    public function run()
    {
        $this->backupName = $this->getParam('backupName', 'none');
        $this->result = Modules_BackupAmazon_Helper::uploadBackup($this->backupName);
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

        $localeArgs = ['backupName' => $this->getParam('backupName', 'none')];

        switch ($this->getStatus()) {
            case static::STATUS_NOT_STARTED:
                return pm_Locale::lmsg('uploadTaskNotStarted', $localeArgs) . $linkToExtensionHome;
            case static::STATUS_RUNNING:
                return pm_Locale::lmsg('uploadTaskRunning', $localeArgs);
            case static::STATUS_DONE:
                return pm_Locale::lmsg('uploadTaskDone', $localeArgs) . $linkToExtensionHome;
            case static::STATUS_ERROR:
                $this->hasDangerousMessage = false;
                $localeArgs['error'] = $this->getParam('onError', 'none');
                return pm_Locale::lmsg('uploadTaskError', $localeArgs) . $linkToExtensionHome;
        }
        return pm_Locale::lmsg('unknownTaskStatus', ['status' => $this->getStatus()]);
    }
}