// Check if user is logged in
if (!localStorage.getItem("token")) {
  // If not, redirect to login page
  window.location.href = "login.html";
}

$("#send-button").click(function () {
  var message = $("#message-input").val();
  var chatroomID = localStorage.getItem("chatroomID"); // retrieve chatroomID from localStorage
  var username = localStorage.getItem("username"); // retrieve username from localStorage
  $.post("/sendMessage", {
    chatroomID: chatroomID,
    username: username,
    message: message,
  })
    .done(function () {
      alert("Message sent");
    })
    .fail(function () {
      alert("Error sending message");
    });
});

$("#receive-button").click(function () {
  var chatroomID = localStorage.getItem("chatroomID"); // retrieve chatroomID from localStorage
  $.post("/retrieveMessages", { chatroomID: chatroomID })
    .done(function (messages) {
      // Split the messages string into an array of individual messages
      var messagesArray = messages.split("\n");

      // Clear the chat room
      $("#chat-room").empty();

      // Append each message to the chat room
      for (var i = 0; i < messagesArray.length; i++) {
        $("#chat-room").append(
          '<p class="chat-message">' + messagesArray[i] + "</p>"
        );
      }
    })
    .fail(function () {
      alert("Error retrieving messages");
    });
});

window.onbeforeunload = function () {
  localStorage.removeItem("token");
  localStorage.removeItem("chatroomID");
};
