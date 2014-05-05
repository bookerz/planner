'use strict';

/* Services */


// Demonstrate how to register services
// In this case it is a simple value service.
var services = angular.module('plannerApp.services', ['ngResource']).value('version', '0.1');

services.factory('getCustomers', ['$resource',
  function($resource){
    return $resource('/data/customers', {}, {
      query: {method:'GET', params:{}, isArray:false}
    });
  }]);

services.factory('getEmployees', ['$resource',
  function($resource){
    return $resource('/data/employees', {}, {
      query: {method:'GET', params:{}, isArray:false}
    });
  }]);

services.factory('getEmployee', ['$resource',
  function($resource){
    return $resource('/data/employee/101', {}, {
      query: {method:'GET', params:{}, isArray:false}
    });
  }]);