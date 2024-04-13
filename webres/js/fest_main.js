var sessionData ='{}';

window.onmessage = function(event){
    var msg = event.data;
    sessionData=msg;
    var ssId = localStorage.getItem('ssID');
    document.cookie = 'X-WBR=' + (ssId==null ? '':ssId) + '; path=/ws';
    document.cookie = 'X-TAData=' + (sessionData==null ? '':sessionData) + '; path=/ws';
    ws = new WebSocket('wss://'+document.location.host+'/ws?ss='+ssId+'&sd='+btoa(sessionData) );
    ws.onerror = wsOnError;
    ws.onopen = wsOnOpen;
    ws.onmessage = wsOnMessage;
};

isWsLog = true;

$(document).ready(function(){
    resize();
});
//#####################################################################################

var wsOnError = function(event) {
    //notify('error','WS соединение с сервером не установленно.');
    console.log('ws: Error');
};
var wsOnOpen = function(e) {
    // e.target.send("WS: connection success");
    e.target.send('{"action":"connect","success":true, "data": '+sessionData+' }');
};
var wsOnMessage = function(ws) {
    var msg = ws.data;
    if (isWsLog) console.log('ws:' + msg);
    var cmd = jQuery.parseJSON(msg);
    if(cmd) {

        if(cmd.action == 'createSession' && cmd.success) {
            localStorage.setItem('ssID', cmd.device);
        }
        if(cmd.action == 'updateSession' && cmd.success) {
            localStorage.setItem('ssID', cmd.device);
        }
        if(cmd.action == 'closeSession' ) {
            localStorage.setItem('ssID', '');
            sendMsgToParent(cmd.action);
        }
    }

};

var sendMsgToParent = function(msg) {
    window.top.postMessage(msg, '*');
}

var notify = function(type, msg) {
    var nn = Math.floor(Math.random() * 1000);
    var attr = { 'id':'notif'+nn};

    attr.class = 'notifyBox body'+type;
    //attr.text = '<div class="header"+type></div>'+'<div class="msg"+type>'+msg+'</div>';
    $("body").append( $('<div>', attr) );
    var box = $("div#notif"+nn);

    attr = { class: "header"+type, text: '!!!!!!!!'}
    box.append($('<div>', attr));

    attr = { class: "msg"+type, text: msg}
    box.append($('<div>', attr));

    box.animate({top: '20px' }, 500, function () {
        setTimeout(function() {
            box.animate({right: '-350px' }, 1000, function () {
                box.remove();
            });
        }, 4000);
    });
};

var resize = function() {
    $('div#tabs').attr('style', 'height:'+parseInt(window.innerHeight-26)+'px;');
};
window.onresize = resize;

function getElementPosition(elemId) {
    var elem = document.getElementById(elemId);

    var w = elem.offsetWidth;
    var h = elem.offsetHeight;

    var l = 0;
    var t = 0;

    while (elem)
    {
        l += elem.offsetLeft;
        t += elem.offsetTop;
        elem = elem.offsetParent;
    }

    return {"left":l, "top":t, "width": w, "height":h};
}