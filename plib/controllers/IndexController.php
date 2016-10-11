<?php
// Copyright 1999-2016. Parallels IP Holdings GmbH.

class IndexController extends pm_Controller_Action
{
    public function init()
    {
        $this->_accessLevel = 'admin';
        
        parent::init();
        
        
        $this->view->pageTitle = $this->lmsg('pageTitle');
        
        $this->view->tabs = [
            [
                'title' => $this->lmsg('tabStorage'),
                'action' => 'storage',
            ],
            [
                'title' => $this->lmsg('tabSettings'),
                'action' => 'settings',
            ],
            [
                'title' => $this->lmsg('tabAbout'),
                'action' => 'about',
            ],
        ];
    }

    public function indexAction()
    {
        if (!pm_Settings::get('backupToAmazonEnabled') || pm_Settings::get('apiKeyBecameInvalid')) {
            $this->_forward('settings');
            return;
        }

        $this->_forward('storage');
    }

    public function storageAction()
    {
        if (!pm_Settings::get('backupToAmazonEnabled')) {
            $this->_forward('settings');
            return;
        }

        $this->view->tools = [
            [
                'title' => $this->lmsg('uploadAll'),
                'description' => $this->lmsg('uploadAllButtonDescription'),
                'icon' => pm_Context::getBaseUrl() . '/icons/upload-files.png',
                'link' => $this->view->getHelper('baseUrl')->moduleUrl(['action' => 'uploadAll']),
            ],
        ];

        $this->view->list = $this->_getStorageList();
    }


    public function storageDataAction()
    {
        $list = $this->_getStorageList();
        // Json data from pm_View_List_Simple
        $this->_helper->json($list->fetchData());
    }
    
    public function settingsAction()
    {
        if (pm_Settings::get('apiKeyBecameInvalid') && !$this->_status->hasMessage($this->lmsg('apiKeyBecameInvalid'))) {
            $this->_status->addError($this->lmsg('apiKeyBecameInvalid'));
        }
        
        $this->view->settingsAwsAccessKeyIdHelp = $this->lmsg('settingsAwsAccessKeyIdHelp');

        $form = new Modules_BackupAmazon_SettingsForm();

        $form->addElement('checkbox', 'backupToAmazonEnabled', [
            'label' => $this->lmsg('enabled'),
            'value' => pm_Settings::get('backupToAmazonEnabled'),
        ]);

        $form->addElement('text', 'awsBucket', [
            'label' => $this->lmsg('awsBucket'),
            'value' => pm_Settings::get('awsBucket'),
            'required' => true,
            'validators' => [
                ['NotEmpty', true],
            ],
        ]);

        $form->addElement('text', 'awsAccessKeyId', [
            'label' => $this->lmsg('awsAccessKeyId'),
            'value' => pm_Settings::get('awsAccessKeyId'),
            'required' => true,
            'validators' => [
                ['NotEmpty', true],
            ],
        ]);

        $form->addElement('text', 'awsSecretAccessKey', [
            'label' => $this->lmsg('awsSecretAccessKey'),
            'value' => pm_Settings::get('awsSecretAccessKey'),
            'required' => true,
            'validators' => [
                ['NotEmpty', true],
            ],
        ]);

        $form->addElement('text', 'awsSubFolder', [
            'label' => $this->lmsg('awsSubFolder'),
            'value' => pm_Settings::get('awsSubFolder'),
            'required' => false,
        ]);

        $form->addElement('checkbox', 'uploadNewBackups', [
            'label' => $this->lmsg('uploadNewBackups'),
            'value' => pm_Settings::get('uploadNewBackups'),
            'required' => false,
            'validators' => [
                ['NotEmpty', true],
            ],
        ]);

        $form->addElement('checkbox', 'deleteSourceBackup', [
            'label' => $this->lmsg('deleteSourceBackup'),
            'value' => pm_Settings::get('deleteSourceBackup'),
            'required' => false,
            'validators' => [
                ['NotEmpty', true],
            ],
        ]);

        $form->addControlButtons([
            'cancelLink' => pm_Context::getModulesListUrl(),
        ]);

        if ($this->getRequest()->isPost() && $form->isValid($this->getRequest()->getPost())) {
            pm_Settings::set('apiKeyBecameInvalid', '');
            pm_Settings::set('backupToAmazonEnabled', $form->getValue('backupToAmazonEnabled'));
            pm_Settings::set('awsBucket', $form->getValue('awsBucket'));
            pm_Settings::set('awsAccessKeyId', $form->getValue('awsAccessKeyId'));
            pm_Settings::set('awsSecretAccessKey', $form->getValue('awsSecretAccessKey'));
            pm_Settings::set('awsSubFolder', $form->getValue('awsSubFolder'));
            pm_Settings::set('uploadNewBackups', $form->getValue('uploadNewBackups'));
            pm_Settings::set('deleteSourceBackup', $form->getValue('deleteSourceBackup'));

            $this->_status->addMessage('info', $this->lmsg('settingsWasSuccessfullySaved'));
            $this->_helper->json(['redirect' => pm_Context::getBaseUrl()]);
        }

        $this->view->form = $form;
    }

