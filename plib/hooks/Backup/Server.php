<?php
// Copyright 1999-2016. Parallels IP Holdings GmbH.

class Modules_BackupAmazon_Backup_Server extends pm_Hook_Backup_Server  // Since Plesk 17.0
{
    public function postBackup()
    {
        if (!pm_Settings::get('backupToAmazonEnabled')) {
            pm_Log::debug('Backup on Amazon disabled');
            return;
        }

        if (!pm_Settings::get('uploadNewBackups')) {
            pm_Log::debug('No need to upload new backups');
            return;
        }

        pm_Log::debug('Start task to upload all new backups from ' . __METHOD__);
        $taskManager = new pm_LongTask_Manager();
        $task = new Modules_BackupAmazon_Task_UploadAll();
        $task->setParam('delay', 300); // Wait 5 min to backup process ends
        $taskManager->start($task);
    }
}