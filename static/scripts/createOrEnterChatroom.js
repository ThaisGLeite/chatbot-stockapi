// Check if user is logged in
if (!sessionStorage.getItem("token")) {
  // If not, redirect to login page
  console.log("User not logged in");
  window.location.href = "login.html";
}

// Helper function to check if the entered chatroom exists
function isChatroomExist(chatroomName, chatroomNames) {
  return chatroomNames && chatroomNames.includes(chatroomName);
}

$(document).ready(function () {
  let chatroomNames = [];

  $.get("/getAllChatrooms")
    .done(function (data) {
      if (data) {
        try {
          chatroomNames = JSON.parse(data) || [];
          if (chatroomNames.length > 0) {
            chatroomNames.forEach(function (chatroomName) {
              $("#chatroom-list").append("<li>" + chatroomName + "</li>");
            });
          } else {
            console.log("No chatrooms exist.");
          }
        } catch (err) {
          console.error("An error occurred while processing the data: ", err);
        }
      } else {
        console.log("No chatrooms exist.");
        chatroomNames = []; // ensure chatroomNames is an empty array when no data received
      }
    })
    .fail(function () {
      alert("Error getting chatrooms");
    });

  $("#create-chatroom-button").click(function () {
    var chatroomName = $("#chatroom-input").val();

    // Check if chatroomName already exists
    $.post("/checkChatroomExist", { chatroomName: chatroomName })
      .done(function (exists) {
        if (exists === "true") {
          alert("Chatroom already exists");
          return;
        }

        $.post("/createChatroom", { chatroomName: chatroomName })
          .done(function (chatroomId) {
            // Save chatroom Name in session storage
            sessionStorage.setItem("chatroomName", chatroomName);

            // Add new chatroom to the list
            $("#chatroom-list").append("<li>" + chatroomName + "</li>");
            chatroomNames.push(chatroomName);

            // Redirect to chatroom.html
            window.location.href = "chatroom.html";
          })
          .fail(function () {
            alert("Error creating chatroom");
          });
      })
      .fail(function () {
        alert("Error checking if chatroom exists");
      });
  });

  $("#enter-chatroom-button").click(function () {
    var chatroomName = $("#chatroom-input").val();
    if (chatroomName) {
      // Check if entered chatroom exists
      if (!isChatroomExist(chatroomName, chatroomNames)) {
        alert("Chatroom does not exist");
        return;
      }

      // Save chatroom Name in session storage
      sessionStorage.setItem("chatroomName", chatroomName);

      // Redirect to chatroom.html
      window.location.href = "chatroom.html";
    } else {
      alert("Please enter a chatroom name");
    }
  });
});