    public function aboutAction()
    {
        if (pm_Settings::get('apiKeyBecameInvalid') && !$this->_status->hasMessage($this->lmsg('apiKeyBecameInvalid'))) {
            $this->_status->addError($this->lmsg('apiKeyBecameInvalid'));
        }
        
        $this->view->about = $this->lmsg('about');
        $this->view->feedback = $this->lmsg('feedback');
        $this->view->faq = $this->lmsg('faq');

        $this->view->faqItems = [
            [
                'question' => $this->lmsg('faqFullNotDeletingQuestion'),
                'answer' => $this->lmsg('faqFullNotDeletingAnswer'),
            ],
            [
                'question' => $this->lmsg('faqPartialDownloadPartialUploadQuestion'),
                'answer' => $this->lmsg('faqPartialDownloadPartialUploadAnswer'),
            ],
            [
                'question' => $this->lmsg('faqPartialUploadStoringQuestion'),
                'answer' => $this->lmsg('faqPartialUploadStoringAnswer'),
            ],
            [
                'question' => $this->lmsg('faqPartialDownloadStoringQuestion'),
                'answer' => $this->lmsg('faqPartialDownloadStoringAnswer'),
            ],
            [
                'question' => $this->lmsg('faqResumeUploadDownloadQuestion'),
                'answer' => $this->lmsg('faqResumeUploadDownloadAnswer'),
            ],
            [
                'question' => $this->lmsg('faqResumeBackupPlanQuestion'),
                'answer' => $this->lmsg('faqResumeBackupPlanAnswer'),
            ],
            [
                'question' => $this->lmsg('faqResumeRestorePlanQuestion'),
                'answer' => $this->lmsg('faqResumeRestorePlanAnswer'),
            ],
            [
                'question' => $this->lmsg('faqLogFileLocationQuestion'),
                'answer' => $this->lmsg('faqLogFileLocationAnswer'),
            ],
        ];

    }

    private function _throwUiError($errStruct) {
        $msg = $errStruct['message'];

        if ($errStruct['locale_key'] <> '') {
            $msg = pm_Locale::lmsg($errStruct['locale_key'], $errStruct['locale_args']);
        }

        $this->_status->addError($msg);
    }

