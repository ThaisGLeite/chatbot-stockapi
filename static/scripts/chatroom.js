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
      // message sent
    })
    .fail(function () {
      alert("Error sending message");
    });
});

$("#receive-button").click(function () {
  var chatroomID = sessionStorage.getItem("chatroomID"); // retrieve chatroomID from sessionStorage
  $.post("/retrieveMessages", { chatroomID: chatroomID })
    .done(function (messages) {
      // Parse the JSON response
      var messagesArray = JSON.parse(messages);

      // Clear the chat room
      $("#chat-room").empty();

      // Append each message to the chat room
      for (var i = 0; i < messagesArray.length; i++) {
        var messageTime = new Date(
          messagesArray[i].timestamp * 1000
        ).toLocaleString();
        // Convert the timestamp to milliseconds and then to a local date/time string
        $("#chat-room").append(
          '<p class="chat-message">' +
            messageTime +
            " - " +
            messagesArray[i].username +
            ": " +
            messagesArray[i].message +
            "</p>"
        );
      }
    })
    .fail(function () {
      alert("Error retrieving messages");
    });
});
