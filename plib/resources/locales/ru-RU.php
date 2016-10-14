<?php
// Copyright 1999-2016. Parallels IP Holdings GmbH.

$messages = array(
    'tabStorage' => 'Хранилище',
    'tabSettings' => 'Настройки',
    'tabAbout' => 'О программе',
    'pageTitle' => 'Backup on Amazon S3',

    'settingsAwsAccessKeyIdHelp' => 'Пользователь <a target="_blank" rel="noopener noreferrer" href="https://console.aws.amazon.com/iam/home?#users">Amazon IAM</a> должен иметь права "AmazonS3FullAccess"',
    'settingsTestFailed' => 'Ошибка при тестировании соедениня с Amazon S3: %%message%%',
    'enabled' => 'Включено',
    'awsRegion' => 'Регион',
    'awsBucket' => 'Bucket',
    'awsAccessKeyId' => 'Access Key ID',
    'awsSecretAccessKey' => 'Secret Access Key',
    'awsSubFolder' => 'Поддиректория в Bucket\'e',
    'uploadNewBackups' => 'Выгружать новые сервеные бекапы в удаленное хранилище',
    'deleteSourceBackup' => 'Удалять исходные бекапы Plesk\'a после выгрузки (только инкрементальные)',
    'emailNotificationUploadFailedTitle' => 'Выгрузка бекапа на Amazon S3 завершилась с ошибкой',
    'emailNotificationSubjectUploadNewBackupFailed' => 'Ошибка при выгрузке новых бекапов Plesk\'a в Amazon S3',
    'emailNotificationBodyUploadNewBackupFailed' => 'Выгузка нового бекапа Plesk\'a в Amazon S3 завершилась с ошибкой: <error>',

    'us-east-1'      => 'us-east-1',
    'us-west-1'      => 'us-west-1',
    'us-west-2'      => 'us-west-2',
    'eu-west-1'      => 'eu-west-1',
    'eu-central-1'   => 'eu-central-1',
    'ap-southeast-1' => 'ap-southeast-1',
    'ap-southeast-2' => 'ap-southeast-2',
    'ap-northeast-1' => 'ap-northeast-1',
    'sa-east-1'      => 'sa-east-1',

    'SignatureDoesNotMatch' => 'Access Key или Secret Key не совпадают',
    'InvalidAccessKeyId' => 'Неправильный Access Key для пользователя ',
    'AccessDenied' => 'Доступ запрещен на Amazon S3, проверьте права и политику безопасности пользователя Amazon',
    'failedQueryNtpServer' => 'Ошибка при проверке локального времени: %%error%%',
    'gapInTime' => 'Обнаружена разница в локальном и реальном времени, локальное время %%local%%, реальное время %%real%%, разница во времени %%diff%%',
    'RequestTimeTooSkewed' => 'Слишком большая разница между текущим временем и временем сервера.',
    'settingsWasSuccessfullySaved' => 'Настройки были успешно сохранены',

    'date' => 'Дата',
    'backupName' => 'Имя бекапа',
    'isRemote' => 'В удаленном хранилище',
    'isPartialDownload' => 'Частичная загрузка',
    'isLocal' => 'В локальном хранилище',
    'isPartialUpload' => 'Частичная выгрузка',
    'size' => 'Размер',
    'mb' => 'Mb',
    'bytes' => 'байт',
    'actions' => '',
    'backupInvalid' => 'Бекап был создан с ошибками',

    'refreshPage' => 'Обновить',
    'uploadAll' => 'Выгрузить всё',
    'uploadAllButtonDescription' => 'Выгрузить все серверные бекапы в удаленное хранилище',
    'upload' => 'Выгрузить',
    'download' => 'Загрузить',
    'resumeUpload' => 'Продолжить Выгрузку',
    'resumeDownload' => 'Продолжить Загрузку',
    'cancelUpload' => 'Отменить Выгрузку',
    'cancelDownload' => 'Отменить Загрузку',

    'uploadAllTaskDone' => 'Все новые бекапы выгружены.',
    'uploadAllTaskError' => 'Ошибка при выгрузке новых бекапов: %%error%%',
    'uploadAllTaskRunning' => 'Выполняется выгрузка новых бекапов',
    'uploadAllTaskNotStarted' => 'Выгрузка новых бекапов не начата',

    'uploadTaskDone' => 'Выгрузка бекапа %%backupName%% завершена.',
    'uploadTaskError' => 'Ошибка при выгрузке бекапа %%backupName%%: %%error%%',
    'uploadTaskRunning' => 'Выполняется выгрузка бекапа %%backupName%%',
    'uploadTaskNotStarted' => 'Выгрузка бекапа %%backupName%% не начата',

    'resumeUploadTaskDone' => 'Продолженная выгрузка бекапа %%backupName%% завершена.',
    'resumeUploadTaskError' => 'Ошибка при продолжении выгрузки бекапа %%backupName%%: %%error%%',
    'resumeUploadTaskRunning' => 'Выполняется продолжение выгрузка бекапа %%backupName%%',
    'resumeUploadTaskNotStarted' => 'Продолжение выгрузки бекапа %%backupName%% не начата',

    'cancelUploadTaskDone' => 'Отмена выгрузки бекапа %%backupName%% завершена.',
    'cancelUploadTaskError' => 'Ошибка отмены выгрузки бекапа %%backupName%%: %%error%%',
    'cancelUploadTaskRunning' => 'Выполняется отмена выгрузки бекапа %%backupName%%',
    'cancelUploadTaskNotStarted' => 'Отмена выгрузки бекапа %%backupName%% не начата',

    'downloadTaskDone' => 'Загрузка бекапа %%remotePath%% завершена.',
    'downloadTaskError' => 'Ошибка загузки бекапа %%remotePath%%: %%error%%',
    'downloadTaskRunning' => 'Выполняется загузка бекапа %%remotePath%%',
    'downloadTaskNotStarted' => 'Загрузка бекапа %%remotePath%% не начата',

    'resumeDownloadTaskDone' => 'Продолжение загрузки бекапа %%remotePath%% завершено.',
    'resumeDownloadTaskError' => 'Ошибка при продолжении загрузки бекапа %%remotePath%%: %%error%%',
    'resumeDownloadTaskRunning' => 'Выполняется продолжение загрузки бекапа %%remotePath%%',
    'resumeDownloadTaskNotStarted' => 'Продолжение загрузки бекапа %%remotePath%% не начато',

    'cancelDownloadTaskDone' => 'Отмена выгрузки бекапа %%remotePath%% завершено.',
    'cancelDownloadTaskError' => 'Ошибка отмены выгрузки бекапа %%remotePath%%: %%error%%',
    'cancelDownloadTaskRunning' => 'Выполняется отмена выгрузки бекапа %%remotePath%%',
    'cancelDownloadTaskNotStarted' => 'Отмена выгрузки бекапа %%remotePath%% не начата',

    'BackupUploadFailedNoDiskSpaceForExport' => 'Свободного дискового пространства %%freeSpace%% байт в директории выгрузки %%uploadDir%% недостаточно для экспорта бекапа %%backupName%% размером %%backupSize%% байт',

    'PleskBackupReturnCodeError' => 'Ошибка импорта/экспорта бекапа',

    'PleskBackupReturnCodeImportedExist' => 'Импортируемый бекап %%filePath%% уже существует в локальном хранилише Plesk\'a %%dumpD%%',
    'PleskBackupReturnCodeImportedObjectNotMatch' => 'Импортируемый бекап %%filePath%% не соответвует',
    'PleskBackupReturnCodeImportWrongPassword' => 'Неправильный пароль бекапа %%filePath%%',
    'PleskBackupReturnCodeImportDeprecatedDumpVersion' => 'Версия бекапа %%filePath%% не поддерживается',
    'PleskBackupReturnCodeImportWinNativeMailContentSkipped' => 'Нативный почтовый контент был пропущен при импорте %%filePath%%',
    'PleskBackupReturnCodeImportErrorSign' => 'Плохая подпись бекапа %%filePath%%',
    'PleskBackupReturnCodeImportNotWellFormedXml' => 'Непаврильно сформированный XML при импорте бекапа %%filePath%%',
    'PleskBackupReturnCodeImportDenied' => 'Запрет импорта бекапа %%filePath%%',

    'PleskBackupReturnCodeTransportPermissionDenied' => 'Transport Permission Denied',
    'PleskBackupReturnCodeTransportWrongPassword' => 'Неправильный пароль транспорт',
    'PleskBackupReturnCodeTransportWrongLogin' => 'Неправильный логин транспорта',
    'PleskBackupReturnCodeTransportResolveHost' => 'Ошибка разрешения имени хоста',
    'PleskBackupReturnCodeTransportUnableConnect' => 'Ошибка подключения',
    'PleskBackupReturnCodeTransportNetworkError' => 'Сетевая ошибка транспорта',
    'PleskBackupReturnCodeTransportFileNotExist' => 'Файл не существует',

    'PleskBackupReturnCodeRepoDumpNotExist' => 'Бекап %%backupName%% не существует',
    'PleskBackupReturnCodeRepoBadDump' => 'Бекап импортированный из %%dstFilePath%% плохой: <a href=\'/admin/backup/restore/?type=local&dumpId=%%backupName%%\'>%%backupName%%</a>',
    'PleskBackupReturnCodeRepoDumpExist' => 'Бекап уже существует %%backupName%%',
    'PleskBackupReturnCodeRepoPathTooLong' => 'Слишком длинный путь %%dumpD%%',

    'about' => 'Это расширение позволяет сохранять Ваши серверные бекапы на Amazon S3',
    'feedback' => 'В случае каких-либо проблем с расширением Вы можете создать issue в репозитории на <a rel="noopener noreferrer" target="_blank" href="https://github.com/plesk/ext-backup-amazon">GitHub</a>',
    'faq' => 'FAQ',
    'faqFullNotDeletingQuestion' => 'Q: Почему "полные" бекапы не удаляются после выгрузки?',
    'faqFullNotDeletingAnswer' => 'A: Потому что Plesk может создавать новые инкрементальные бекапы только на основе полного бекапа. "Инкрементальные" бекапы позволяют достичь более эффектиыного использования хранилиша и уменьшить затраты на хранение.',

    'faqPartialDownloadPartialUploadQuestion' => 'Q: Что значат поля "Частичная загрузка" и "Частичная выгрузка"?',
    'faqPartialDownloadPartialUploadAnswer' => 'A: Эти поля указывают на неудачную и незавершенную загрузку или выгрузку бекапа.',

    'faqPartialUploadStoringQuestion' => 'Q: В случае неудачной выгрузки на Amazon S3, где хранятся выгруженные данные?',
    'faqPartialUploadStoringAnswer' => 'A: Экспортированный бекап продолжает хранится в директории DUMP_TMP_D/upload, по-умолчанию это /tmp/upload. Также Amazon S3 продолжает хранить все вызруженные данные (и взымать плату за них) до тех пор пока выгрузка не будет завершена или отменена.',

    'faqPartialDownloadStoringQuestion' => 'Q:  В случае неудачной загрузки с Amazon S3, где хранятся загруженные данные?',
    'faqPartialDownloadStoringAnswer' => 'A: В директории DUMP_TMP_D/download, по умолчанию это /tmp/download. Вы можете изменить настройку DUMP_TMP_D в /etc/psa/psa.conf или реестре Windows HKLM\SOFTWARE\Wow6432Node\PLESK\PSA Config\Config',

    'faqResumeUploadDownloadQuestion' => 'Q: Что я могу сделать в случае неудачной загрузки или выгрузки?',
    'faqResumeUploadDownloadAnswer' => 'A: Расширение предоставляет возможность для Отмены или Продолжения неудачной загрузки или выгрузки бекапа. Выполнение операции продолжится с последних успешно отправленных байтов.',

    'faqResumeBackupPlanQuestion' => 'Q: Какая моя лучшая стратегия бекапа?',
    'faqResumeBackupPlanAnswer' => 'A: 
        <ol>
          <li>Включите это расширение на вкладке "Settings"</li>
          <li>Сделайте полный серверный бекап средствами Plesk\'a</li>  
          <li>Выгрузите этот полный бекап на Amazon S3 используя это расширение с вкладки "Хранилище"</li>  
          <li><a href="/admin/backup/schedule" target="_blank">Настройте регулярный инкрементное резервное копирование</a> для создания новых инрементальных бекапов каждый день.</li>  
          <li>Начиная с этого момента все новые инрементальные бекапы будут выгружатся в Amazon S3. Имейте ввиду что Вы всегда должны продолжать хранить хотябы один полный бекап сервера в локальном хранилище Plesk\'a.</li>  
        </ol>',

    'faqResumeRestorePlanQuestion' => 'Q: Какая моя лучшая стратегия для восстановления из бекапа?',
    'faqResumeRestorePlanAnswer' => 'A: Используя это расширение загрузите последний полный бекап и последний инрементальный бекап из Amazon S3. После загурзки эти бекапы станут доступны в <a target="_blank" href="/admin/backup/list">хранилище Plesk\'a</a> и Вы сможете восставить сервер из последнего инкреентального бекапа.',

    'faqLogFileLocationQuestion' => 'Q: Где я могу найти журнал операций с Amazon S3?',
    'faqLogFileLocationAnswer' => 'A: Все операции с Amazon S3 сохраняются в <b>/usr/local/psa/admin/sbin/modules/backup-amazon/backup-on-amazon.log</b>',

);