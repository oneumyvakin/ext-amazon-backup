// Copyright 1999-2016. Parallels IP Holdings GmbH.

var BackupAmazonAnimationSrc = '/modules/backup-amazon/icons/indicator.gif';

function BackupAmazonPostData(node) {

    if (node.disabled) {
        return
    }

    node.disabled = true;
    innerImg = node.querySelector("img");
    innerImg.src = BackupAmazonAnimationSrc;

    var data = {};
    [].forEach.call(node.attributes, function(attr) {
        if (/^data-/.test(attr.name)) {
            var camelCaseName = attr.name.substr(5).replace(/-(.)/g, function ($0, $1) {
                return $1.toUpperCase();
            });
            data[camelCaseName] = attr.value;
        }
    });

    console.log(data);

    var XHR = new XMLHttpRequest();
    var FD  = new FormData();

    var forgeryToken = $('forgery_protection_token');
    FD.append('forgery_protection_token', forgeryToken.content);

    // We push our data into our FormData object
    for(name in data) {
        FD.append(name, data[name]);
    }

    // We define what will happen if the data are successfully sent
    XHR.addEventListener('load', function(event) {
        setTimeout("location.reload();", 3000);

    });

    // We define what will happen in case of error
    XHR.addEventListener('error', function(event) {
        alert('Oups! Something goes wrong.');
    });

    // We setup our request
    XHR.open('POST', data['action']);

    // We just send our FormData object, HTTP headers are set automatically
    XHR.send(FD);


}

