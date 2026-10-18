// Delegated listener registered once, so it keeps working when htmx swaps in a new header/modal
// and this script gets executed again.
if (!window.loginModalInitialized) {
  window.loginModalInitialized = true;

  document.addEventListener("click", (event) => {
    const loginModal = document.getElementById("login-modal");
    if (!loginModal) {
      return;
    }

    if (event.target.closest("#login-modal-open")) {
      loginModal.classList.add("show");
      return;
    }

    if (event.target !== loginModal && !loginModal.contains(event.target)) {
      loginModal.classList.remove("show");
    }
  });
}
