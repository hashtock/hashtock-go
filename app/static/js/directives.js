'use strict';

/* Directives */

var hashtockDirectives = angular.module('hashtockDirectives', []);

hashtockDirectives.directive('ionslider',function(){
    return{
        restrict:'E',
        scope:{
            min:'=',
            max:'=',
            step:'@',
            action:'&'
        },
        controller: function($rootScope, $scope, $element){
            (function init(){
                $($element).ionRangeSlider({
                    type: "single",
                    grid: false,
                    min: $scope.min,
                    max: $scope.max,
                    step: $scope.step,
                    onChange: function(data) {
                        $scope.action({'data':data.from});
                    }
                });
            })();

            var slider = $($element).data("ionRangeSlider");

            $scope.$watch('max', function(value) {
                var disabled = $scope.max === $scope.min;
                if (value !== undefined) {
                    slider.update({
                        max: value,
                        step: $scope.step,
                        disable: disabled
                    });
                }
            }, true);       

            $scope.$watch('min', function(value) {
                var disabled = $scope.max === $scope.min;
                if (value !== undefined) {
                    slider.update({
                        min: value,
                        step: $scope.step,
                        disable: disabled
                    });
                }
            }, true);
        }
    }
});
