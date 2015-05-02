app.controller('HomeCtl', function($scope, $http) {
  $scope.submitted = "";
  $scope.deploy_cf = function() {
    $scope.error = {};
    $scope.error.status = false;
    deploy = {};
    deploy.aws_access_key = $scope.aws_access_key;
    deploy.aws_secret_key = $scope.aws_secret_key;
    deploy.email = $scope.email;
    $http.post('/deploy', deploy).success(function(data) {
      $scope.submitted = "disabled";
      $scope.success = true;
    }).
    error(function(data, status) {
      $scope.error.status = true;
      $scope.error.message = data.error;
    });
  }
})
