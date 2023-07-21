// Check if user is logged in
if (!localStorage.getItem("token")) {
  // If not, redirect to login page
  window.location.href = "login.html";
}

// Rest of the chat room functionality here

window.onbeforeunload = function () {
  localStorage.removeItem("token");
};

$("#create-chatroom-button").click(function () {
  var chatroomName = $("#chatroom-input").val();
  $.post("/createChatroom", { chatroomName: chatroomName })
    .done(function (chatroomId) {
      alert("Chatroom created with ID: " + chatroomId);
    })
    .fail(function () {
      alert("Error creating chatroom");
    });
});
