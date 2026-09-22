const loginOpenBtn = document.getElementById("login-modal-open");
const loginModal = document.getElementById("login-modal");

loginOpenBtn.addEventListener("click", () => {
  loginModal.classList.add("show");
});

document.addEventListener("click", (event) => {
  if (event.target === loginModal || loginOpenBtn.contains(event.target)) {
    return;
  }
  if (!loginModal.contains(event.target)) {
    loginModal.classList.remove("show");
  }
});
