app.controller('HomeCtl', function($scope, $http) {
  $scope.email_form = {}
  $scope.email_form.submitted = false;
  $scope.code_form = {}
  $scope.code_form.submitted = false;
  $scope.post_id = -1;
  $scope.post_url = "";
  $scope.post_title = "";
  $scope.set_message_after_email_submit = true;
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
      $scope.email_form.submitted = true;
      if ($scope.set_message_after_email_submit) {
        $scope.success.status = true;
        $scope.success.message = data.success;
        $scope.set_message_after_email_submit = false;
      }
      $scope.post_id = data.id;
      $scope.post_url = data.url;
      $scope.entry_id = data.entry_id;
      $scope.post_title = "READ THIS POST:  ".concat(data.title);
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
    submittal.entry_id = $scope.entry_id;
    submittal.post_id = $scope.post_id;
    submit_data = {}
    submit_data.submit = submittal
    $http.post('/submit/entry', submit_data).success(function(data) {
      $scope.submit_email();
      $scope.success.status = true;
      $scope.success.message = "Correct!  \"".concat($scope.code).concat("\" is the correct code.  ").concat(data.success);
      $scope.code = "";
    }).
    error(function(data, status) {
      $scope.error.status = true;
      $scope.error.message = data.error;
    });
  }
})
