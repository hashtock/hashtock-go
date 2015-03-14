'use strict';

/* Controllers */

var hashtockControllers = angular.module('hashtockControllers', []);

hashtockControllers.controller('UserCtrl', ['$scope', 'User',
    function($scope, User) {
        $scope.user = User.get();
    }
]);

hashtockControllers.controller('TagListCtrl', ['$scope', 'Tag',
    function($scope, Tag) {
        $scope.tags = Tag.query();
        $scope.orderProp = 'hashtag';
        $scope.reverse = false;
    }
]);

hashtockControllers.controller('OrderListCtrl', ['$scope', 'Order',
    function($scope, Order) {
        $scope.historicalOrders = Order.history();
        $scope.pendingOrders = Order.pending();

        $scope.cancelOrder = function(order) {
            order.$cancel(function(value) {
                $scope.pendingOrders = Order.pending(function(){
                    // In case GAE returns shity old data
                    for (var i = 0; i < $scope.pendingOrders.length; i++) {
                        if ($scope.pendingOrders[i].uuid === order.uuid) {
                           $scope.pendingOrders.splice(i, 1);
                           break;
                        }
                    };
                });
            });
        };
    }
]);

hashtockControllers.controller('PortfolioCtrl', ['$scope', 'Portfolio',
    function($scope, Portfolio) {
        $scope.portfolioTags = Portfolio.query();
    }
]);

hashtockControllers.controller('TagDetailCtrl',
    ['$scope', '$routeParams', '$q', 'User', 'Tag', 'Portfolio', 'Order', 'moment',
  function($scope, $routeParams, $q, User, Tag, Portfolio, Order, moment) {
    $scope.user = User.get();
    $scope.share = Portfolio.tag({tag: $routeParams.tag}, function() {
        $scope.maxSharesToSell = $scope.share.quantity;
    });
    $scope.tag = Tag.get({tag: $routeParams.tag});

    $q.all([$scope.user.$promise, $scope.tag.$promise]).then(function(){
        var val = Math.floor(10 * $scope.user.founds / $scope.tag.value)/10;
        $scope.maxSharesToBuy = Math.min(val, $scope.tag.in_bank);
    })

    $scope.shareQuantityChanged = function(data) {
        $scope.newOrder.quantity = data;
        $scope.$apply();
    }

    $scope.canExecuteOrder = function() {
        if ($scope.newOrder === undefined || $scope.newOrder.quantity <= 0) {
            return false
        }

        if ($scope.newOrder.bank_order === true && $scope.newOrder.action === 'buy') {
            return (0 < $scope.newOrder.quantity && $scope.newOrder.quantity <= $scope.maxSharesToBuy)
        }

        if ($scope.newOrder.bank_order === true && $scope.newOrder.action === 'sell') {
            return (0 < $scope.newOrder.quantity && $scope.newOrder.quantity <= $scope.maxSharesToSell)
        }

        // Other cases...
        return false
    }

    $scope.canBuyFromBank = function(action) {
        return $scope.maxSharesToBuy > 0.1;
    }

    $scope.hasSharesToSell = function() {
        return $scope.maxSharesToSell > 0.1;
    }

    $scope.isDealingWithBank = function() {
        return ($scope.newOrder && $scope.newOrder.bank_order);
    }

    $scope.isOrderInProgress = function() {
        return ($scope.newOrder !== undefined);   
    }

    $scope.maxSharesInCurrentOrder = function () {
        var maxValue = 0
        if ($scope.isDealingWithBank()) {
            if ($scope.newOrder.action === 'buy') {
                maxValue = $scope.maxSharesToBuy;
            }

            if ($scope.newOrder.action === 'sell') {
                maxValue = $scope.maxSharesToSell;
            }
        }

        return maxValue;
    }

    $scope.cancelOrder = function() {
        $scope.newOrder = undefined;
    }

    $scope.newBankOrder = function(action) {
        $scope.newOrder = new Order({
            action: action,
            bank_order: true,
            hashtag: $scope.tag.hashtag,
            quantity: 0
        });
    }

    $scope.executeOrder = function() {
        $scope.newOrder.$save(function(data){
            $scope.newOrder = undefined;
        });
    }
}]);

hashtockControllers.controller('TagValuesCtrl',
    ['$scope', '$routeParams', '$q', 'TagValues', 'moment',
  function($scope, $routeParams, $q, TagValues, moment) {

    $scope.durationOptions = [1, 7, 14, 30];
    $scope.showingDays = 1;

    $scope.showDays = function(days) {
        $scope.loadingTagValues = true;
        $scope.tagValues = TagValues.query({tag: $routeParams.tag, days: days}, function(values) {
            for (var i = 0; i < values.length; i++) {
                var date = new Date(values[i].date)
                values[i].label = date;
            };

            $scope.showingDays = days;
            $scope.data = [{
                key: '#' + $scope.tag.hashtag,
                values: values
            }];
            $scope.loadingTagValues = false;
        });
    }
    $scope.showDays($scope.showingDays);

    $scope.formatDate = function(d) {
        switch ($scope.showingDays) {
            case 1: return moment(d).format("HH:mm");
            default: return moment(d).format("D MMM");
        }
    }

    $scope.options = {
        chart: {
            type: 'lineChart',
            height: 250,
            margin : {top: 20, right: 20, bottom: 60, left: 60},
            x: function(d){ return d.label; },
            y: function(d){ return d.value; },
            showLegend: false,
            xAxis: {
                axisLabel: 'Time',
                showMaxMin: false,
                tickFormat: $scope.formatDate,
            },
            yAxis: {
                axisLabel: 'Value',
                tickFormat: function(d){
                    return d.toFixed(1);
                }
            }
        }
    }
}]);
