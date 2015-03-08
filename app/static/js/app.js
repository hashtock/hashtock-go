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

hashtockApp.run(function ($rootScope, $location, User) {
    $rootScope.$on('$routeChangeStart', function () {
        if ($rootScope.loggedIn) {
            return;
        }

        User.get(function(user) {
            if (user.id === undefined) {
                window.location = "/auth/login/?continue=" + encodeURIComponent($location.absUrl());
            } else {
                $rootScope.loggedIn = true;
            }
        }, function() {
            window.location = "/auth/login/?continue=" + encodeURIComponent($location.absUrl());
        });
    });
});

hashtockApp.config(['$resourceProvider', function($resourceProvider) {
    // Don't strip trailing slashes from calculated URLs
    $resourceProvider.defaults.stripTrailingSlashes = false;
}]);
