'use strict';

/* Controllers */

angular.module('plannerApp.controllers', []).
  controller('MyCtrl1', [function() {

  }])
  .controller('MyCtrl2', [function() {

  }])
  .controller('EmployeeListCtrl', ['$scope', 'Employees', function($scope, Employees) {
  	$scope.employees = Employees.query();
  }])
  .controller('EmployeeCtrl', ['$scope', 'Employee', function($scope, Employee) {
  	$scope.employee = Employee.query();
  }]);