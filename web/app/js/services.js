'use strict';

/* Services */


// Demonstrate how to register services
// In this case it is a simple value service.
var employeeService = angular.module('plannerApp.services', ['ngResource']).value('version', '0.1');


employeeService.factory('Employees', ['$resource',
  function($resource){
    return $resource('/data/employees', {}, {
      query: {method:'GET', params:{}, isArray:false}
    });
  }]);

employeeService.factory('Employee', ['$resource',
  function($resource){
    return $resource('/data/employee/101', {}, {
      query: {method:'GET', params:{}, isArray:false}
    });
  }]);