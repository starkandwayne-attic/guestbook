app.controller('HomeCtl', function($scope, $http) {
  $scope.submitted = "";
  $scope.post_id = -1;
  $scope.post_url = "";
  $scope.post_title = "";
  $scope.submit_code = function() {
    $scope.error = {};
    $scope.error.status = false;
    $scope.success = {};
    $scope.success.status = false;
    submittal = {};
    submittal.code = $scope.code;
    submittal.name = $scope.name;
    submittal.email = $scope.email;
    submittal.comment = $scope.comment;
    submittal.post_id = $scope.post_id;
    submit_data = {}
    submit_data.submit = submittal
    $http.post('/submit', submit_data).success(function(data) {
      $scope.submitted = "disabled";
      $scope.success.status = true;
      $scope.success.message = data.success
    }).
    error(function(data, status) {
      $scope.error.status = true;
      $scope.error.message = data.error;
    });
  };
  $scope.load_random_post = function() {
    $http.get('/posts/random').success(function(data) {
      $scope.post_id = data.id;
      $scope.post_url = data.url;
      $scope.post_title = data.title;
    });
  };
  $scope.load_random_post();
})
