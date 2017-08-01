var app = angular.module("terraboard", ['ngRoute', 'ngSanitize', 'ui.select'], function($locationProvider, $routeProvider){
    $locationProvider.html5Mode(true);

    $routeProvider.when("/", {
        templateUrl: "static/main.html",
        controller: "tbMainCtrl"
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

app.directive("sparklinechart", function () {
	return {
		restrict: "E",
		scope: {
			data: "@"
		},
		compile: function (tElement, tAttrs, transclude) {
			tElement.replaceWith("<span>" + tAttrs.data + "</span>");
			return function (scope, element, attrs) {
				attrs.$observe("data", function (newValue) {
					element.html(newValue);
					element.sparkline(
							'html', {
								type: 'line',
								width: '200px',
								height: 'auto',
								barWidth: 11,
								barColor: 'blue'
							});
				});
			};
		}
	};
});

app.controller("tbMainCtrl", ['$scope', '$http', function($scope, $http) {
    $scope.itemsPerPage = 20;

    $scope.getStats = function(page) {
        var params = {};
        if (page != undefined) {
            params.page = page;
        }
        var query = $.param(params);
        $http.get('api/states/stats?'+query).then(function(response){
            $scope.results = response.data;
            $scope.pages = Math.ceil($scope.results.total / $scope.itemsPerPage);
            $scope.page = $scope.results.page;
            $scope.prevPage = (page <= 1) ? undefined : $scope.page - 1;
            $scope.nextPage = (page >= $scope.pages) ? undefined : $scope.page + 1;
            $scope.startItems = $scope.itemsPerPage*($scope.page-1)+1;
            $scope.itemsInPage = Math.min($scope.itemsPerPage*$scope.page, $scope.results.total)
        });
    };

    // On page load
    $scope.getStats(1);

    $scope.getActivity = function(idx, path) {
        $http.get('api/state/activity/'+path).then(function(response){
            var states = response.data;
            var activityData = [];
            for (i=0; i < states.length; i++) {
                var date = new Date(states[i].last_modified).getTime() / 1000;
                activityData.push(date+":"+states[i].resource_count);
            }
            var activity = activityData.join(",");

            $scope.results.states[idx].activity = activity;
        });
    };

    $http.get('api/locks').then(function(response){
        $scope.locks = response.data;

        $scope.isLocked = function(path) {
            if (path in $scope.locks) {
                return true;
            }
            return false;
        };
    });
}]);

app.controller("tbListCtrl", ['$scope', '$http', '$location', function($scope, $http, $location) {
    if ($location.path().startsWith("/state/")) {
        $scope.placeholder = $location.path().replace('/state/', '');
    } else {
        $scope.placeholder = 'Enter a state file path...';
    }
    $http.get('api/states').then(function(response){
        $scope.states = response.data;
    });
}]);

app.controller("tbStateCtrl", ['$scope', '$http', '$location', function($scope, $http, $location) {
    $scope.Utils = { keys : Object.keys };
    $scope.display = {};

    // Init
    $scope.selectedVersion = $location.search().versionid;

    var key = $location.url().replace('/state/', '');
    $http.get('api/state/activity/'+key).then(function(response){
        $scope.versions = {};
        for (i=0; i<response.data.length; i++) {
            $scope.versions[response.data[i].version_id] = new Date(response.data[i].last_modified).toLocaleString();
        }
        $scope.$watch('selectedVersion', function(ver) {
            $location.search('versionid', ver);
        });
    });

    $http.get('api'+$location.url(), {cache: true}).then(function(response){
        $scope.path = $location.path().replace('/state/', '');
        $scope.details = response.data;
        $scope.selectedVersion = $scope.details.version.version_id;
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

    $http.get('api/locks').then(function(response){
        $scope.locks = response.data;

        $scope.isLocked = function(path) {
            if ($scope.path in $scope.locks) {
                return true;
            }
            return false;
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

    $scope.refreshAttrKeys = function() {
        console.log("Refreshing keys");
        $http.get('api/attribute/keys?resource_type='+$scope.resType).then(function(response){
            console.log(response.data);
            $scope.attribute_keys = response.data;
        });
    };

    $scope.itemsPerPage = 20;

    $scope.doSearch = function(page) {
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
        if (page != undefined) {
            params.page = page;
        }
        var query = $.param(params);
        $http.get('api/search/attribute?'+query).then(function(response){
            $scope.results = response.data;
            $scope.pages = Math.ceil($scope.results.total / $scope.itemsPerPage);
            $scope.page = $scope.results.page;
            $scope.prevPage = (page <= 1) ? undefined : $scope.page - 1;
            $scope.nextPage = (page >= $scope.pages) ? undefined : $scope.page + 1;
            $scope.startItems = $scope.itemsPerPage*($scope.page-1)+1;
            $scope.itemsInPage = Math.min($scope.itemsPerPage*$scope.page, $scope.results.total)
        });
    }

    // On page load
    $scope.doSearch(1);

    $scope.clearForm = function() {
        $scope.resType = undefined;
        $scope.resID = undefined;
        $scope.attrKey = undefined;
        $scope.attrVal = undefined;
        $scope.results = undefined;
        $scope.doSearch(1);
    }
}]);
