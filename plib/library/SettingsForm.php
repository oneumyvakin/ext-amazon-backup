<?php
// Copyright 1999-2016. Parallels IP Holdings GmbH.

class Modules_BackupAmazon_SettingsForm extends pm_Form_Simple
{
    /**
     * Validate the form
     *
     * @param  array $data
     * @return boolean
     */
    function isValid($data)
    {
        parent::isValid($data);

        $backupToAmazonEnabled = $this->getElement('backupToAmazonEnabled')->getValue();
        $awsBucket = $this->getElement('awsBucket')->getValue();
        $awsAccessKeyId = $this->getElement('awsAccessKeyId')->getValue();
        $awsSecretAccessKey = $this->getElement('awsSecretAccessKey')->getValue();
        $awsSubFolder = $this->getElement('awsSubFolder')->getValue();
        $deleteSourceBackup = $this->getElement('deleteSourceBackup')->getValue();

        if ($backupToAmazonEnabled) {
            $result = Modules_BackupAmazon_Helper::test($awsBucket, $awsAccessKeyId, $awsSecretAccessKey, $awsSubFolder);
            if (isset($result['is_error'])) {

                $msg = $result['message'];

                if ($result['locale_key'] <> '') {
                    $msg = pm_Locale::lmsg($result['locale_key'], $result['locale_args']);
                }

                $this->getElement('backupToAmazonEnabled')->addError(pm_Locale::lmsg('settingsTestFailed', ['message' => $msg]));
                $this->markAsError();
                return false;
            }
        }

        return true;
    }
}