'use strict';

/* Controllers */

var hashtockControllers = angular.module('hashtockControllers', []);

hashtockControllers.controller('TagListCtrl', ['$scope', 'Tag',
    function($scope, Tag) {
        $scope.tags = Tag.query();
        $scope.orderProp = 'hashtag';
        $scope.reverse = false;
    }
]);

hashtockControllers.controller('UserCtrl', ['$scope', 'User',
    function($scope, User) {
        $scope.user = User.query();
    }
]);
