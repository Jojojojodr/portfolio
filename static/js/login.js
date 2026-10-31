// Delegated listener registered once, so it keeps working when htmx swaps in a new header/modal
// and this script gets executed again.
if (!window.loginModalInitialized) {
  window.loginModalInitialized = true;

  // Shows the login or register tab of the auth modal
  const showAuthTab = (authModal, tab) => {
    authModal.querySelectorAll("[data-auth-panel]").forEach((panel) => {
      panel.hidden = panel.dataset.authPanel !== tab;
    });
    authModal.querySelectorAll("[data-auth-tab]").forEach((button) => {
      const selected = button.dataset.authTab === tab;
      button.classList.toggle("active", selected);
      button.setAttribute("aria-selected", selected);
    });
    authModal.querySelector(`[data-auth-panel="${tab}"] input:not([type="hidden"])`)?.focus();
  };

  document.addEventListener("click", (event) => {
    const authModal = document.getElementById("auth-modal");
    const loginModal = document.getElementById("login-modal");

    const authOpen = event.target.closest("[data-auth-open]");
    if (authModal && authOpen) {
      event.preventDefault();
      loginModal?.classList.remove("show");
      if (!authModal.open) authModal.showModal();
      showAuthTab(authModal, authOpen.dataset.authOpen);
      return;
    }

    const authTab = event.target.closest("[data-auth-tab]");
    if (authModal && authTab) {
      showAuthTab(authModal, authTab.dataset.authTab);
      return;
    }

    // Close with the close button or a click on the backdrop, which targets the dialog itself
    if (authModal && (event.target.closest("[data-auth-close]") || event.target === authModal)) {
      authModal.close();
      return;
    }

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
