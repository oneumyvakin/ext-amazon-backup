<?php
// Copyright 1999-2016. Parallels IP Holdings GmbH.

class Modules_BackupAmazon_LongTasks extends pm_Hook_LongTasks  // Since Plesk 17.0
{
    public function getLongTasks()
    {
        return [
            new Modules_BackupAmazon_Task_UploadAll(),

            new Modules_BackupAmazon_Task_Upload(),
            new Modules_BackupAmazon_Task_ResumeUpload(),
            new Modules_BackupAmazon_Task_CancelUpload(),

            new Modules_BackupAmazon_Task_Download(),
            new Modules_BackupAmazon_Task_ResumeDownload(),
            new Modules_BackupAmazon_Task_CancelDownload(),
        ];
    }
}