    /*
     * Format action buttons for Storage List
     * @param $backup array
     * @return string
     */
    private function _getListColumnActionButtons(array $backup)
    {
        if ($backup['IsRemote'] && $backup['IsLocal'] && !$backup['IsPartialUpload'] && !$backup['IsPartialDownload']) {
            return '';
        }

        $buttons = '';
        $htdocs = pm_Context::getBaseUrl();

        if ($backup['InProgress']) {
            $buttons = "<img src='{$htdocs}/icons/indicator.gif' /> " . $buttons;
            return $buttons;
        }

        if ($backup['IsLocal'] && !$backup['IsLocalInvalid'] && !$backup['IsRemote'] && !$backup['IsPartialUpload'] && !$backup['IsPartialDownload']) {
            $buttons .= ' <a href="#" '
                .' data-action="' . $this->_helper->url('upload') . '"'
                .' data-backup-name="' . $backup['Backup']['Fullname'] . '"'
                .' onclick="BackupAmazonPostData(this)"'
                .'>'
                ."<img src='{$htdocs}/icons/upload.png' /> "
                . pm_Locale::lmsg('upload') . '</a>';
        }

        if (!$backup['IsLocal'] && $backup['IsRemote'] && !$backup['IsPartialUpload'] && !$backup['IsPartialDownload']) {
            $buttons .= ' <a href="#" '
                .' data-action="' . $this->_helper->url('download') . '"'
                .' data-remote-path="' . $backup['RemotePath'] . '"'
                .' onclick="BackupAmazonPostData(this)"'
                .'>'
                ."<img src='{$htdocs}/icons/download.png' /> "
                . pm_Locale::lmsg('download') . '</a>';
        }

        if ($backup['IsLocal'] && !$backup['IsRemote'] && $backup['IsPartialUpload'] && !$backup['IsPartialDownload']) {
            $buttons .= ' <a href="#" '
                .' data-action="' . $this->_helper->url('resumeUpload') . '"'
                .' data-backup-name="' . $backup['Backup']['Fullname'] . '"'
                .' data-upload-id="' . $backup['UploadId'] . '"'
                .' onclick="BackupAmazonPostData(this)"'
                .'>'
                ."<img src='{$htdocs}/icons/upload.png' /> "
                . pm_Locale::lmsg('resumeUpload') . '</a>';
        }

        if (!$backup['IsLocal'] && $backup['IsRemote'] && !$backup['IsPartialUpload'] && $backup['IsPartialDownload']) {
            $buttons .= ' <a href="#" '
                .' data-action="' . $this->_helper->url('resumeDownload') . '"'
                .' data-remote-path="' . $backup['RemotePath'] . '"'
                .' onclick="BackupAmazonPostData(this)"'
                .'>'
                ."<img src='{$htdocs}/icons/download.png' /> "
                . pm_Locale::lmsg('resumeDownload') . '</a>';
        }

        if ($backup['IsPartialUpload']) {
            $buttons .= ' <a href="#" '
                .' data-action="' . $this->_helper->url('cancelUpload') . '"'
                .' data-backup-name="' . $backup['Name'] . '"'
                .' data-remote-path="' . $backup['RemotePath'] . '"'
                .' onclick="BackupAmazonPostData(this)"'
                .'>'
                ."<img src='{$htdocs}/icons/cancel.png' /> "
                . pm_Locale::lmsg('cancelUpload') . '</a>';
        }

        if ($backup['IsPartialDownload']) {
            $buttons .= ' <a href="#" '
                .' data-action="' . $this->_helper->url('cancelDownload') . '"'
                .' data-backup-name="' . $backup['Name'] . '"'
                .' data-remote-path="' . $backup['RemotePath'] . '"'
                .' onclick="BackupAmazonPostData(this)"'
                .'>'
                ."<img src='{$htdocs}/icons/cancel.png' /> "
                . pm_Locale::lmsg('cancelDownload') . '</a>';
        }

        return $buttons;
    }

    /*
     * Format backup name for List
     * @param $backup array
     * @return string
     */
    private function _getListColumnBackupName(array $backup) {
        $backupName = $backup['Name'];
        if ($backup['IsLocal']) {
            $backupName = "<a href='/admin/backup/restore/?type=local&dumpId=${backupName}'>${backupName}</a>";
        }

        if ($backup['IsRemote'] && !$backup['IsPartialUpload']) {
            $bucket = pm_Settings::get('awsBucket');
            $prefix = pm_Settings::get('awsSubFolder');
            $awsUrl = "https://console.aws.amazon.com/s3/home?bucket=${bucket}&prefix=${prefix}";
            $backupName = "<a target='_blank' href='${awsUrl}'>${backupName}</a>";
        }

        return $backupName;
    }

