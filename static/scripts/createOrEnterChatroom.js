// Check if user is logged in
if (!localStorage.getItem("token")) {
  // If not, redirect to login page
  window.location.href = "login.html";
}

$(document).ready(function () {
  $("#create-chatroom-button").click(function () {
    var chatroomName = $("#chatroom-input").val();
    $.post("/createChatroom", { chatroomName: chatroomName })
      .done(function (chatroomId) {
        alert("Chatroom created with ID: " + chatroomId);

        // Save chatroom ID in local storage
        localStorage.setItem("chatroomID", chatroomId);

        // Redirect to chatroom.html
        window.location.href = "chatroom.html";
      })
      .fail(function () {
        alert("Error creating chatroom");
      });
  });
});

$("#enter-chatroom-button").click(function () {
  var chatroomID = $("#chatroom-input").val();
  if (chatroomID) {
    // Save chatroom ID in local storage
    localStorage.setItem("chatroomID", chatroomID);

    // Redirect to chatroom.html
    window.location.href = "chatroom.html";
  } else {
    alert("Please enter a chatroom ID");
  }
});
