document
  .getElementById("login-form")
  .addEventListener("submit", function (event) {
    event.preventDefault();

    var username = document.getElementById("login-username").value;
    var password = document.getElementById("login-password").value;

    fetch("/login", {
      method: "POST",
      body: new URLSearchParams({
        username: username,
        password: password,
      }),
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
    }).then(function (response) {
      if (response.status !== 200) {
        alert("Invalid username or password");
        return;
      }

      response.text().then(function (token) {
        sessionStorage.setItem("token", token);
        sessionStorage.setItem("username", username); // Save username in session storage
        window.location.href = "createOrEnterChatroom.html";
      });
    });
  });
