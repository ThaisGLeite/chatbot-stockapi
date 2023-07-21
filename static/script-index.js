// Check if user is logged in
if (!localStorage.getItem("token")) {
  // If not, redirect to login page
  window.location.href = "login.html";
}

// Rest of the chat room functionality here

window.onbeforeunload = function () {
  localStorage.removeItem("token");
};
