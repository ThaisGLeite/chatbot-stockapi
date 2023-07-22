// Check if user is logged in
if (!sessionStorage.getItem("token")) {
  // If not, redirect to login page
  console.log("User not logged in");
  window.location.href = "login.html";
}

$(document).ready(function () {
  $.get("/getAllChatrooms")
    .done(function (data) {
      var chatroomNames = JSON.parse(data);
      chatroomNames.forEach(function (chatroomName) {
        $("#chatroom-list").append("<li>" + chatroomName + "</li>");
      });
    })
    .fail(function () {
      alert("Error getting chatrooms");
    });

  $("#create-chatroom-button").click(function () {
    var chatroomName = $("#chatroom-input").val();
    $.post("/createChatroom", { chatroomName: chatroomName })
      .done(function (chatroomId) {
        // Save chatroom ID in session storage
        sessionStorage.setItem("chatroomID", chatroomId);

        // Redirect to chatroom.html
        window.location.href = "chatroom.html";
      })
      .fail(function () {
        alert("Error creating chatroom");
      });
  });

  $("#enter-chatroom-button").click(function () {
    var chatroomID = $("#chatroom-input").val();
    if (chatroomID) {
      // Save chatroom ID in session storage
      sessionStorage.setItem("chatroomID", chatroomID);

      // Redirect to chatroom.html
      window.location.href = "chatroom.html";
    } else {
      alert("Please enter a chatroom ID");
    }
  });
});
