var app = angular.module("terraboard", ['ngRoute', 'ngSanitize', 'ui.select', 'chart.js'], function($locationProvider, $routeProvider){
    $locationProvider.html5Mode(true);

    $routeProvider.when("/", {
        templateUrl: "static/main.html",
        controller: "tbMainCtrl"
    }).when("/state/:path*", {
        templateUrl: "static/state.html",
        controller: "tbStateCtrl",
        reloadOnSearch: false
    }).when("/search", {
        templateUrl: "static/search.html",
        controller: "tbSearchCtrl",
        reloadOnSearch: false
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

app.directive("hlcode", ['$timeout', function($timeout) {
    return {
        restrict: "E",
        scope: {
            code: '=code',
            lang: '=lang'
        },
        link: function() {
            $timeout(sh_highlightDocument, 0, false);
        },
        template: "<pre class=\"sh_{{lang}} sh_sourceCode\">{{code}}</pre>"
    }
}]);

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
    $scope.searchType = function(points, ev) {
        var type = points[0]._chart.data.labels[points[0]._index];
        if ($.isArray(type)) {
            console.log("Clicked zone is an array, not searching");
            return;
        }
        $location.url('search/?type='+type);
        $scope.$apply();
    };



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
    $scope.searchTfVersion = function(points, ev) {
        var version = points[0]._chart.data.labels[points[0]._index][0];
        if ($.isArray(version)) {
            console.log("Clicked zone is an array, not searching");
            return;
        }
        $location.url('search/?tf_version='+version);
        $scope.$apply();
    };


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

app.controller("tbListCtrl",
        ['$scope', '$http', '$location', '$routeParams',
        function($scope, $http, $location, $routeParams) {
    $scope.placeholder = 'Enter a state file path...';

    $scope.$on('$routeChangeSuccess', function() {
        if ($routeParams.path != undefined) {
            $scope.placeholder = $routeParams.path;
        }
    });

    $http.get('api/states').then(function(response){
        $scope.states = response.data;
    });
}]);

app.controller("tbStateCtrl",
        ['$scope', '$http', '$location', '$routeParams',
        function($scope, $http, $location, $routeParams) {
    $scope.Utils = { keys : Object.keys };
    $scope.display = {
      welcome: true,
      details: false,
      compare: false
    };

    // Init
    $scope.selectedVersion = {
        versionId: $location.search().versionid
    };

    $scope.compareVersion = {
        versionId: $location.search().compare
    };

    $scope.$on('$routeChangeSuccess', function() {
        $http.get('api/state/activity/'+$routeParams.path).then(function(response){
            $scope.versions = [];
            for (i=0; i<response.data.length; i++) {
                var ver = {
                    versionId: response.data[i].version_id,
                    date: new Date(response.data[i].last_modified.toLocaleString())
                };
                $scope.versions.unshift(ver);
            }

            $scope.$watch('compareVersion', function(ver) {
                if (ver != undefined && ver.versionId != undefined) {
                    $location.search('compare', ver.versionId);
                    $scope.display.welcome = false;
                    $scope.display.details = false;
                    $scope.display.compare = true;
                    $http.get('api/state/compare/'+$routeParams.path+'?from='+$scope.selectedVersion.versionId+'&to='+ver.versionId).then(function(response){
                        $scope.compare = response.data;

                        $scope.only_in_old = Object.keys($scope.compare.differences.only_in_old).length;
                        $scope.only_in_new = Object.keys($scope.compare.differences.only_in_new).length;
                        $scope.differences = Object.keys($scope.compare.differences.resource_diff).length;
                    });
                } else {
                    $location.search('compare', null);
                    $scope.display.compare = false;
                    $scope.display.details = true;
                }
            });
        });
    });

    $scope.getDetails = function(versionId) {
        $http.get('api/state/'+$routeParams.path+'?versionid='+versionId+'#'+$location.hash()).then(function(response){
            $scope.path = $routeParams.path;
            $scope.details = response.data;
            var mods = $scope.details.modules;

            // Init
            if ($location.hash() != "") {
                // Default
                $scope.display.welcome = false;
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
                $scope.selectedmod = m;
                $scope.selectedres = r;
                var mod = $scope.details.modules[m];
                var res = mod.resources[r];
                var res_title = res.type+'.'+res.name;
                var hash = (mod == 0) ? res_title : mod.path+'.'+res_title;
                $location.hash(hash);
            };
        });
    };

    $scope.$on('$routeChangeSuccess', function() {
        $scope.getDetails($location.search().version_id);
    });

    $scope.$watch('selectedVersion', function(ver) {
        $scope.getDetails(ver.versionId);
        $location.search('versionid', ver.versionId);
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

app.controller("tbSearchCtrl",
        ['$scope', '$http', '$location', '$routeParams',
        function($scope, $http, $location) {
    $http.get('api/tf_versions').then(function(response){
        $scope.tf_versions = response.data;
    });
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
        $http.get('api/attribute/keys?resource_type='+$scope.resType).then(function(response){
            $scope.attribute_keys = response.data;
        });
    };

    $scope.itemsPerPage = 20;

    $scope.doSearch = function(page) {
        var params = {};
        if ($scope.resType != undefined) {
            params.type = $scope.resType;
        }
        if ($scope.resID != undefined) {
            params.name = $scope.resID;
        }
        if ($scope.attrKey != undefined) {
            params.key = $scope.attrKey;
        }
        if ($scope.attrVal != undefined) {
            params.value = $scope.attrVal;
        }
        if ($scope.tfVersion != undefined) {
            params.tf_version = $scope.tfVersion;
        }
        if (page != undefined) {
            params.page = page;
        }
        var query = $.param(params);
        $location.path($location.path()).search(params);
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
    if ($location.search().type != undefined) {
        $scope.resType = $location.search().type;
    }
    if ($location.search().name != undefined) {
        $scope.resID = $location.search().name;
    }
    if ($location.search().key != undefined) {
        $scope.attrKey = $location.search().key;
    }
    if ($location.search().value != undefined) {
        $scope.attrVal = $location.search().value;
    }
    if ($location.search().tf_version != undefined) {
        $scope.tfVersion = $location.search().tf_version;
    }

    $scope.doSearch(1);

    $scope.clearForm = function() {
        $scope.tfVersion = undefined;
        $scope.resType = undefined;
        $scope.resID = undefined;
        $scope.attrKey = undefined;
        $scope.attrVal = undefined;
        $scope.results = undefined;
        $location.url($location.path());
    }
}]);
