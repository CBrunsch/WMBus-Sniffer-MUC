'use strict';

function MainCtrl($scope, $rootScope, $filter, $http, $location, Frame, ngTableParams) {
    $scope.tableParams = new ngTableParams({
        page: 1,            // show first page
        total: 0,           // length of data
        count: 10,          // count per page
        sorting: {
            ID: 'desc'     // initial sorting
        }
    });

    $scope.visible = function(frameID) {
        if ($scope.visible[frameID] == 1) {
            $scope.visible[frameID] = 0;
        } else {
            $scope.visible[frameID] = 1;
        }
    };

    var cachedData = $rootScope.frames;
    var data;
    var firstSort = true;
    $scope.$watch('tableParams', function(params) {
        var firstRun = true;
        if ($rootScope.frames == null) {
            $rootScope.frames = Frame;
            $rootScope.frames.setOnMessageCB(
                function (m) {
                    console.log(m.data);
                    $scope.$apply(function(){
                    	var oldData = $rootScope.frames;
                    	$rootScope.frames = JSON.parse(m.data);
                    	if(!firstRun) {
                    		angular.forEach(oldData, function(oldData) {
        	      				$rootScope.frames.push(oldData);
        	      			});
                    	}
        	      		firstRun = false;
                        cachedData = $rootScope.frames;
                    })
                });
        } else {
            if (!firstSort) {
                $rootScope.frames = [];
                angular.forEach(cachedData, function(cachedData) {
                                $rootScope.frames.push(cachedData);
                            });
            }
            firstSort = false;

            data = params.filter ? $filter('filter')($rootScope.frames, params.filter) : $rootScope.frames;
            data = params.sorting ? $filter('orderBy')(data, params.orderBy()) : data;

            var total = data.length;
            $rootScope.frames = data.slice((params.page - 1) * params.count, params.page * params.count);
        }
        $scope.tableParams.total = $rootScope.frames.length;

    }, false);

    // Replay a specific frame
    $scope.SendFrame = function(replayedFrame) {
        $http({
            url: "/send", 
            method: "GET",
            params: {frame: replayedFrame}
        });
    };

    // Send to CMD = Store as rootScope and change to template
    $scope.SendToCMD = function(frameValue) {
        $rootScope.payload = frameValue;
        $location.path("/cmd"); 
    };
}

function CmdCtrl($scope, $http) {
    $scope.sendSinglePayload = function() {
        $http({
            url: "/send", 
            method: "GET",
            params: {frame: $scope.payload}
        });
    };
    $scope.repeatedPayloadSend = function() {
        $http({
            url: "/send", 
            method: "GET",
            params: {frame: $scope.payload, number: $scope.payloadRepeatNo}
        });
    };
    $scope.repeatedPayloadSend = function() {
        $http({
            url: "/send", 
            method: "GET",
            params: {frame: $scope.payload, frame2: $scope.payload2, number: $scope.payloadRepeatNo}
        });
    };
    $scope.partiallyEncrypted = function() {
        $http({
            url: "/sendPartiallyHandler", 
            method: "GET",
            params: {frame: $scope.payload, appendedData: $scope.payload2}
        });
    };
}


function SettingsCtrl($scope, $http) {

}