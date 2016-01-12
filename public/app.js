$(document).ready(function() {
    var lock = new Auth0Lock(
      // All these properties are set in auth0-variables.js
      AUTH0_CLIENT_ID,
      AUTH0_DOMAIN
    );

    var userProfile;

    $('.btn-login').click(function(e) {
      e.preventDefault();
      lock.showSignin(function(err, profile, token) {
        if (err) {
          // Error callback
          console.log("There was an error");
          alert("There was an error logging in");
        } else {
          // Success calback

          // Save the JWT token.
          sessionStorage.setItem('userToken', token);

          // Save the profile
          userProfile = profile;

          $('.login-box').hide();
          $('.logged-in-box').show();
          $('.nickname').text(profile.name);
        }
      });
    });

    $.ajaxSetup({
      'beforeSend': function(xhr) {
        if (sessionStorage.getItem('userToken')) {
          xhr.setRequestHeader('Authorization',
                'Bearer ' + sessionStorage.getItem('userToken'));
        }
      }
    });

    $('.btn-api').click(function(e) {
      // Just call your API here. The header will be sent
      $.ajax({
        url: 'http://localhost:3001/api/products',
        method: 'GET'
      }).then(function(data, textStatus, jqXHR) {
        alert("ça marche");
      }, function() {
        alert("You need to download the server seed and start it to call this API");
      });
    });


});
