app.controller('HomeCtl', function($scope, $http) {
  $scope.submitted = "";
  $scope.post_id = -1;
  $scope.post_url = "";
  $scope.post_title = "";
  $scope.submit_email = function() {
    $scope.error = {};
    $scope.error.status = false;
    $scope.success = {};
    $scope.success.status = false;
    submittal = {};
    submittal.name = $scope.name;
    submittal.email = $scope.email;
    submittal.comment = $scope.comment;
    submit_data = {}
    submit_data.submit = submittal
    $http.post('/submit/email', submit_data).success(function(data) {
      $scope.submitted = "disabled";
      $scope.success.status = true;
      $scope.success.message = data.success
      $scope.post_id = data.id;
      $scope.post_url = data.url;
      $scope.post_title = data.title;
    }).
    error(function(data, status) {
      $scope.error.status = true;
      $scope.error.message = data.error;
    });
  };
  $scope.submit_code = function() {
    $scope.error = {};
    $scope.error.status = false;
    $scope.success = {};
    $scope.success.status = false;
    submittal = {};
    submittal.code = $scope.code;
    submittal.email = $scope.email;
    submit_data = {}
    submit_data.submit = submittal
    $http.post('/submit/code', submit_data).success(function(data) {
      $scope.submitted = "disabled";
      $scope.success.status = true;
      $scope.success.message = data.success
    }).
    error(function(data, status) {
      $scope.error.status = true;
      $scope.error.message = data.error;
    });
  }
})
