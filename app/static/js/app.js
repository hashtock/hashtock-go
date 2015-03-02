'use strict';

/* App Module */

var hashtockApp = angular.module('hashtockApp', [
    'ngRoute',

    'hashtockControllers',
    'hashtockServices'
]);

hashtockApp.config(['$routeProvider',
    function($routeProvider) {
        $routeProvider.
        when('/tags', {
            templateUrl: '/static/partials/tag-table.html',
            controller: 'TagListCtrl'
        }).
        otherwise({
            redirectTo: '/tags'
        });
    }
]);

hashtockApp.config(['$resourceProvider', function($resourceProvider) {
  // Don't strip trailing slashes from calculated URLs
  $resourceProvider.defaults.stripTrailingSlashes = false;
}]);
