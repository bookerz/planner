'use strict';

/* Controllers */

angular.module('plannerApp.controllers', []).
  controller('CustomersListCtrl', ['$scope', 'getCustomers', function($scope, getCustomers) {
   $scope.customers = getCustomers.query();
  }])
  .controller('MyCtrl2', [function() {
  }])
  .controller('EmployeeListCtrl', ['$scope', 'getEmployees', function($scope, getEmployees) {
  	$scope.employees = getEmployees.query();
  }])
  .controller('EmployeeCtrl', ['$scope', 'getEmployee', function($scope, getEmployee) {
  	$scope.employee = getEmployee.query();
  }]);