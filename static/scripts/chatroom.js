// Check if user is logged in
if (!sessionStorage.getItem("token")) {
  // If not, redirect to login page
  window.location.href = "login.html";
} else {
  // Get chatroom ID and append it to the title
  var chatroomID = sessionStorage.getItem("chatroomID");
  $("#chatroom-title").text("Chatroom: " + chatroomID);
}

$("#send-button").click(function () {
  console.log("clicked the send button");
  var message = $("#message-input").val();
  var chatroomID = sessionStorage.getItem("chatroomID"); // retrieve chatroomID from sessionStorage
  var username = sessionStorage.getItem("username"); // retrieve username from sessionStorage
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
  var chatroomID = sessionStorage.getItem("chatroomID"); // retrieve chatroomID from sessionStorage
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
