<?php
// Copyright 1999-2016. Parallels IP Holdings GmbH.

class Modules_BackupAmazon_Task_UploadAll extends pm_LongTask_Task // Since Plesk 17.0
{
    const UID = 'uploadAll';

    public $hidden = false;
    public $trackProgress = false;
    public $hasDangerousMessage = true;

    private $result = [];

    public function run()
    {
        sleep($this->getParam('delay', 0));

        $this->result = Modules_BackupAmazon_Helper::uploadAllBackups();
        if (isset($this->result['is_error']) && $this->result['is_error']) {
            $msg = $this->result['message'];

            if (isset($this->result['locale_key']) && $this->result['locale_key'] <> '') {
                $msg = pm_Locale::lmsg($this->result['locale_key'], $this->result['locale_args']);
            }

            $notifier = new pm_Notification();
            $notifier->send('uploadFailed', ['error' => $msg], null);

            $this->setParam('onError', $msg);
            throw new pm_Exception();
        }
    }

    public function statusMessage()
    {
        $htdocs = pm_Context::getBaseUrl();
        $localeRefreshPage = pm_Locale::lmsg('refreshPage');
        $linkToExtensionHome = "<div><a href='${htdocs}'>${localeRefreshPage}</a></div>";
        $localeArgs = [];

        switch ($this->getStatus()) {
            case static::STATUS_NOT_STARTED:
                return pm_Locale::lmsg('uploadAllTaskNotStarted', $localeArgs) . $linkToExtensionHome;
            case static::STATUS_RUNNING:
                return pm_Locale::lmsg('uploadAllTaskRunning', $localeArgs);
            case static::STATUS_DONE:
                return pm_Locale::lmsg('uploadAllTaskDone', $localeArgs) . $linkToExtensionHome;
            case static::STATUS_ERROR:
                $this->hasDangerousMessage = false;
                $localeArgs['error'] = $this->getParam('onError', 'none');
                return pm_Locale::lmsg('uploadAllTaskError', $localeArgs) . $linkToExtensionHome;
        }
        return pm_Locale::lmsg('unknownAllTaskStatus', ['status' => $this->getStatus()]);
    }

}