'use strict';

/* App Module */

var hashtockApp = angular.module('hashtockApp', [
    'ngRoute',

    'nvd3',

    'hashtockControllers',
    'hashtockServices',
    'hashtockDirectives'
]);

hashtockApp.config(['$routeProvider',
    function($routeProvider) {
        $routeProvider.
        when('/portfolio', {
            templateUrl: '/static/partials/portfolio-table.html',
            controller: 'PortfolioCtrl'
        }).
        when('/tags', {
            templateUrl: '/static/partials/tag-table.html',
            controller: 'TagListCtrl'
        }).
        when('/tags/:tag', {
            templateUrl: '/static/partials/tag-details.html',
            controller: 'TagDetailCtrl'
        }).
        when('/orders', {
            templateUrl: '/static/partials/order-table.html',
            controller: 'OrderListCtrl'
        }).
        otherwise({
            redirectTo: '/portfolio'
        });
    }
]);

hashtockApp.config(['$resourceProvider', function($resourceProvider) {
  // Don't strip trailing slashes from calculated URLs
  $resourceProvider.defaults.stripTrailingSlashes = false;
}]);
