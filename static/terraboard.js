var app = angular.module("terraboard", ['ngRoute', 'ngSanitize', 'ui.select', 'chart.js'], function($locationProvider, $routeProvider){
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
                            barColor: 'blue',
                            tooltipFormatter: function(sp, opts, fields) {
                                console.log(fields.x);
                                var date = new Date(0);
                                date.setUTCSeconds(fields.x);
                                return date.toLocaleString()+' - '+fields.y+' resources';
                            }
                        }
                    );
                    element.bind('sparklineClick', function(ev) {
                        var sparkline = ev.sparklines[0],
                        region = sparkline.getCurrentRegionFields();
                        var path = element[0].attributes.path.value;
                        scope.$parent.$parent.goToState(path, region.x);
                    });
                });
            };
        }
    };
});

app.controller("tbMainCtrl", ['$scope', '$http', '$location', function($scope, $http, $location) {
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

    // Version map for sparklines click events
    $scope.versionMap = {};
    $scope.getActivity = function(idx, path) {
        $http.get('api/state/activity/'+path).then(function(response){
            var states = response.data;
            $scope.versionMap[path] = {};
            var activityData = [];
            for (i=0; i < states.length; i++) {
                var date = new Date(states[i].last_modified).getTime() / 1000;
                activityData.push(date+":"+states[i].resource_count);
                $scope.versionMap[path][date] = states[i].version_id;
            }
            var activity = activityData.join(",");

            $scope.results.states[idx].activity = activity;
        });
    };

    $scope.goToState = function(path, epoch) {
        var versionId = $scope.versionMap[path][epoch];
        var url = 'state/'+path+'?versionid='+versionId;
        $location.url(url);
        $scope.$apply();
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

    pieResourceTypesLabels   = [[], [], [], [], [], [], ["Total"]];
    pieResourceTypesData     = [0, 0, 0, 0, 0, 0, 0];
    $http.get('api/resource/types/count').then(function(response){
        data = response.data;
        angular.forEach(data, function(value, i) {
            if(i < 6) {
                pieResourceTypesLabels[i] = value.name;
                pieResourceTypesData[i]   = parseInt(value.count, 10);
            } else {
                pieResourceTypesLabels[6].push(value.name+": "+value.count);
                pieResourceTypesData[6] += parseInt(value.count, 10);
            }
        });
    });
    $scope.pieResourceTypesData    = pieResourceTypesData;
    $scope.pieResourceTypesLabels  = pieResourceTypesLabels;
    $scope.pieResourceTypesOptions = { legend: { display: false } };



    pieTfVersionsLabels   = [[], [], [], [], [], [], ["Total"]];
    pieTfVersionsData     = [0, 0, 0, 0, 0, 0, 0];
    $http.get('api/states/tfversion/count?orderBy=version').then(function(response){
        data = response.data;
        angular.forEach(data, function(value, i) {
            if(i < 6) {
                pieTfVersionsLabels[i] = [value.name];
                pieTfVersionsData[i]   = parseInt(value.count, 10);
            } else {
                pieTfVersionsData[6] += parseInt(value.count, 10);
                pieTfVersionsLabels[6].push(value.name+": "+value.count);
            }
        });
    });

    $scope.pieTfVersionsLabels  = pieTfVersionsLabels;
    $scope.pieTfVersionsData    = pieTfVersionsData;
    $scope.pieTfVersionsOptions = { legend: { display: false } };


    $scope.pieLockedStatesLabels = ["Locked", "Unlocked"];
    $scope.pieLockedStatesData   = [0, 0];
    $scope.$watch('locks', function(nv, ov){
        $scope.pieLockedStatesData[0] = Object.keys(nv).length;
        $scope.pieLockedStatesData[1] -= Object.keys(nv).length;
    });
    $scope.$watch('results.total', function(nv, ov){
        $scope.pieLockedStatesData[1] = nv - $scope.pieLockedStatesData[0];
    });
    $scope.pieLockedStatesOptions = { legend: { display: false } };


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
    $scope.selectedVersion = {
        versionId: $location.search().versionid
    };

    var key = $location.url().replace('/state/', '');
    $http.get('api/state/activity/'+key).then(function(response){
        $scope.versions = [];
        for (i=0; i<response.data.length; i++) {
            var ver = {
                versionId: response.data[i].version_id,
                date: new Date(response.data[i].last_modified.toLocaleString())
            };
            $scope.versions.unshift(ver);
        }

        $scope.$watch('selectedVersion', function(ver) {
            $location.search('versionid', ver.versionId);
        });
    });

    $http.get('api'+$location.url(), {cache: true}).then(function(response){
        $scope.path = $location.path().replace('/state/', '');
        $scope.details = response.data;
        $scope.selectedVersion = {
            versionId: $scope.details.version.version_id
        };
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
