// Index page: the header and the edge panels open when the cursor reaches that edge of the screen.
// The handles on the edges open the panels with a click or tap as well.
(() => {
  const EDGE = 16;
  const CLOSE_DELAY = 300;

  const header = document.querySelector("header");
  const backdrop = document.getElementById("edge-backdrop");
  const panels = {
    left: document.getElementById("skills-panel"),
    right: document.getElementById("about-panel"),
    bottom: document.getElementById("education-panel"),
  };

  let openPanel = null;
  let closeTimer = null;

  function open(panel) {
    if (openPanel === panel) return;
    close();
    openPanel = panel;
    panel.classList.add("open");
    backdrop.classList.add("open");
    header.classList.remove("show");
  }

  function close() {
    clearTimeout(closeTimer);
    if (!openPanel) return;
    openPanel.classList.remove("open");
    backdrop.classList.remove("open");
    openPanel = null;
  }

  document.addEventListener("mousemove", (event) => {
    if (openPanel || document.getElementById("auth-modal")?.open) return;
    const x = event.clientX;
    const y = event.clientY;

    if (y <= EDGE) {
      header.classList.add("show");
      return;
    }
    if (header.classList.contains("show")) {
      const loginOpen = document.getElementById("login-modal")?.classList.contains("show");
      if (y <= header.offsetHeight || loginOpen) return;
      header.classList.remove("show");
    }

    if (x <= EDGE) open(panels.left);
    else if (x >= window.innerWidth - EDGE) open(panels.right);
    else if (y >= window.innerHeight - EDGE) open(panels.bottom);
  });

  Object.values(panels).forEach((panel) => {
    panel.addEventListener("mouseleave", () => {
      closeTimer = setTimeout(close, CLOSE_DELAY);
    });
    panel.addEventListener("mouseenter", () => clearTimeout(closeTimer));
  });

  document.querySelectorAll("[data-edge-open]").forEach((handle) => {
    handle.addEventListener("click", () => {
      const panel = document.getElementById(handle.dataset.edgeOpen);
      open(panel);
      panel.querySelector("[data-edge-close]").focus();
    });
  });

  document.querySelectorAll("[data-edge-close]").forEach((button) => {
    button.addEventListener("click", close);
  });

  backdrop.addEventListener("click", close);

  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape") {
      close();
      header.classList.remove("show");
    }
  });
})();
