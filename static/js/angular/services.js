'use strict';

angular.module('app').factory('Frame', function () {
	var onOpenCB, onCloseCB, onMessageCB;

	var loc = window.location, location;
	location = "ws:";
	location += "//" + loc.host;
	location += loc.pathname + "socket";
	var ws = new WebSocket(location);
	ws.onopen = function () {
		if(onOpenCB !== undefined)
		{
			onOpenCB();
		}
	};
	ws.onclose = function () {
		if(onCloseCB !== undefined)
		{
			onCloseCB();
		}
	};
	ws.onmessage = function (m) {
		if(onMessageCB !== undefined)
		{
			onMessageCB(m);
		}
	};

	return{
		setOnOpenCB: function(cb){
			onOpenCB = cb;
		},
		setOnCloseCB: function(cb){
			onCloseCB = cb;
		},
		setOnMessageCB: function(cb){
			onMessageCB = cb;
		}
	};
});