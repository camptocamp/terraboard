var app = angular.module("terraboard", [], function($locationProvider){
    $locationProvider.html5Mode(true);
});

app.controller("tbBreadCtrl", ['$scope', '$location', function($scope, $location) {
    $scope.$on('$locationChangeSuccess', function() {
        $scope.path = $location.path();
    });
}]);

app.controller("tbListCtrl", ['$scope', '$http', '$location', function($scope, $http, $location) {
    $http.get('api/states').then(function(response){
        $scope.keys = response.data;
    });
}]);

app.controller("tbStateCtrl", ['$scope', '$http', '$location', function($scope, $http, $location) {
    $scope.Utils = { keys : Object.keys };
    $scope.display = {};

    // Init
    $scope.selectedVersion = $location.search().versionid;

    var key = $location.url().replace('/state/', '');
    $http.get('api/history/'+key).then(function(response){
        $scope.history = response.data;
        $scope.versions = {};
        for (i=0; i<response.data.length; i++) {
            $scope.versions[response.data[i].VersionId] = new Date(response.data[i].LastModified).toLocaleString();
        }
        $scope.$watch('selectedVersion', function(ver) {
            $location.search('versionid', ver);
        });
    });

    $scope.$on('$locationChangeSuccess', function() {
        $http.get('api'+$location.url()).then(function(response){
            $scope.path = $location.path();
            $scope.details = response.data;
            $scope.setSelected = function(mod, res) {
                $scope.selectedmod = mod;
                $scope.selectedres = res;
            };
        });
    });
}]);

