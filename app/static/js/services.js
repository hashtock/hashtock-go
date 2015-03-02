'use strict';

/* Services */

var hashtockServices = angular.module('hashtockServices', ['ngResource']);

hashtockServices.factory('Tag', ['$resource', function($resource){
    return $resource('/api/tag/:tag/', {}, {
      query: {method:'GET', params:{tag:'@hashtag'}, isArray:true}
    });
}]);

hashtockServices.factory('User', ['$resource', function($resource){
    return $resource('/api/user/', {}, {
      query: {method:'GET', isArray:false}
    });
}]);