    /*
     * Format size for List
     * @param $size integer
     * @return string
     */
    private function _getListColumnSize($size) {
        if ((int)($size) < 1000000) {
            return (string)$size . ' ' . pm_Locale::lmsg('bytes');
        }

        $mbSize = round((float)($size) / 1000000);
        return $mbSize . ' ' . pm_Locale::lmsg('mb');
    }

    /*
     * Format Local for List
     * @param $backup array
     * @return string
     */
    private function _getListColumnLocal(array $backup) {

        $htdocs = pm_Context::getBaseUrl();
        $local = '';

        if ($backup['IsLocal']) {
            $local = "<img src='{$htdocs}/icons/on.png' />";
        }

        if ($backup['IsLocalInvalid']) {
            $tooltip = pm_Locale::lmsg('backupInvalid');
            $local = "<img src='{$htdocs}/icons/warning.png' title='${tooltip}' />";
        }

        return $local;
    }

    private function _getStorageList()
    {

        $data = [];
        $storage = Modules_BackupAmazon_Helper::storage();
        if (isset($storage['is_error'])) {
            $this->_throwUiError($storage);
            return new pm_View_List_Simple($this->view, $this->_request);
        }


        $htdocs = pm_Context::getBaseUrl();
        foreach ($storage as $backupName => $item) {

            $data[$backupName] = [
                'column-1' => $item['Backup']['CreationDate'],
                'column-2' => $this->_getListColumnBackupName($item),
                'column-3' => $item['IsRemote'] ? "<img src='{$htdocs}/icons/on.png' />" : '',
                'column-4' => $item['IsPartialDownload'] && !$item['InProgress'] ? "<img src='{$htdocs}/icons/warning.png' />" : '',
                'column-5' => $this->_getListColumnLocal($item),
                'column-6' => $item['IsPartialUpload'] && !$item['InProgress'] ? "<img src='{$htdocs}/icons/warning.png' />" : '',
                'column-7' => $this->_getListColumnSize($item['Backup']['Size']),
                'column-8' => $this->_getListColumnActionButtons($item),
            ];
        }

        if (!count($data) > 0) {
            return new pm_View_List_Simple($this->view, $this->_request);
        }

        $options = [
            'defaultSortField' => 'column-1',
            'defaultSortDirection' => pm_View_List_Simple::SORT_DIR_DOWN,
        ];
        $list = new pm_View_List_Simple($this->view, $this->_request, $options);
        $list->setData($data);
        $list->setColumns([
            pm_View_List_Simple::COLUMN_SELECTION,
            'column-1' => [
                'title' => $this->lmsg('date'),
                'noEscape' => true,
                'searchable' => true,
                'sortable' => true,
            ],
            'column-2' => [
                'title' => $this->lmsg('backupName'),
                'noEscape' => true,
                'searchable' => true,
                'sortable' => true,
            ],
            'column-3' => [
                'title' => $this->lmsg('isRemote'),
                'noEscape' => true,
                'searchable' => false,
                'sortable' => true,
            ],
            'column-4' => [
                'title' => $this->lmsg('isPartialDownload'),
                'noEscape' => true,
                'searchable' => false,
                'sortable' => true,
            ],
            'column-5' => [
                'title' => $this->lmsg('isLocal'),
                'noEscape' => true,
                'sortable' => true,
            ],
            'column-6' => [
                'title' => $this->lmsg('isPartialUpload'),
                'noEscape' => true,
                'searchable' => false,
                'sortable' => true,
            ],
            'column-7' => [
                'title' => $this->lmsg('size'),
                'sortable' => true,
            ],
            'column-8' => [
                'title' => $this->lmsg('actions'),
                'noEscape' => true,
                'searchable' => false,
                'sortable' => false,
                
            ],
        ]);

        $list->setDataUrl(['action' => 'storage-data']);
        return $list;
    }

