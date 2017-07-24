var app = angular.module("terraboard", [], function($locationProvider){
    $locationProvider.html5Mode(true);
});

app.controller("tbBreadCtrl", ['$scope', '$location', function($scope, $location) {
    $scope.$on('$locationChangeSuccess', function() {
        $scope.path = $location.path().replace('/state/', '');
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
        $http.get('api'+$location.url(), {cache: true}).then(function(response){
            $scope.path = $location.path();
            $scope.details = response.data;
            var mods = $scope.details.modules;

            // Init
            if ($location.hash() != "") {
                // Default
                $scope.selectedmod = 0;

                // Search for module in selected res
                var targetRes = $location.hash();
                for (i=0; i < mods.length; i++) {
                    if (targetRes.startsWith(mods[i].path+'.')) {
                        $scope.selectedmod = i;
                    }
                }

                targetRes = targetRes.replace(mods[$scope.selectedmod].path+'.', '');
                var resources = mods[$scope.selectedmod].resources;
                for (j=0; j < resources.length; j++) {
                    if (targetRes == resources[j].type+'.'+resources[j].name) {
                        $scope.selectedres = j;
                        break;
                    }
                }

                // Init display.mod
                $scope.display.mod = $scope.selectedmod;
            }

            $scope.setSelected = function(m, r) {
                var mod = $scope.details.modules[m];
                var res = mod.resources[r];
                var res_title = res.type+'.'+res.name;
                var hash = (mod == 0) ? res_title : mod.path+'.'+res_title;
                $location.hash(hash);
            };
        });
    });
}]);

