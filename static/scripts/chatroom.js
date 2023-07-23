// Check if user is logged in
if (!sessionStorage.getItem("token")) {
  // If not, redirect to login page
  window.location.href = "login.html";
} else {
  // Get chatroom name and append it to the title
  var chatroomName = sessionStorage.getItem("chatroomName");
  if (chatroomName) {
    $("#chatroom-title").text("Chatroom: " + chatroomName);
  } else {
    console.error("No chatroomName found in sessionStorage");
  }
}

$("#send-button").click(function () {
  console.log("clicked the send button");
  var message = $("#message-input").val();
  var chatroomName = sessionStorage.getItem("chatroomName"); // retrieve chatroomName from sessionStorage
  var username = sessionStorage.getItem("username"); // retrieve username from sessionStorage

  $.post("/sendMessage", {
    chatroomName: chatroomName,
    username: username,
    message: message,
  })
    .done(function () {
      // message sent successfully, clear the input field
      $("#message-input").val("");

      // retrieve the latest messages automatically
      retrieveMessages();
    })
    .fail(function () {
      alert("Error sending message");
    });
});

function retrieveMessages() {
  var chatroomName = sessionStorage.getItem("chatroomName"); // retrieve chatroomName from sessionStorage
  $.post("/retrieveMessages", { chatroomName: chatroomName })
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
}

// Call retrieveMessages when the receive button is clicked
$("#receive-button").click(retrieveMessages);
