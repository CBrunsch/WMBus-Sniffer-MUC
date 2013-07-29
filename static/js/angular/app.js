var app = angular.module('app', ['ngTable', '$strap.directives']).
  config(['$routeProvider', function($routeProvider) {
  $routeProvider.
      when('/', {templateUrl: 'templates/partials/sniffer.html',  controller: MainCtrl}).
      when('/cmd', {templateUrl: 'templates/partials/cmd.html',  controller: CmdCtrl}).
      otherwise({redirectTo: '/'});
}]);

app.run(function($rootScope, $http) {
	$rootScope.snifferActive = snifferActive;
	$rootScope.stopSniffer = function() {
       $rootScope.snifferActive = false;
       $http({
            url: "/statusSniffer", 
            method: "GET",
            params: {status: "stop"}
        });
    };
   	$rootScope.startSniffer = function() {
       $rootScope.snifferActive = true;
        $http({
            url: "/statusSniffer", 
            method: "GET",
            params: {status: "start"}
        });
    };
    $rootScope.export = function() {
        window.open("/export");
    };
    $rootScope.truncate = function() {
       $http({
            url: "/truncate", 
            method: "GET",
        });
       $rootScope.frames = null;
     };
});