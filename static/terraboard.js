var app = angular.module("terraboard", ['ngRoute', 'ngSanitize', 'ui.select'], function($locationProvider, $routeProvider){
    $locationProvider.html5Mode(true);

    $routeProvider.when("/", {
        templateUrl: "static/main.html"
    }).when("/state/:path*", {
        templateUrl: "static/state.html",
        controller: "tbStateCtrl"
    }).when("/search", {
        templateUrl: "static/search.html",
        controller: "tbSearchCtrl"
    }).otherwise({
        redirectTo: "/"
    });
});

app.controller("tbListCtrl", ['$scope', '$http', '$location', function($scope, $http, $location) {
    if ($location.path().startsWith("/state/")) {
        $scope.placeholder = $location.path().replace('/state/', '');
    } else {
        $scope.placeholder = 'Enter a state file path...';
    }
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

    $http.get('api'+$location.url(), {cache: true}).then(function(response){
        $scope.path = $location.path();
        $scope.details = response.data;
        $scope.selectedVersion = $scope.details.version_id;
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
}]);

app.controller("tbSearchCtrl", ['$scope', '$http', '$location', '$routeParams', function($scope, $http) {
    $http.get('api/resource/types').then(function(response){
        $scope.resource_keys = response.data;
    });
    $http.get('api/resource/names').then(function(response){
        $scope.resource_names = response.data;
    });
    $http.get('api/attribute/keys').then(function(response){
        $scope.attribute_keys = response.data;
    });

    $scope.doSearch= function() {
        var params = {};
        if ($scope.resType != "") {
            params.type = $scope.resType;
        }
        if ($scope.resID != "") {
            params.name = $scope.resID;
        }
        if ($scope.attrKey != "") {
            params.key = $scope.attrKey;
        }
        if ($scope.attrVal != "") {
            params.value = $scope.attrVal;
        }
        var query = $.param(params);
        console.log(query);
        $http.get('api/search/attribute?'+query).then(function(response){
            $scope.results = response.data;
        });
    }
}]);
