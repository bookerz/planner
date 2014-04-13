'use strict';

/* Controllers */

angular.module('plannerApp.controllers', []).
  controller('MyCtrl1', [function() {

  }])
  .controller('MyCtrl2', [function() {

  }])
  .controller('EmployeeCtrl', ['$scope', 'Employees', function($scope, Employees) {
  	$scope.employees = Employees.query();
  }]);