var sessionData ='{}';

window.onmessage = function(event){
	var msg = event.data;
	var uData = JSON.parse(msg);
	sessionData = uData.login + '|' +  uData.token;
	var ssId = localStorage.getItem('ssID');
	// document.cookie = 'X-WBR=' + (ssId==null ? '':ssId) + '; path=/';
	// document.cookie = 'X-TAData=' + (sessionData==null ? '':sessionData) + '; path=/';
	ws = new WebSocket('wss://'+document.location.host+'/ws?ss='+ssId+'&sd='+btoa(sessionData) );
	ws.onerror = wsOnError;
	ws.onopen = wsOnOpen;
	ws.onmessage = wsOnMessage;
};

isWsLog = true;

//#####################################################################################

var wsOnError = function(event) {
	//notify('error','WS соединение с сервером не установленно.');
	console.log('ws: Error');
};
var wsOnOpen = function(e) {
	// e.target.send("WS: connection success");
	e.target.send('{"action":"connect","success":true, "data":{} }');
};
var wsOnMessage = function(ws) {
	var msg = ws.data;
	if (isWsLog) console.log('ws:' + msg);
	var cmd = jQuery.parseJSON(msg);
	if(cmd) {
		// if(cmd.action == 'seatStateUpdate') {
		//     //updateSeatState(cmd.data);
		// }
		if(cmd.action == 'setDeviID') {
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
};

$(document).ready(function(){
	//$('#fMap').height(document.body.clientHeight);
	var svg = document.getElementById("fMap");
	svg.addEventListener("load",function(){
		var svgDoc = svg.contentDocument;
		//console.log(svgDoc);
		//$("g#canvas",svgDoc).mousedown(mapDarg);

		$("rect.cover",svgDoc).mouseover(defOver);
		$("rect.cover",svgDoc).mouseout(defOver);


		// Get one of the SVG items by ID;
		$("rect.cover",svgDoc).click(defId);


		$("div#upBut").click( function() {
			mapMoveY(svgDoc,-100);
		});
		$("div#downBut").click( function() {
			mapMoveY(svgDoc,100);
		});
		$("div#leftBut").click( function() {
			mapMoveX(svgDoc,-100);
		});
		$("div#rightBut").click( function() {
			mapMoveX(svgDoc,100);
		});

		$("div#magBut").click( function() {
			magMap(svgDoc,'+');
		});
		$("div#umagBut").click( function() {
			magMap(svgDoc,"-");
		});

	}, false);

});
n=1;
m=1;
var defId = function(e) {
	var r = e.target;
	if(r) {
		console.log( r.id );
	}	
	
};

var mapDarg = function(e) {
	var r = e.target;
	var svgDoc = svg.contentDocument;
	var ss = $("svg",svgDoc);
	var vb = ss.attr("viewBox");
	var mm = vb.split(' ');

	r.onmousemove = function(e) {
		var dy = e.movementY;
		var cy = parseInt(mm[1]);
		var maxy = parseInt(mm[3]);
		var y = ( cy - dy*300);
		if(y<0) y =0;
		if(y>maxy) y =maxy;
		var tt = mm[0] + ' ' + y + ' ' +mm[2]+ ' ' +mm[3];
		ss.attr("viewBox", tt);

		return false;
	};
	r.onmouseup = function() {
		//document.removeEventListener('mousemove', onMouseMove);
		r.onmousemove = null;
	};

};

var  defOver= function(e) {
	var r = e.target;
	if(r) {
		if(e.type =="mouseover") {
			r.classList.add('block');
		} else {
			r.classList.remove('block');
		}

	}

};

var mapMoveX = function (svgDoc, dx) {
	var ss = $("svg",svgDoc);
	var w = parseInt(ss.attr("width"));
	var vb = ss.attr("viewBox");
	var mm = vb.split(' ');
	var cx = parseInt(mm[0]);
	var maxx = w - parseInt(mm[2]);
	var x = ( cx + dx);
	if(x<0) x = 0;
	if(x>maxx) x = maxx;
	var tt =  x + ' ' + mm[1] + ' ' +mm[2]+ ' ' +mm[3];
	ss.attr("viewBox", tt);
};
var mapMoveY = function (svgDoc, dy) {
	var ss = $("svg",svgDoc);
	var h = parseInt(ss.attr("height"));
	var vb = ss.attr("viewBox");
	var mm = vb.split(' ');
	var cy = parseInt(mm[1]);
	var maxy = h-parseInt(mm[3]);
	var y = ( cy + dy);
	if(y<0) y =0;
	if(y>maxy) y =maxy;
	var tt = mm[0] + ' ' + y + ' ' +mm[2]+ ' ' +mm[3];
	ss.attr("viewBox", tt);
}

var magMap = function (svgDoc, c) {
	var ss = $("svg",svgDoc);
	var vb = ss.attr("viewBox");
	var w = parseInt(ss.attr("width"))*0.1;
	var h= parseInt(ss.attr("height"))*0.1;
	var mm = vb.split(' ');
	var cx = parseInt(mm[2]);
	var cy = parseInt(mm[3]);
	//var maxy = parseInt(mm[3]);

	if(c=='+') {
		cx -= w;
		cy -= h;
	};
	if(c=='-'){
		cx += w;
		cy += h;
	};
	var tt = mm[0] + ' ' + mm[1] + ' ' +cx+ ' ' +cy;
	ss.attr("viewBox", tt);
};

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
};