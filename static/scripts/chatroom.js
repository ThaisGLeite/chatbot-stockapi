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

// Retrieve messages from the server
$("#send-button").click(function () {
  var message = $("#message-input").val().trim(); // trim whitespace from beginning and end
  var chatroomName = sessionStorage.getItem("chatroomName"); // retrieve chatroomName from sessionStorage
  var username = sessionStorage.getItem("username"); // retrieve username from sessionStorage

  // If message is empty or only whitespace, return and do nothing
  if (!message) {
    return;
  }

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

// Retrieve messages from the server
function retrieveMessages() {
  var chatroomName = sessionStorage.getItem("chatroomName"); // retrieve chatroomName from sessionStorage
  $.post("/retrieveMessages", { chatroomName: chatroomName })
    .done(function (messages) {
      // Parse the JSON response
      var messagesArray = JSON.parse(messages);

      if (Array.isArray(messagesArray)) {
        // Check if messagesArray is an array
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
      } else if (messagesArray !== null) {
        // if messagesArray is not an array, but not null
        console.error("The server response is not an array: ", messagesArray);
      }
    })
    .fail(function () {
      alert("Error retrieving messages");
    });
}

// Call retrieveMessages when the receive button is clicked
$("#receive-button").click(retrieveMessages);
