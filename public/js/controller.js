app.controller('HomeCtl', function($scope, $http) {
  $scope.submitted = "";
  $scope.submit_code = function() {
    $scope.error = {};
    $scope.error.status = false;
    submittal = {};
    submittal.code = $scope.code;
    submittal.name = $scope.name;
    submittal.email = $scope.email;
    submittal.comment = $scope.comment;
    submittal.post_id = 2;
    submit_data = {}
    submit_data.submit = submittal
    $http.post('/submit', submit_data).success(function(data) {
      $scope.submitted = "disabled";
      $scope.success = true;
    }).
    error(function(data, status) {
      $scope.error.status = true;
      $scope.error.message = data.error;
    });
  }
})