    public function uploadallAction()
    {
        $taskManager = new pm_LongTask_Manager();
        $task = new Modules_BackupAmazon_Task_UploadAll();
        $taskManager->start($task);

        $this->_redirect(pm_Context::getBaseUrl());
    }

    public function uploadAction()
    {
        $messages[] = ['status' => 'info', 'content' => "Upload started"];

        $backupName = $this->_getParam('backupName');

        $taskManager = new pm_LongTask_Manager();
        $task = new Modules_BackupAmazon_Task_Upload();
        $task->setParam('backupName', $backupName);
        $taskManager->start($task);

        $this->_helper->json([
                'status' => 'success',
                'statusMessages' => $messages,
                'backupName' => $backupName,
            ]
        );
    }

    public function resumeuploadAction()
    {
        $messages[] = ['status' => 'info', 'content' => "Resume upload started"];

        $backupName = $this->_getParam('backupName');

        $taskManager = new pm_LongTask_Manager();
        $task = new Modules_BackupAmazon_Task_ResumeUpload();
        $task->setParam('backupName', $backupName);
        $taskManager->start($task);

        $this->_helper->json([
                'status' => 'success',
                'statusMessages' => $messages,
                'backupName' => $backupName,
            ]
        );
    }

    public function canceluploadAction()
    {
        $messages[] = ['status' => 'info', 'content' => "Cancel upload started"];

        $backupName = $this->_getParam('backupName');
        $uploadId =  $this->_getParam('uploadId');

        $taskManager = new pm_LongTask_Manager();
        $task = new Modules_BackupAmazon_Task_CancelUpload();
        $task->setParam('backupName', $backupName);
        $task->setParam('uploadId', $uploadId);
        $taskManager->start($task);

        $this->_helper->json([
                'status' => 'success',
                'statusMessages' => $messages,
                'backupName' => $backupName,
            ]
        );
    }

    public function downloadAction()
    {
        $messages[] = ['status' => 'info', 'content' => "Download started"];

        $remotePath = $this->_getParam('remotePath');

        $taskManager = new pm_LongTask_Manager();
        $task = new Modules_BackupAmazon_Task_Download();
        $task->setParam('remotePath', $remotePath);
        $taskManager->start($task);

        $this->_helper->json([
                'status' => 'success',
                'statusMessages' => $messages,
                'backupName' => $remotePath,
            ]
        );
    }

    public function canceldownloadAction()
    {
        $messages[] = ['status' => 'info', 'content' => "Cancel upload started"];

        $backupName = $this->_getParam('backupName');
        $remotePath = $this->_getParam('remotePath');

        $taskManager = new pm_LongTask_Manager();
        $task = new Modules_BackupAmazon_Task_CancelDownload();
        $task->setParam('backupName', $backupName);
        $task->setParam('remotePath', $remotePath);
        $taskManager->start($task);

        $this->_helper->json([
                'status' => 'success',
                'statusMessages' => $messages,
                'backupName' => $backupName,
            ]
        );
    }

    public function resumedownloadAction()
    {
        $messages[] = ['status' => 'info', 'content' => "Resume upload started"];

        $backupName = $this->_getParam('backupName');
        $remotePath =  $this->_getParam('remotePath');

        $taskManager = new pm_LongTask_Manager();
        $task = new Modules_BackupAmazon_Task_ResumeDownload();
        $task->setParam('backupName', $backupName);
        $task->setParam('remotePath', $remotePath);

        $taskManager->start($task);

        $this->_helper->json([
                'status' => 'success',
                'statusMessages' => $messages,
                'backupName' => $backupName,
            ]
        );
    }
}
