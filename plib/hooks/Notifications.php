<?php
// Copyright 1999-2016. Parallels IP Holdings GmbH.

class Modules_BackupAmazon_Notifications extends pm_Hook_Notifications  // Since Plesk 17.0
{
    /**
     * @return array
     */
    public function getNotifications()
    {
        $notifications = [
            'uploadFailed' => [
                'title' => pm_Locale::lmsg('emailNotificationUploadFailedTitle'),
                'subject' => pm_Locale::lmsg('emailNotificationSubjectUploadNewBackupFailed'),
                'message' => pm_Locale::lmsg('emailNotificationBodyUploadNewBackupFailed'),
                'notifyAdmin' => true,
                'notifyResellers' => false,
                'notifyClients' => false,
                'notifyCustomEmail' => false,
                'customEmail' => '',
            ],
        ];

        return $notifications;
    }
}