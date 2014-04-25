'use strict';


// Declare app level module which depends on filters, and services
angular.module('plannerApp', [
  'ngRoute',
  'plannerApp.filters',
  'plannerApp.services',
  'plannerApp.directives',
  'plannerApp.controllers'
]).
config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/view1', {templateUrl: 'partials/partial1.html', controller: 'MyCtrl1'});
  $routeProvider.when('/view2', {templateUrl: 'partials/partial2.html', controller: 'MyCtrl2'});
  $routeProvider.when('/employee', {templateUrl: 'partials/employee.html', controller: 'EmployeeCtrl'});
  $routeProvider.when('/employees', {templateUrl: 'partials/employees.html', controller: 'EmployeeListCtrl'});
  $routeProvider.otherwise({redirectTo: '/view1'});
}]);
