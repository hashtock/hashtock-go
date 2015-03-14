'use strict';

/* Services */

var hashtockServices = angular.module('hashtockServices', ['ngResource']);

hashtockServices.factory('Tag', ['$resource', function($resource){
    return $resource('/api/tag/:tag/', {}, {
        query: {method:'GET', params:{tag:'@hashtag'}, isArray:true}
    });
}]);

hashtockServices.factory('TagValues', ['$resource', function($resource){
    return $resource('/api/tag/:tag/values/', {}, {
        query: {method:'GET', params:{tag:'@hashtag', days:undefined}, isArray:true}
    });
}]);

hashtockServices.factory('User', ['$resource', function($resource){
    return $resource('/api/user/', {}, {
        get: {method:'GET', isArray:false}
    });
}]);

hashtockServices.factory('Order', ['$resource', function($resource){
    return $resource('/api/order/:uuid/', {}, {
        pending: {method:'GET', isArray:true},
        history: {method: 'GET', params:{uuid:'history'}, isArray:true},
        get: {method: 'GET', params:{uuid:'@uuid'}, isArray:false},
        save: {method: 'POST', params:{uuid:'@uuid'}, isArray:false},
        cancel: {method: 'DELETE', params:{uuid:'@uuid'}, isArray:false}
    });
}]);

hashtockServices.factory('Portfolio', ['$resource', function($resource){
    return $resource('/api/portfolio/:tag/', {}, {
        query: {method:'GET', isArray:true},
        tag: {method:'GET', params:{tag:'@hashtag'}, isArray:false}
    });
}]);
