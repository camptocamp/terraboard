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
            var mods = $scope.details.modules;

            // Init
            if ($location.hash() != "") {
                // Default
                $scope.selectedmod = 0;
                $scope.selectedres = $location.hash();

                // Search for module in selected res
                for (i=0; i < mods.length; i++) {
                    if ($scope.selectedres.startsWith(mods[i].path[1]+'.')) {
                        $scope.selectedmod = i;
                        $scope.selectedres = $scope.selectedres.replace(mods[i].path[1]+'.', '');
                        break;
                    }
                }
            }

            $scope.setSelected = function(mod, res) {
                var hash = (mod == 0) ? res : $scope.details.modules[mod].path[1]+'.'+res;
                $location.hash(hash);
            };
        });
    });
}]);

